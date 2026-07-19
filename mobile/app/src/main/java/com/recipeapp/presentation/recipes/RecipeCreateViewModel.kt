package com.recipeapp.presentation.recipes

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.recipeapp.domain.common.UiResult
import com.recipeapp.domain.model.NewRecipe
import com.recipeapp.domain.repository.RecipeRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class RecipeCreateUiState(
    val title: String = "",
    val description: String = "",
    val cookTime: String = "",
    val isSaving: Boolean = false,
    val error: String? = null,
    val createdRecipeId: String? = null
)

@HiltViewModel
class RecipeCreateViewModel @Inject constructor(
    private val recipeRepository: RecipeRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(RecipeCreateUiState())
    val uiState: StateFlow<RecipeCreateUiState> = _uiState.asStateFlow()

    fun onTitleChange(value: String) {
        _uiState.value = _uiState.value.copy(title = value, error = null)
    }

    fun onDescriptionChange(value: String) {
        _uiState.value = _uiState.value.copy(description = value, error = null)
    }

    fun onCookTimeChange(value: String) {
        _uiState.value = _uiState.value.copy(cookTime = value, error = null)
    }

    fun save() {
        val state = _uiState.value
        _uiState.value = state.copy(isSaving = true, error = null)
        viewModelScope.launch {
            val newRecipe = NewRecipe(
                title = state.title,
                description = state.description,
                cookTime = state.cookTime.toIntOrNull() ?: 0
            )
            when (val result = recipeRepository.createRecipe(newRecipe)) {
                is UiResult.Success -> _uiState.value =
                    _uiState.value.copy(isSaving = false, createdRecipeId = result.data.id)
                is UiResult.Error -> _uiState.value =
                    _uiState.value.copy(isSaving = false, error = result.message)
            }
        }
    }
}
