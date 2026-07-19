package com.recipeapp.data.recipe

import com.google.gson.Gson
import com.recipeapp.data.common.toUiError
import com.recipeapp.domain.common.UiResult
import com.recipeapp.domain.model.NewRecipe
import com.recipeapp.domain.model.Recipe
import com.recipeapp.domain.model.RecipeIngredient
import com.recipeapp.domain.model.RecipeInstruction
import com.recipeapp.domain.model.RecipeSummary
import com.recipeapp.domain.repository.RecipeRepository
import com.recipeapp.network.api.RecipeApi
import com.recipeapp.network.dto.CreateRecipeRequestDto
import com.recipeapp.network.dto.IngredientDto
import com.recipeapp.network.dto.InstructionDto
import com.recipeapp.network.dto.RecipeDto
import com.recipeapp.network.dto.RecipeSummaryDto
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class RecipeRepositoryImpl @Inject constructor(
    private val recipeApi: RecipeApi,
    private val gson: Gson
) : RecipeRepository {

    override suspend fun getRecipes(search: String?, limit: Int, offset: Int): UiResult<List<RecipeSummary>> =
        runCatching { recipeApi.getRecipes(search, limit, offset) }.fold(
            onSuccess = { list -> UiResult.Success(list.map { it.toDomain() }) },
            onFailure = { it.toUiError(gson) }
        )

    override suspend fun getRecipe(id: String): UiResult<Recipe> =
        runCatching { recipeApi.getRecipe(id) }.fold(
            onSuccess = { UiResult.Success(it.toDomain()) },
            onFailure = { it.toUiError(gson) }
        )

    override suspend fun createRecipe(recipe: NewRecipe): UiResult<Recipe> =
        runCatching {
            recipeApi.createRecipe(
                CreateRecipeRequestDto(
                    title = recipe.title,
                    description = recipe.description,
                    prepTime = recipe.prepTime,
                    cookTime = recipe.cookTime,
                    servings = recipe.servings,
                    difficulty = recipe.difficulty,
                    category = recipe.category,
                    cuisine = recipe.cuisine,
                    imageUrl = recipe.imageUrl
                )
            )
        }.fold(
            onSuccess = { UiResult.Success(it.toDomain()) },
            onFailure = { it.toUiError(gson) }
        )

    private fun RecipeSummaryDto.toDomain() = RecipeSummary(
        id = id,
        userId = userId,
        isOwner = isOwner,
        title = title,
        description = description ?: "",
        cookTime = cookTime,
        difficulty = difficulty,
        category = category ?: "",
        cuisine = cuisine ?: "",
        imageUrl = imageUrl,
        createdAt = createdAt
    )

    private fun RecipeDto.toDomain() = Recipe(
        id = id,
        userId = userId,
        isOwner = isOwner,
        title = title,
        description = description,
        prepTime = prepTime,
        cookTime = cookTime,
        servings = servings,
        difficulty = difficulty,
        category = category,
        cuisine = cuisine,
        imageUrl = imageUrl,
        ingredients = ingredients.map { it.toDomain() },
        instructions = instructions.map { it.toDomain() },
        tags = tags,
        createdAt = createdAt,
        updatedAt = updatedAt
    )

    private fun IngredientDto.toDomain() = RecipeIngredient(
        id = id,
        name = name,
        amount = amount ?: "",
        unit = unit ?: "",
        notes = notes ?: "",
        position = position
    )

    private fun InstructionDto.toDomain() = RecipeInstruction(
        id = id,
        text = text,
        position = position,
        durationMinutes = duration,
        temperature = temperature
    )
}
