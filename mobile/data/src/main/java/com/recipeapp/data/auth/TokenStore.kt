package com.recipeapp.data.auth

import android.content.Context
import androidx.security.crypto.EncryptedSharedPreferences
import androidx.security.crypto.MasterKey
import com.google.gson.Gson
import com.recipeapp.domain.model.AuthSession
import com.recipeapp.network.AuthTokenProvider
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Persists the auth session in EncryptedSharedPreferences (Android Keystore
 * backed) rather than plain DataStore, since this specifically stores a
 * bearer token/credential, not general app preferences.
 */
@Singleton
class TokenStore @Inject constructor(
    @ApplicationContext context: Context,
    private val gson: Gson
) : SessionStore {

    private val prefs = EncryptedSharedPreferences.create(
        context,
        PREFS_NAME,
        MasterKey.Builder(context).setKeyScheme(MasterKey.KeyScheme.AES256_GCM).build(),
        EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
        EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
    )

    private val _session = MutableStateFlow(readSession())
    override val session: StateFlow<AuthSession?> = _session

    override fun currentToken(): String? = _session.value?.token

    override fun save(authSession: AuthSession) {
        prefs.edit().putString(KEY_SESSION, gson.toJson(authSession)).apply()
        _session.value = authSession
    }

    override fun updateToken(token: String, expiresIn: Long) {
        val updated = _session.value?.copy(token = token, expiresIn = expiresIn) ?: return
        save(updated)
    }

    override fun clear() {
        prefs.edit().remove(KEY_SESSION).apply()
        _session.value = null
    }

    private fun readSession(): AuthSession? {
        val json = prefs.getString(KEY_SESSION, null) ?: return null
        return try {
            gson.fromJson(json, AuthSession::class.java)
        } catch (e: Exception) {
            null
        }
    }

    private companion object {
        const val PREFS_NAME = "recipeapp_secure_prefs"
        const val KEY_SESSION = "session"
    }
}
