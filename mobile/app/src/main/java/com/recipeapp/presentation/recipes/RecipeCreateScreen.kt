package com.recipeapp.presentation.recipes

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel

@Composable
fun RecipeCreateScreen(
    onCreated: (String) -> Unit,
    viewModel: RecipeCreateViewModel = hiltViewModel()
) {
    val state by viewModel.uiState.collectAsState()

    LaunchedEffect(state.createdRecipeId) {
        state.createdRecipeId?.let(onCreated)
    }

    Column(
        modifier = Modifier.fillMaxWidth().padding(24.dp),
        verticalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        Text(text = "New recipe", style = MaterialTheme.typography.headlineMedium)

        OutlinedTextField(
            value = state.title,
            onValueChange = viewModel::onTitleChange,
            label = { Text("Title") },
            modifier = Modifier.fillMaxWidth()
        )
        OutlinedTextField(
            value = state.description,
            onValueChange = viewModel::onDescriptionChange,
            label = { Text("Description") },
            modifier = Modifier.fillMaxWidth()
        )
        OutlinedTextField(
            value = state.cookTime,
            onValueChange = viewModel::onCookTimeChange,
            label = { Text("Cook time (minutes)") },
            modifier = Modifier.fillMaxWidth()
        )

        state.error?.let { Text(text = it, color = MaterialTheme.colorScheme.error) }

        Button(
            onClick = viewModel::save,
            enabled = !state.isSaving && state.title.isNotBlank(),
            modifier = Modifier.fillMaxWidth()
        ) {
            if (state.isSaving) {
                CircularProgressIndicator(modifier = Modifier.padding(2.dp))
            } else {
                Text("Save")
            }
        }
    }
}
