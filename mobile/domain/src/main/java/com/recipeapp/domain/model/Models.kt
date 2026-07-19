package com.recipeapp.domain.model

/**
 * Domain models mirroring the real backend contract (backend/internal/models,
 * backend/api/openapi.yaml) — string/UUID ids, flat [tags], a separate ordered
 * [instructions] list, and RFC3339 timestamps kept as String (the backend
 * serializes Go time.Time as RFC3339; parsing to a richer date type is left to
 * the presentation layer if/when needed, so no kotlinx-datetime dependency is
 * required here).
 */

data class User(
    val id: String,
    val email: String,
    val username: String,
    val firstName: String,
    val lastName: String,
    val avatarUrl: String? = null
)

data class AuthSession(
    val token: String,
    val user: User,
    val expiresIn: Long
)

data class LoginCredentials(
    val email: String,
    val password: String
)

data class RegisterDetails(
    val email: String,
    val username: String,
    val password: String,
    val firstName: String,
    val lastName: String
)

data class RecipeIngredient(
    val id: String,
    val name: String,
    val amount: String,
    val unit: String,
    val notes: String,
    val position: Int
)

data class RecipeInstruction(
    val id: String,
    val text: String,
    val position: Int,
    val durationMinutes: Int? = null,
    val temperature: Int? = null
)

data class Recipe(
    val id: String,
    val userId: String,
    val isOwner: Boolean,
    val title: String,
    val description: String,
    val prepTime: Int,
    val cookTime: Int,
    val servings: Int,
    val difficulty: String,
    val category: String,
    val cuisine: String,
    val imageUrl: String?,
    val ingredients: List<RecipeIngredient> = emptyList(),
    val instructions: List<RecipeInstruction> = emptyList(),
    val tags: List<String> = emptyList(),
    val createdAt: String,
    val updatedAt: String
)

/** Narrower shape returned by the recipe list endpoint (GET /recipes/). */
data class RecipeSummary(
    val id: String,
    val userId: String,
    val isOwner: Boolean,
    val title: String,
    val description: String,
    val cookTime: Int,
    val difficulty: String,
    val category: String,
    val cuisine: String,
    val imageUrl: String?,
    val createdAt: String
)

data class NewRecipe(
    val title: String,
    val description: String,
    val prepTime: Int = 0,
    val cookTime: Int = 0,
    val servings: Int = 0,
    val difficulty: String = "",
    val category: String = "",
    val cuisine: String = "",
    val imageUrl: String = ""
)
