package handlers

import (
	"encoding/json"
	"net/http"

	"recipe-app/internal/logger"
)

type APIHandler struct{}

func NewAPIHandler() *APIHandler {
	return &APIHandler{}
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
	w.Header().Set("Content-Type", "application/json")
	recipes := []map[string]interface{}{
		{"id": "1", "name": "Spaghetti Bolognese", "prep_time": 30, "difficulty": "medium"},
		{"id": "2", "name": "Chicken Curry", "prep_time": 45, "difficulty": "hard"},
	}
	json.NewEncoder(w).Encode(recipes)
}

func (h *APIHandler) createRecipe(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	logger.FromContext(r.Context()).Info("Creating new recipe")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Recipe created successfully",
		"id":      "3",
	})
}

func (h *APIHandler) getRecipe(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	w.Header().Set("Content-Type", "application/json")
	recipe := map[string]interface{}{
		"id":           "1",
		"name":         "Spaghetti Bolognese",
		"prep_time":    30,
		"cook_time":    45,
		"difficulty":   "medium",
		"ingredients":  []string{"spaghetti", "ground beef", "tomato sauce", "onions"},
		"instructions": "1. Cook pasta. 2. Brown beef. 3. Add sauce. 4. Combine and serve.",
	}
	json.NewEncoder(w).Encode(recipe)
}

func (h *APIHandler) updateRecipe(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	logger.FromContext(r.Context()).Info("Updating recipe")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Recipe updated successfully",
	})
}

func (h *APIHandler) deleteRecipe(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	logger.FromContext(r.Context()).Info("Deleting recipe")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Recipe deleted successfully",
	})
}
