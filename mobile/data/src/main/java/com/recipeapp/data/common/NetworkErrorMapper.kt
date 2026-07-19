package com.recipeapp.data.common

import com.google.gson.Gson
import com.recipeapp.domain.common.UiResult
import com.recipeapp.network.dto.ErrorResponseDto
import retrofit2.HttpException
import java.io.IOException

/** Shared by every repository that talks to the backend's single error envelope. */
fun Throwable.toUiError(gson: Gson): UiResult.Error {
    if (this is HttpException) {
        val message = response()?.errorBody()?.string()?.let { body ->
            runCatching { gson.fromJson(body, ErrorResponseDto::class.java) }.getOrNull()
        }?.let { it.userMessage ?: it.message ?: it.error }
        return UiResult.Error(
            message = message ?: "Request failed (${code()})",
            isUnauthorized = code() == 401
        )
    }
    if (this is IOException) {
        return UiResult.Error(message = "Network error, check your connection")
    }
    return UiResult.Error(message = message ?: "Unknown error")
}
