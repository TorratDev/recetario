package com.recipeapp.presentation

import androidx.compose.runtime.Composable
import androidx.navigation.NavHostController
import androidx.navigation.NavType
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import androidx.navigation.navArgument
import com.recipeapp.presentation.auth.LoginScreen
import com.recipeapp.presentation.auth.RegisterScreen
import com.recipeapp.presentation.recipes.RecipeCreateScreen
import com.recipeapp.presentation.recipes.RecipeDetailScreen
import com.recipeapp.presentation.recipes.RecipeListScreen

object Routes {
    const val LOGIN = "login"
    const val REGISTER = "register"
    const val RECIPES = "recipes"
    const val RECIPE_CREATE = "recipes/new"
    const val RECIPE_DETAIL = "recipes/{recipeId}"

    fun recipeDetail(id: String) = "recipes/$id"
}

@Composable
fun RecipeNavigation(navController: NavHostController = rememberNavController()) {
    NavHost(navController = navController, startDestination = Routes.LOGIN) {
        composable(Routes.LOGIN) {
            LoginScreen(
                onLoggedIn = {
                    navController.navigate(Routes.RECIPES) {
                        popUpTo(Routes.LOGIN) { inclusive = true }
                    }
                },
                onNavigateToRegister = { navController.navigate(Routes.REGISTER) }
            )
        }
        composable(Routes.REGISTER) {
            RegisterScreen(
                onRegistered = {
                    navController.navigate(Routes.RECIPES) {
                        popUpTo(Routes.LOGIN) { inclusive = true }
                    }
                },
                onNavigateToLogin = { navController.popBackStack() }
            )
        }
        composable(Routes.RECIPES) {
            RecipeListScreen(
                onOpenRecipe = { id -> navController.navigate(Routes.recipeDetail(id)) },
                onCreateRecipe = { navController.navigate(Routes.RECIPE_CREATE) }
            )
        }
        composable(Routes.RECIPE_CREATE) {
            RecipeCreateScreen(
                onCreated = { id ->
                    navController.navigate(Routes.recipeDetail(id)) {
                        popUpTo(Routes.RECIPE_CREATE) { inclusive = true }
                    }
                }
            )
        }
        composable(
            route = Routes.RECIPE_DETAIL,
            arguments = listOf(navArgument("recipeId") { type = NavType.StringType })
        ) {
            RecipeDetailScreen()
        }
    }
}
