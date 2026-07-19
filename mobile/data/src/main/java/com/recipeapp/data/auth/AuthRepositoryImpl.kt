package com.recipeapp.data.auth

import com.google.gson.Gson
import com.recipeapp.data.common.toUiError
import com.recipeapp.domain.common.UiResult
import com.recipeapp.domain.model.AuthSession
import com.recipeapp.domain.model.LoginCredentials
import com.recipeapp.domain.model.RegisterDetails
import com.recipeapp.domain.model.User
import com.recipeapp.domain.repository.AuthRepository
import com.recipeapp.network.api.AuthApi
import com.recipeapp.network.dto.AuthResponseDto
import com.recipeapp.network.dto.LoginRequestDto
import com.recipeapp.network.dto.RegisterRequestDto
import com.recipeapp.network.dto.UserDto
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class AuthRepositoryImpl @Inject constructor(
    private val authApi: AuthApi,
    private val tokenStore: SessionStore,
    private val gson: Gson
) : AuthRepository {

    override val session: Flow<AuthSession?> = tokenStore.session.map { it }

    override suspend fun login(credentials: LoginCredentials): UiResult<AuthSession> =
        runCatching {
            authApi.login(LoginRequestDto(credentials.email, credentials.password))
        }.fold(
            onSuccess = { dto -> dto.toSession().also { tokenStore.save(it) }.let { UiResult.Success(it) } },
            onFailure = { it.toUiError(gson) }
        )

    override suspend fun register(details: RegisterDetails): UiResult<AuthSession> =
        runCatching {
            authApi.register(
                RegisterRequestDto(
                    email = details.email,
                    username = details.username,
                    password = details.password,
                    firstName = details.firstName,
                    lastName = details.lastName
                )
            )
        }.fold(
            onSuccess = { dto -> dto.toSession().also { tokenStore.save(it) }.let { UiResult.Success(it) } },
            onFailure = { it.toUiError(gson) }
        )

    /**
     * The backend only re-signs a token that is still valid — it has no
     * separate long-lived refresh token. If the current token already
     * expired, this call 401s too; on that failure we clear the local
     * session so the caller falls back to a full re-login instead of
     * retrying forever.
     */
    override suspend fun refresh(): UiResult<Unit> =
        runCatching { authApi.refresh() }.fold(
            onSuccess = { dto ->
                tokenStore.updateToken(dto.token, dto.expiresIn)
                UiResult.Success(Unit)
            },
            onFailure = { error ->
                tokenStore.clear()
                error.toUiError(gson)
            }
        )

    override suspend fun logout() {
        runCatching { authApi.logout() }
        tokenStore.clear()
    }

    private fun AuthResponseDto.toSession() = AuthSession(
        token = token,
        user = user.toDomain(),
        expiresIn = expiresIn
    )

    private fun UserDto.toDomain() = User(
        id = id,
        email = email,
        username = username,
        firstName = firstName,
        lastName = lastName,
        avatarUrl = avatarUrl
    )
}
