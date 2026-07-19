package com.recipeapp.network.api

import com.recipeapp.network.dto.AuthResponseDto
import com.recipeapp.network.dto.CreateRecipeRequestDto
import com.recipeapp.network.dto.LoginRequestDto
import com.recipeapp.network.dto.RecipeDto
import com.recipeapp.network.dto.RecipeSummaryDto
import com.recipeapp.network.dto.RefreshResponseDto
import com.recipeapp.network.dto.RegisterRequestDto
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.Header
import retrofit2.http.Path
import retrofit2.http.POST
import retrofit2.http.Query

interface AuthApi {
    @POST("auth/login")
    suspend fun login(@Body request: LoginRequestDto): AuthResponseDto

    @POST("auth/register")
    suspend fun register(@Body request: RegisterRequestDto): AuthResponseDto

    /**
     * No @Body — the backend re-signs a new token from the current, still-valid
     * one passed via the Authorization header. AuthInterceptor already attaches
     * the stored bearer token to every request, so no explicit header param is
     * needed here either.
     */
    @POST("auth/refresh")
    suspend fun refresh(): RefreshResponseDto

    /**
     * Requires Accept: application/json, or the backend responds with a 303
     * redirect instead of JSON (browser/HTMX default behavior).
     */
    @POST("auth/logout")
    suspend fun logout(@Header("Accept") accept: String = "application/json")
}

interface RecipeApi {
    @GET("recipes/")
    suspend fun getRecipes(
        @Query("search") search: String? = null,
        @Query("limit") limit: Int = 20,
        @Query("offset") offset: Int = 0
    ): List<RecipeSummaryDto>

    @GET("recipes/{id}/")
    suspend fun getRecipe(@Path("id") id: String): RecipeDto

    @POST("recipes/")
    suspend fun createRecipe(@Body request: CreateRecipeRequestDto): RecipeDto
}
