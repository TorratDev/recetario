package com.recipeapp.presentation.auth

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation

@Composable
fun RegisterScreen(
    onRegistered: () -> Unit,
    onNavigateToLogin: () -> Unit,
    viewModel: AuthViewModel = hiltViewModel()
) {
    val state by viewModel.uiState.collectAsState()

    LaunchedEffect(state.isAuthenticated) {
        if (state.isAuthenticated) onRegistered()
    }

    Column(
        modifier = Modifier
            .fillMaxWidth()
            .padding(24.dp),
        verticalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        Text(text = "Create account", style = MaterialTheme.typography.headlineMedium)

        OutlinedTextField(
            value = state.email,
            onValueChange = viewModel::onEmailChange,
            label = { Text("Email") },
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Email),
            modifier = Modifier.fillMaxWidth()
        )
        OutlinedTextField(
            value = state.username,
            onValueChange = viewModel::onUsernameChange,
            label = { Text("Username") },
            modifier = Modifier.fillMaxWidth()
        )
        OutlinedTextField(
            value = state.firstName,
            onValueChange = viewModel::onFirstNameChange,
            label = { Text("First name") },
            modifier = Modifier.fillMaxWidth()
        )
        OutlinedTextField(
            value = state.lastName,
            onValueChange = viewModel::onLastNameChange,
            label = { Text("Last name") },
            modifier = Modifier.fillMaxWidth()
        )
        OutlinedTextField(
            value = state.password,
            onValueChange = viewModel::onPasswordChange,
            label = { Text("Password") },
            visualTransformation = PasswordVisualTransformation(),
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
            modifier = Modifier.fillMaxWidth()
        )

        state.error?.let { Text(text = it, color = MaterialTheme.colorScheme.error) }

        Button(
            onClick = viewModel::register,
            enabled = !state.isLoading &&
                state.email.isNotBlank() &&
                state.username.isNotBlank() &&
                state.password.isNotBlank(),
            modifier = Modifier.fillMaxWidth()
        ) {
            if (state.isLoading) {
                CircularProgressIndicator(modifier = Modifier.padding(2.dp))
            } else {
                Text("Register")
            }
        }

        TextButton(onClick = onNavigateToLogin, modifier = Modifier.fillMaxWidth()) {
            Text("Already have an account? Log in")
        }
    }
}
