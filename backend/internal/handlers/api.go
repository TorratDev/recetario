package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	appmiddleware "recipe-app/internal/appmiddleware"
	"recipe-app/internal/logger"
	"recipe-app/internal/models"
	"recipe-app/internal/validation"
)

type APIHandler struct {
	recipeService *models.RecipeService
}

func NewAPIHandler(db *sql.DB) *APIHandler {
	return &APIHandler{
		recipeService: models.NewRecipeService(db),
	}
}

func (h *APIHandler) HandleRecipes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		h.getRecipes(w, r, ctx)
	case http.MethodPost:
		h.createRecipe(w, r, ctx)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *APIHandler) HandleRecipe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		h.getRecipe(w, r, ctx)
	case http.MethodPut:
		h.updateRecipe(w, r, ctx)
	case http.MethodDelete:
		h.deleteRecipe(w, r, ctx)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *APIHandler) HandleCreateRecipe(w http.ResponseWriter, r *http.Request) {
	h.createRecipe(w, r, r.Context())
}

func (h *APIHandler) HandleUpdateRecipe(w http.ResponseWriter, r *http.Request) {
	h.updateRecipe(w, r, r.Context())
}

func (h *APIHandler) HandleDeleteRecipe(w http.ResponseWriter, r *http.Request) {
	h.deleteRecipe(w, r, r.Context())
}

func (h *APIHandler) getRecipes(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	logger.FromContext(r.Context()).Info("Getting recipes list")

	// Parse query parameters
	filter := models.RecipeFilter{}

	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		if userID, err := strconv.Atoi(userIDStr); err == nil {
			filter.UserID = &userID
		}
	}

	if title := r.URL.Query().Get("title"); title != "" {
		filter.Title = &title
	}

	if difficulty := r.URL.Query().Get("difficulty"); difficulty != "" {
		filter.Difficulty = &difficulty
	}

	if isPublicStr := r.URL.Query().Get("is_public"); isPublicStr != "" {
		if isPublic, err := strconv.ParseBool(isPublicStr); err == nil {
			filter.IsPublic = &isPublic
		}
	}

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	filter.Limit = &limit

	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}
	filter.Offset = &offset

	if sort := r.URL.Query().Get("sort_by"); sort != "" {
		filter.SortBy = sort
	}

	if order := r.URL.Query().Get("sort_order"); order != "" {
		filter.SortOrder = order
	}

	recipes, err := h.recipeService.GetAll(filter)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to get recipes")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func (h *APIHandler) getRecipe(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	logger.FromContext(r.Context()).Info("Getting recipe", "id", id)

	recipe, err := h.recipeService.GetByID(id, true)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to get recipe")
		if err.Error() == "recipe not found" {
			http.Error(w, "Recipe not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

func (h *APIHandler) createRecipe(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	var recipe models.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate recipe
	validator := validation.ValidateRecipe(recipe)
	if validator.HasErrors() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors":  validator.GetErrors(),
			"message": "Validation failed",
		})
		return
	}

	// Get user ID from context
	userID, ok := appmiddleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	recipe.UserID = userID

	logger.FromContext(r.Context()).Info("Creating new recipe", "title", recipe.Title)

	createdRecipe, err := h.recipeService.Create(&recipe)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to create recipe")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Recipe created successfully",
		"recipe":  createdRecipe,
	})
}

func (h *APIHandler) updateRecipe(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	var recipe models.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate recipe
	validator := validation.ValidateRecipe(recipe)
	if validator.HasErrors() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors":  validator.GetErrors(),
			"message": "Validation failed",
		})
		return
	}

	logger.FromContext(r.Context()).Info("Updating recipe", "id", id)

	updatedRecipe, err := h.recipeService.Update(id, &recipe)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to update recipe")
		if err.Error() == "recipe not found" {
			http.Error(w, "Recipe not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Recipe updated successfully",
		"recipe":  updatedRecipe,
	})
}

func (h *APIHandler) deleteRecipe(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	logger.FromContext(r.Context()).Info("Deleting recipe", "id", id)

	if err := h.recipeService.Delete(id); err != nil {
		logger.LogError(r.Context(), err, "Failed to delete recipe")
		if err.Error() == "recipe not found" {
			http.Error(w, "Recipe not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Recipe deleted successfully",
	})
}
