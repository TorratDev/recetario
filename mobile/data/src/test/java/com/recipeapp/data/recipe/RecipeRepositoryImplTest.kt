package com.recipeapp.data.recipe

import com.google.gson.Gson
import com.recipeapp.domain.common.UiResult
import com.recipeapp.domain.model.NewRecipe
import com.recipeapp.network.api.RecipeApi
import com.recipeapp.network.dto.CreateRecipeRequestDto
import com.recipeapp.network.dto.RecipeDto
import com.recipeapp.network.dto.RecipeSummaryDto
import kotlinx.coroutines.test.runTest
import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.ResponseBody.Companion.toResponseBody
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test
import retrofit2.HttpException
import retrofit2.Response

private class FakeRecipeApi(
    private val listResult: Result<List<RecipeSummaryDto>> = Result.failure(IllegalStateException("not stubbed")),
    private val getResult: Result<RecipeDto> = Result.failure(IllegalStateException("not stubbed")),
    private val createResult: Result<RecipeDto> = Result.failure(IllegalStateException("not stubbed"))
) : RecipeApi {
    override suspend fun getRecipes(search: String?, limit: Int, offset: Int): List<RecipeSummaryDto> =
        listResult.getOrThrow()

    override suspend fun getRecipe(id: String): RecipeDto = getResult.getOrThrow()

    override suspend fun createRecipe(request: CreateRecipeRequestDto): RecipeDto = createResult.getOrThrow()
}

private fun httpException(code: Int, body: String): HttpException {
    val responseBody = body.toResponseBody("application/json".toMediaTypeOrNull())
    return HttpException(Response.error<Any>(code, responseBody))
}

private val sampleRecipeDto = RecipeDto(
    id = "r1", userId = "u1", title = "Tarte", description = "desc",
    createdAt = "2026-01-01T00:00:00Z", updatedAt = "2026-01-01T00:00:00Z"
)

class RecipeRepositoryImplTest {

    @Test
    fun `getRecipes maps summaries`() = runTest {
        val summary = RecipeSummaryDto(id = "r1", userId = "u1", title = "Tarte", createdAt = "2026-01-01T00:00:00Z")
        val api = FakeRecipeApi(listResult = Result.success(listOf(summary)))
        val repo = RecipeRepositoryImpl(api, Gson())

        val result = repo.getRecipes()

        assertTrue(result is UiResult.Success)
        result as UiResult.Success
        assertEquals(1, result.data.size)
        assertEquals("Tarte", result.data.first().title)
    }

    @Test
    fun `getRecipe surfaces backend error message`() = runTest {
        val api = object : RecipeApi by FakeRecipeApi() {
            override suspend fun getRecipe(id: String): RecipeDto {
                throw httpException(404, """{"user_message":"Recipe not found"}""")
            }
        }
        val repo = RecipeRepositoryImpl(api, Gson())

        val result = repo.getRecipe("missing")

        assertTrue(result is UiResult.Error)
        assertEquals("Recipe not found", (result as UiResult.Error).message)
    }

    @Test
    fun `createRecipe maps created recipe`() = runTest {
        val api = FakeRecipeApi(createResult = Result.success(sampleRecipeDto))
        val repo = RecipeRepositoryImpl(api, Gson())

        val result = repo.createRecipe(NewRecipe(title = "Tarte", description = "desc"))

        assertTrue(result is UiResult.Success)
        assertEquals("r1", (result as UiResult.Success).data.id)
    }
}
