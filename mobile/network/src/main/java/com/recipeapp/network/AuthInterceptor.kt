package com.recipeapp.network

import okhttp3.Interceptor
import okhttp3.Response
import javax.inject.Inject

/**
 * Attaches `Authorization: Bearer <token>` to every request when a session
 * exists. The backend checks the Bearer header before any cookie
 * (appmiddleware/auth.go), so this is the only auth transport the Android
 * client needs — no cookie jar.
 */
class AuthInterceptor @Inject constructor(
    private val tokenProvider: AuthTokenProvider
) : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val original = chain.request()
        val token = tokenProvider.currentToken()
        val request = if (token != null) {
            original.newBuilder()
                .header("Authorization", "Bearer $token")
                .build()
        } else {
            original
        }
        return chain.proceed(request)
    }
}
