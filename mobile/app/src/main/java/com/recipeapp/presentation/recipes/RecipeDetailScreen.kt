package com.recipeapp.presentation.recipes

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel

@Composable
fun RecipeDetailScreen(
    viewModel: RecipeDetailViewModel = hiltViewModel()
) {
    val state by viewModel.uiState.collectAsState()

    when {
        state.isLoading -> Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
            CircularProgressIndicator()
        }

        state.error != null -> Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
            Column(horizontalAlignment = Alignment.CenterHorizontally) {
                Text(state.error ?: "", color = MaterialTheme.colorScheme.error)
                TextButton(onClick = viewModel::retry) { Text("Retry") }
            }
        }

        state.recipe != null -> {
            val recipe = state.recipe!!
            Column(modifier = Modifier.fillMaxSize().padding(16.dp)) {
                Text(text = recipe.title, style = MaterialTheme.typography.headlineSmall)
                Text(text = recipe.description, style = MaterialTheme.typography.bodyMedium)

                Text(text = "Ingredients", style = MaterialTheme.typography.titleMedium, modifier = Modifier.padding(top = 16.dp))
                recipe.ingredients.sortedBy { it.position }.forEach { ingredient ->
                    Text(text = "- ${ingredient.amount} ${ingredient.unit} ${ingredient.name}".trim())
                }

                Text(text = "Instructions", style = MaterialTheme.typography.titleMedium, modifier = Modifier.padding(top = 16.dp))
                recipe.instructions.sortedBy { it.position }.forEachIndexed { index, instruction ->
                    Text(text = "${index + 1}. ${instruction.text}")
                }
            }
        }
    }
}
