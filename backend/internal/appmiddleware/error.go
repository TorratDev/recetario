package appmiddleware

import (
	"encoding/json"
	"net/http"
	"time"

	"recipe-app/internal/logger"

	"github.com/go-chi/chi/v5/middleware"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

type AppError struct {
	StatusCode int
	Message    string
	Code       string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func NewAppError(statusCode int, message, code string, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
		Code:       code,
		Err:        err,
	}
}

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				ctx := r.Context()
				logger.LogError(ctx, nil, "Panic recovered")

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				json.NewEncoder(w).Encode(ErrorResponse{
					Error:   "Internal server error",
					Message: "An unexpected error occurred",
					Code:    "INTERNAL_ERROR",
				})
			}
		}()

		next.ServeHTTP(&errorResponseWriter{ResponseWriter: w, r: r}, r)
	})
}

type errorResponseWriter struct {
	http.ResponseWriter
	r *http.Request
}

func (erw *errorResponseWriter) WriteHeader(statusCode int) {
	if statusCode >= 400 {
		ctx := erw.r.Context()
		logger.LogError(ctx, nil, "HTTP error response")

		erw.ResponseWriter.Header().Set("Content-Type", "application/json")
		erw.ResponseWriter.WriteHeader(statusCode)

		errorResp := ErrorResponse{
			Error:   http.StatusText(statusCode),
			Message: getErrorMessage(statusCode),
			Code:    getErrorCode(statusCode),
		}

		json.NewEncoder(erw.ResponseWriter).Encode(errorResp)
		return
	}

	erw.ResponseWriter.WriteHeader(statusCode)
}

func getErrorMessage(statusCode int) string {
	messages := map[int]string{
		http.StatusBadRequest:          "Invalid request parameters",
		http.StatusUnauthorized:        "Authentication required",
		http.StatusForbidden:           "Access forbidden",
		http.StatusNotFound:            "Resource not found",
		http.StatusMethodNotAllowed:    "Method not allowed",
		http.StatusTooManyRequests:     "Rate limit exceeded",
		http.StatusInternalServerError: "Internal server error",
		http.StatusBadGateway:          "Service unavailable",
		http.StatusServiceUnavailable:  "Service temporarily unavailable",
	}

	if msg, ok := messages[statusCode]; ok {
		return msg
	}
	return "An error occurred"
}

func getErrorCode(statusCode int) string {
	codes := map[int]string{
		http.StatusBadRequest:          "BAD_REQUEST",
		http.StatusUnauthorized:        "UNAUTHORIZED",
		http.StatusForbidden:           "FORBIDDEN",
		http.StatusNotFound:            "NOT_FOUND",
		http.StatusMethodNotAllowed:    "METHOD_NOT_ALLOWED",
		http.StatusTooManyRequests:     "RATE_LIMIT_EXCEEDED",
		http.StatusInternalServerError: "INTERNAL_ERROR",
		http.StatusBadGateway:          "BAD_GATEWAY",
		http.StatusServiceUnavailable:  "SERVICE_UNAVAILABLE",
	}

	if code, ok := codes[statusCode]; ok {
		return code
	}
	return "UNKNOWN_ERROR"
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ctx := logger.WithCorrelationID(r.Context())
		r = r.WithContext(ctx)

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		logger.LogRequest(ctx, r.Method, r.RequestURI, duration, ww.Status())
	})
}
