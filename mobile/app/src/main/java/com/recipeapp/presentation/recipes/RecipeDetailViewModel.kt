package com.recipeapp.presentation.recipes

import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.recipeapp.domain.common.UiResult
import com.recipeapp.domain.model.Recipe
import com.recipeapp.domain.repository.RecipeRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class RecipeDetailUiState(
    val recipe: Recipe? = null,
    val isLoading: Boolean = true,
    val error: String? = null
)

@HiltViewModel
class RecipeDetailViewModel @Inject constructor(
    private val recipeRepository: RecipeRepository,
    savedStateHandle: SavedStateHandle
) : ViewModel() {

    private val recipeId: String = checkNotNull(savedStateHandle["recipeId"])

    private val _uiState = MutableStateFlow(RecipeDetailUiState())
    val uiState: StateFlow<RecipeDetailUiState> = _uiState.asStateFlow()

    init {
        load()
    }

    fun retry() = load()

    private fun load() {
        _uiState.value = _uiState.value.copy(isLoading = true, error = null)
        viewModelScope.launch {
            when (val result = recipeRepository.getRecipe(recipeId)) {
                is UiResult.Success -> _uiState.value =
                    _uiState.value.copy(isLoading = false, recipe = result.data)
                is UiResult.Error -> _uiState.value =
                    _uiState.value.copy(isLoading = false, error = result.message)
            }
        }
    }
}
