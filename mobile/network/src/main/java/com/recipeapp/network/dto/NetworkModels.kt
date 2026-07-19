package com.recipeapp.network.dto

import com.google.gson.annotations.SerializedName

/**
 * DTOs mirror the real backend JSON exactly (backend/internal/models,
 * backend/api/openapi.yaml) — string ids, flat `tags: []string`, a separate
 * ordered `instructions` array, and `amount` as a free-text string, not a
 * numeric quantity. See mobile/domain/.../model/Models.kt for the mapped
 * domain shapes.
 */

// --- Auth ---

data class LoginRequestDto(
    @SerializedName("email") val email: String,
    @SerializedName("password") val password: String
)

data class RegisterRequestDto(
    @SerializedName("email") val email: String,
    @SerializedName("username") val username: String,
    @SerializedName("password") val password: String,
    @SerializedName("first_name") val firstName: String,
    @SerializedName("last_name") val lastName: String
)

data class UserDto(
    @SerializedName("id") val id: String,
    @SerializedName("email") val email: String,
    @SerializedName("username") val username: String,
    @SerializedName("first_name") val firstName: String,
    @SerializedName("last_name") val lastName: String,
    @SerializedName("avatar_url") val avatarUrl: String? = null
)

data class AuthResponseDto(
    @SerializedName("token") val token: String,
    @SerializedName("user") val user: UserDto,
    @SerializedName("expires_in") val expiresIn: Long
)

/** No `user` field — unlike AuthResponseDto, matching the real refresh response. */
data class RefreshResponseDto(
    @SerializedName("token") val token: String,
    @SerializedName("expires_in") val expiresIn: Long
)

// --- Recipes ---

data class IngredientDto(
    @SerializedName("id") val id: String,
    @SerializedName("recipe_id") val recipeId: String? = null,
    @SerializedName("name") val name: String,
    @SerializedName("amount") val amount: String? = null,
    @SerializedName("unit") val unit: String? = null,
    @SerializedName("notes") val notes: String? = null,
    @SerializedName("position") val position: Int = 0
)

data class InstructionDto(
    @SerializedName("id") val id: String,
    @SerializedName("recipe_id") val recipeId: String? = null,
    @SerializedName("text") val text: String,
    @SerializedName("position") val position: Int = 0,
    @SerializedName("duration") val duration: Int? = null,
    @SerializedName("temperature") val temperature: Int? = null
)

/** Full shape returned by GET /recipes/{id}/ and by create/update. */
data class RecipeDto(
    @SerializedName("id") val id: String,
    @SerializedName("user_id") val userId: String,
    @SerializedName("is_owner") val isOwner: Boolean = false,
    @SerializedName("title") val title: String,
    @SerializedName("description") val description: String,
    @SerializedName("prep_time") val prepTime: Int = 0,
    @SerializedName("cook_time") val cookTime: Int = 0,
    @SerializedName("servings") val servings: Int = 0,
    @SerializedName("difficulty") val difficulty: String = "",
    @SerializedName("category") val category: String = "",
    @SerializedName("cuisine") val cuisine: String = "",
    @SerializedName("image_url") val imageUrl: String? = null,
    @SerializedName("ingredients") val ingredients: List<IngredientDto> = emptyList(),
    @SerializedName("instructions") val instructions: List<InstructionDto> = emptyList(),
    @SerializedName("tags") val tags: List<String> = emptyList(),
    @SerializedName("created_at") val createdAt: String,
    @SerializedName("updated_at") val updatedAt: String
)

/** Narrower shape returned by GET /recipes/ (list) — fields are dynamically built server-side. */
data class RecipeSummaryDto(
    @SerializedName("id") val id: String,
    @SerializedName("user_id") val userId: String,
    @SerializedName("is_owner") val isOwner: Boolean = false,
    @SerializedName("title") val title: String,
    @SerializedName("description") val description: String? = null,
    @SerializedName("cook_time") val cookTime: Int = 0,
    @SerializedName("difficulty") val difficulty: String = "",
    @SerializedName("category") val category: String? = null,
    @SerializedName("cuisine") val cuisine: String? = null,
    @SerializedName("image_url") val imageUrl: String? = null,
    @SerializedName("created_at") val createdAt: String
)

data class CreateRecipeRequestDto(
    @SerializedName("title") val title: String,
    @SerializedName("description") val description: String,
    @SerializedName("prep_time") val prepTime: Int = 0,
    @SerializedName("cook_time") val cookTime: Int = 0,
    @SerializedName("servings") val servings: Int = 0,
    @SerializedName("difficulty") val difficulty: String = "",
    @SerializedName("category") val category: String = "",
    @SerializedName("cuisine") val cuisine: String = "",
    @SerializedName("image_url") val imageUrl: String = ""
)

// --- Errors ---

/** Mirrors backend/internal/appmiddleware.ErrorResponse — the single shared error envelope. */
data class ErrorResponseDto(
    @SerializedName("error") val error: String? = null,
    @SerializedName("message") val message: String? = null,
    @SerializedName("code") val code: String? = null,
    @SerializedName("description") val description: String? = null,
    @SerializedName("user_message") val userMessage: String? = null
)
