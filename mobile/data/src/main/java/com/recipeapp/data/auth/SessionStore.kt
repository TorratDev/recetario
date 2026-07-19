package com.recipeapp.data.auth

import com.recipeapp.domain.model.AuthSession
import com.recipeapp.network.AuthTokenProvider
import kotlinx.coroutines.flow.StateFlow

/**
 * Split out from TokenStore so AuthRepositoryImpl can be unit tested with a
 * fake — TokenStore itself needs an Android Context (EncryptedSharedPreferences)
 * and can't be constructed in a plain JUnit test.
 */
interface SessionStore : AuthTokenProvider {
    val session: StateFlow<AuthSession?>
    fun save(authSession: AuthSession)
    fun updateToken(token: String, expiresIn: Long)
    fun clear()
}
