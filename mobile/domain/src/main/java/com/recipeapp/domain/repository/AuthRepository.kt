package com.recipeapp.domain.repository

import com.recipeapp.domain.common.UiResult
import com.recipeapp.domain.model.AuthSession
import com.recipeapp.domain.model.LoginCredentials
import com.recipeapp.domain.model.RegisterDetails
import kotlinx.coroutines.flow.Flow

interface AuthRepository {
    /** Emits the current session, or null when signed out. */
    val session: Flow<AuthSession?>

    suspend fun login(credentials: LoginCredentials): UiResult<AuthSession>

    suspend fun register(details: RegisterDetails): UiResult<AuthSession>

    /**
     * Re-signs a new token from the current one. The backend has no separate
     * long-lived refresh token — if the current token has already expired,
     * this also fails, and the caller must clear the session and force a
     * full re-login.
     */
    suspend fun refresh(): UiResult<Unit>

    /** Clears the local session regardless of whether the network call succeeds. */
    suspend fun logout()
}
