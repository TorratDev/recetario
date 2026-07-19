package com.recipeapp.data.auth

import com.google.gson.Gson
import com.recipeapp.domain.common.UiResult
import com.recipeapp.domain.model.AuthSession
import com.recipeapp.domain.model.LoginCredentials
import com.recipeapp.domain.model.User
import com.recipeapp.network.api.AuthApi
import com.recipeapp.network.dto.AuthResponseDto
import com.recipeapp.network.dto.LoginRequestDto
import com.recipeapp.network.dto.RefreshResponseDto
import com.recipeapp.network.dto.RegisterRequestDto
import com.recipeapp.network.dto.UserDto
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.test.runTest
import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.ResponseBody.Companion.toResponseBody
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertNull
import org.junit.Assert.assertTrue
import org.junit.Test
import retrofit2.HttpException
import retrofit2.Response

private class FakeSessionStore : SessionStore {
    private val _session = MutableStateFlow<AuthSession?>(null)
    override val session: StateFlow<AuthSession?> = _session

    override fun currentToken(): String? = _session.value?.token
    override fun save(authSession: AuthSession) {
        _session.value = authSession
    }
    override fun updateToken(token: String, expiresIn: Long) {
        _session.value = _session.value?.copy(token = token, expiresIn = expiresIn)
    }
    override fun clear() {
        _session.value = null
    }
}

private class FakeAuthApi(
    private val loginResult: Result<AuthResponseDto> = Result.failure(IllegalStateException("not stubbed")),
    private val refreshResult: Result<RefreshResponseDto> = Result.failure(IllegalStateException("not stubbed"))
) : AuthApi {
    override suspend fun login(request: LoginRequestDto): AuthResponseDto = loginResult.getOrThrow()
    override suspend fun register(request: RegisterRequestDto): AuthResponseDto = loginResult.getOrThrow()
    override suspend fun refresh(): RefreshResponseDto = refreshResult.getOrThrow()
    override suspend fun logout(accept: String) { /* no-op */ }
}

private val sampleUser = UserDto(
    id = "user-1", email = "a@b.com", username = "ab",
    firstName = "A", lastName = "B"
)

class AuthRepositoryImplTest {

    private fun httpException(code: Int, body: String): HttpException {
        val responseBody = body.toResponseBody("application/json".toMediaTypeOrNull())
        return HttpException(Response.error<Any>(code, responseBody))
    }

    @Test
    fun `login success saves session`() = runTest {
        val store = FakeSessionStore()
        val api = FakeAuthApi(loginResult = Result.success(AuthResponseDto("token-1", sampleUser, 3600)))
        val repo = AuthRepositoryImpl(api, store, Gson())

        val result = repo.login(LoginCredentials("a@b.com", "pw"))

        assertTrue(result is UiResult.Success)
        assertEquals("token-1", store.currentToken())
    }

    @Test
    fun `login failure surfaces backend message and does not save session`() = runTest {
        val store = FakeSessionStore()
        val api = object : AuthApi by FakeAuthApi() {
            override suspend fun login(request: LoginRequestDto): AuthResponseDto {
                throw httpException(401, """{"user_message":"Invalid credentials"}""")
            }
        }
        val repo = AuthRepositoryImpl(api, store, Gson())

        val result = repo.login(LoginCredentials("a@b.com", "wrong"))

        assertTrue(result is UiResult.Error)
        result as UiResult.Error
        assertEquals("Invalid credentials", result.message)
        assertTrue(result.isUnauthorized)
        assertNull(store.currentToken())
    }

    @Test
    fun `refresh failure clears the session`() = runTest {
        val store = FakeSessionStore()
        store.save(AuthSession(token = "stale", user = User(sampleUser.id, sampleUser.email, sampleUser.username, sampleUser.firstName, sampleUser.lastName), expiresIn = 0))
        val api = object : AuthApi by FakeAuthApi() {
            override suspend fun refresh(): RefreshResponseDto = throw httpException(401, "{}")
        }
        val repo = AuthRepositoryImpl(api, store, Gson())

        val result = repo.refresh()

        assertTrue(result is UiResult.Error)
        assertNull(store.currentToken())
    }

    @Test
    fun `logout always clears the local session`() = runTest {
        val store = FakeSessionStore()
        store.save(AuthSession(token = "t", user = User("1", "a@b.com", "ab", "A", "B"), expiresIn = 60))
        val api = object : AuthApi by FakeAuthApi() {
            override suspend fun logout(accept: String) { throw httpException(500, "{}") }
        }
        val repo = AuthRepositoryImpl(api, store, Gson())

        repo.logout()

        assertFalse(store.session.value != null)
    }
}
