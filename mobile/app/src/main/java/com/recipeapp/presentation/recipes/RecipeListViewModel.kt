package com.recipeapp.presentation.recipes

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.recipeapp.domain.common.UiResult
import com.recipeapp.domain.model.RecipeSummary
import com.recipeapp.domain.repository.RecipeRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class RecipeListUiState(
    val recipes: List<RecipeSummary> = emptyList(),
    val search: String = "",
    val isLoading: Boolean = false,
    val error: String? = null
)

@HiltViewModel
class RecipeListViewModel @Inject constructor(
    private val recipeRepository: RecipeRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(RecipeListUiState())
    val uiState: StateFlow<RecipeListUiState> = _uiState.asStateFlow()

    init {
        load()
    }

    fun onSearchChange(value: String) {
        _uiState.value = _uiState.value.copy(search = value)
        load()
    }

    fun retry() = load()

    private fun load() {
        val query = _uiState.value.search.ifBlank { null }
        _uiState.value = _uiState.value.copy(isLoading = true, error = null)
        viewModelScope.launch {
            when (val result = recipeRepository.getRecipes(search = query)) {
                is UiResult.Success -> _uiState.value =
                    _uiState.value.copy(isLoading = false, recipes = result.data)
                is UiResult.Error -> _uiState.value =
                    _uiState.value.copy(isLoading = false, error = result.message)
            }
        }
    }
}
