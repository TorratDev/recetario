package com.recipeapp.domain.repository

import com.recipeapp.domain.common.UiResult
import com.recipeapp.domain.model.NewRecipe
import com.recipeapp.domain.model.Recipe
import com.recipeapp.domain.model.RecipeSummary

interface RecipeRepository {
    suspend fun getRecipes(search: String? = null, limit: Int = 20, offset: Int = 0): UiResult<List<RecipeSummary>>

    suspend fun getRecipe(id: String): UiResult<Recipe>

    suspend fun createRecipe(recipe: NewRecipe): UiResult<Recipe>
}
