package com.recipeapp.domain.common

/**
 * Shared success/error wrapper for repository calls, so every ViewModel
 * handles network failures the same way (issue #54's "errores de red
 * gestionados de forma visible" requirement) instead of each feature
 * inventing its own error type.
 */
sealed class UiResult<out T> {
    data class Success<T>(val data: T) : UiResult<T>()
    data class Error(val message: String, val isUnauthorized: Boolean = false) : UiResult<Nothing>()
}
