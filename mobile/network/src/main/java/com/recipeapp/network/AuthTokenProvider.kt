package com.recipeapp.network

/**
 * Port implemented by the data module's token store (data/.../auth/TokenStore.kt)
 * so AuthInterceptor can read the current bearer token without the network
 * module depending on data (dependency direction stays network <- data <- app).
 */
interface AuthTokenProvider {
    fun currentToken(): String?
}
