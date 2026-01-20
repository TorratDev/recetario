package handlers

import (
	"encoding/json"
	"net/http"

	"recipe-app/internal/appmiddleware"
	"recipe-app/internal/logger"
)

type UserHandler struct{}

type ProfileUpdateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := appmiddleware.GetUserID(ctx)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	user := User{
		ID:    userID,
		Email: "user@example.com",
		Name:  "Test User",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := appmiddleware.GetUserID(ctx)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	var req ProfileUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	logger.FromContext(ctx).Info("Profile update requested", "user_id", userID, "name", req.Name)

	user := User{
		ID:    userID,
		Email: req.Email,
		Name:  req.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
