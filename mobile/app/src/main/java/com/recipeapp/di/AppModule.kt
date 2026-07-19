package com.recipeapp.di

import com.recipeapp.BuildConfig
import com.recipeapp.network.BaseUrl
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent

@Module
@InstallIn(SingletonComponent::class)
object AppModule {

    @Provides
    @BaseUrl
    fun provideBaseUrl(): String = BuildConfig.API_BASE_URL
}
