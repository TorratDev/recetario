package com.recipeapp.presentation.auth

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.recipeapp.domain.common.UiResult
import com.recipeapp.domain.model.LoginCredentials
import com.recipeapp.domain.model.RegisterDetails
import com.recipeapp.domain.repository.AuthRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class AuthUiState(
    val email: String = "",
    val password: String = "",
    val username: String = "",
    val firstName: String = "",
    val lastName: String = "",
    val isLoading: Boolean = false,
    val error: String? = null,
    val isAuthenticated: Boolean = false
)

@HiltViewModel
class AuthViewModel @Inject constructor(
    private val authRepository: AuthRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(AuthUiState())
    val uiState: StateFlow<AuthUiState> = _uiState.asStateFlow()

    fun onEmailChange(value: String) {
        _uiState.value = _uiState.value.copy(email = value, error = null)
    }

    fun onPasswordChange(value: String) {
        _uiState.value = _uiState.value.copy(password = value, error = null)
    }

    fun onUsernameChange(value: String) {
        _uiState.value = _uiState.value.copy(username = value, error = null)
    }

    fun onFirstNameChange(value: String) {
        _uiState.value = _uiState.value.copy(firstName = value, error = null)
    }

    fun onLastNameChange(value: String) {
        _uiState.value = _uiState.value.copy(lastName = value, error = null)
    }

    fun login() {
        val state = _uiState.value
        _uiState.value = state.copy(isLoading = true, error = null)
        viewModelScope.launch {
            when (val result = authRepository.login(LoginCredentials(state.email, state.password))) {
                is UiResult.Success -> _uiState.value = _uiState.value.copy(isLoading = false, isAuthenticated = true)
                is UiResult.Error -> _uiState.value = _uiState.value.copy(isLoading = false, error = result.message)
            }
        }
    }

    fun register() {
        val state = _uiState.value
        _uiState.value = state.copy(isLoading = true, error = null)
        viewModelScope.launch {
            val details = RegisterDetails(
                email = state.email,
                username = state.username,
                password = state.password,
                firstName = state.firstName,
                lastName = state.lastName
            )
            when (val result = authRepository.register(details)) {
                is UiResult.Success -> _uiState.value = _uiState.value.copy(isLoading = false, isAuthenticated = true)
                is UiResult.Error -> _uiState.value = _uiState.value.copy(isLoading = false, error = result.message)
            }
        }
    }
}
