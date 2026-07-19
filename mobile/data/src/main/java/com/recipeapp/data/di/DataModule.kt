package com.recipeapp.data.di

import com.recipeapp.data.auth.AuthRepositoryImpl
import com.recipeapp.data.auth.SessionStore
import com.recipeapp.data.auth.TokenStore
import com.recipeapp.data.recipe.RecipeRepositoryImpl
import com.recipeapp.domain.repository.AuthRepository
import com.recipeapp.domain.repository.RecipeRepository
import com.recipeapp.network.AuthTokenProvider
import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
abstract class DataModule {

    @Binds
    @Singleton
    abstract fun bindAuthRepository(impl: AuthRepositoryImpl): AuthRepository

    @Binds
    @Singleton
    abstract fun bindRecipeRepository(impl: RecipeRepositoryImpl): RecipeRepository

    @Binds
    @Singleton
    abstract fun bindSessionStore(impl: TokenStore): SessionStore

    @Binds
    @Singleton
    abstract fun bindAuthTokenProvider(impl: TokenStore): AuthTokenProvider
}
