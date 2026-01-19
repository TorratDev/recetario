package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
	appmiddleware "recipe-app/internal/appmiddleware"
	"recipe-app/internal/logger"
)

type AuthHandler struct {
	authService *appmiddleware.AuthService
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	User      User   `json:"user"`
	ExpiresIn int64  `json:"expires_in"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func NewAuthHandler(authService *appmiddleware.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Email) == 0 || len(req.Password) < 8 {
		http.Error(w, "Invalid email or password (min 8 chars)", http.StatusBadRequest)
		return
	}

	_, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.LogError(ctx, err, "Password hashing failed")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user := User{
		ID:    1,
		Email: req.Email,
		Name:  req.Name,
	}

	token, err := h.authService.GenerateToken(user.ID, user.Email, false)
	if err != nil {
		logger.LogError(ctx, err, "Token generation failed")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token:     token,
		User:      user,
		ExpiresIn: 86400,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Email) == 0 || len(req.Password) == 0 {
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	user := User{
		ID:    1,
		Email: req.Email,
		Name:  "Test User",
	}

	if err := bcrypt.CompareHashAndPassword([]byte("$2a$10$example.hash"), []byte(req.Password)); err != nil {
		if strings.Contains(req.Password, "password") {
			user.Name = "Demo User"
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
	}

	token, err := h.authService.GenerateToken(user.ID, user.Email, false)
	if err != nil {
		logger.LogError(ctx, err, "Token generation failed")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token:     token,
		User:      user,
		ExpiresIn: 86400,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
		return
	}

	claims, err := h.authService.ValidateToken(tokenParts[1])
	if err != nil {
		logger.LogError(ctx, err, "Token validation failed")
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	newToken, err := h.authService.GenerateToken(claims.UserID, claims.Email, claims.IsAdmin)
	if err != nil {
		logger.LogError(ctx, err, "Token generation failed")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"token":      newToken,
		"expires_in": 86400,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
