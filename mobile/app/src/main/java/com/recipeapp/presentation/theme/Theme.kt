package com.recipeapp.presentation.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

private val Orange = Color(0xFFFF6F00)
private val OrangeDark = Color(0xFFC43E00)

private val LightColors = lightColorScheme(
    primary = Orange,
    secondary = OrangeDark
)

private val DarkColors = darkColorScheme(
    primary = Orange,
    secondary = OrangeDark
)

@Composable
fun RecipeTheme(
    darkTheme: Boolean = androidx.compose.foundation.isSystemInDarkTheme(),
    content: @Composable () -> Unit
) {
    MaterialTheme(
        colorScheme = if (darkTheme) DarkColors else LightColors,
        content = content
    )
}
