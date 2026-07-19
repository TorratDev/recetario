package com.recipeapp.network

import javax.inject.Qualifier

/** Bound by the app module (di/AppModule.kt) to the backend's base URL. */
@Qualifier
@Retention(AnnotationRetention.BINARY)
annotation class BaseUrl
