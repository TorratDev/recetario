package com.recipeapp.presentation.recipes

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material3.Card
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.recipeapp.domain.model.RecipeSummary

@Composable
fun RecipeListScreen(
    onOpenRecipe: (String) -> Unit,
    onCreateRecipe: () -> Unit,
    viewModel: RecipeListViewModel = hiltViewModel()
) {
    val state by viewModel.uiState.collectAsState()

    Scaffold(
        floatingActionButton = {
            FloatingActionButton(onClick = onCreateRecipe) {
                Icon(Icons.Filled.Add, contentDescription = "Add recipe")
            }
        }
    ) { padding ->
        Column(modifier = Modifier.fillMaxSize().padding(padding).padding(16.dp)) {
            OutlinedTextField(
                value = state.search,
                onValueChange = viewModel::onSearchChange,
                label = { Text("Search recipes") },
                modifier = Modifier.fillMaxWidth()
            )

            when {
                state.isLoading -> Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) { CircularProgressIndicator() }

                state.error != null -> Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Text(state.error ?: "", color = MaterialTheme.colorScheme.error)
                        TextButton(onClick = viewModel::retry) { Text("Retry") }
                    }
                }

                state.recipes.isEmpty() -> Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) { Text("No recipes yet") }

                else -> LazyColumn(
                    modifier = Modifier.fillMaxSize(),
                    verticalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    items(state.recipes, key = { it.id }) { recipe ->
                        RecipeSummaryRow(recipe = recipe, onClick = { onOpenRecipe(recipe.id) })
                    }
                }
            }
        }
    }
}

@Composable
private fun RecipeSummaryRow(recipe: RecipeSummary, onClick: () -> Unit) {
    Card(modifier = Modifier.fillMaxWidth(), onClick = onClick) {
        Column(modifier = Modifier.padding(16.dp)) {
            Text(text = recipe.title, style = MaterialTheme.typography.titleMedium)
            if (recipe.description.isNotBlank()) {
                Text(text = recipe.description, style = MaterialTheme.typography.bodyMedium)
            }
        }
    }
}
