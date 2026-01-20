package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"recipe-app/internal/logger"
	"recipe-app/internal/models"
	"recipe-app/internal/validation"
)

type IngredientHandler struct {
	ingredientService *models.IngredientService
}

func NewIngredientHandler(db *sql.DB) *IngredientHandler {
	return &IngredientHandler{
		ingredientService: models.NewIngredientService(db),
	}
}

func (h *IngredientHandler) HandleIngredients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		h.getIngredients(w, r, ctx)
	case http.MethodPost:
		h.createIngredient(w, r, ctx)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *IngredientHandler) HandleIngredient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		h.getIngredient(w, r, ctx)
	case http.MethodPut:
		h.updateIngredient(w, r, ctx)
	case http.MethodDelete:
		h.deleteIngredient(w, r, ctx)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *IngredientHandler) getIngredients(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	logger.FromContext(r.Context()).Info("Getting ingredients list")

	query := r.URL.Query().Get("search")
	var ingredients []models.Ingredient
	var err error

	if query != "" {
		ingredients, err = h.ingredientService.Search(query)
	} else {
		ingredients, err = h.ingredientService.GetAll()
	}

	if err != nil {
		logger.LogError(r.Context(), err, "Failed to get ingredients")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ingredients)
}

func (h *IngredientHandler) getIngredient(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ingredient ID", http.StatusBadRequest)
		return
	}

	logger.FromContext(r.Context()).Info("Getting ingredient", "id", id)

	ingredient, err := h.ingredientService.GetByID(id)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to get ingredient")
		if err.Error() == "ingredient not found" {
			http.Error(w, "Ingredient not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ingredient)
}

func (h *IngredientHandler) createIngredient(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	var req struct {
		Name     string `json:"name"`
		Category string `json:"category"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate ingredient
	validator := validation.ValidateIngredient(req)
	if validator.HasErrors() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors":  validator.GetErrors(),
			"message": "Validation failed",
		})
		return
	}

	logger.FromContext(r.Context()).Info("Creating new ingredient", "name", req.Name)

	ingredient, err := h.ingredientService.Create(req.Name, req.Category)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to create ingredient")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ingredient)
}

func (h *IngredientHandler) updateIngredient(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ingredient ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Name     string `json:"name"`
		Category string `json:"category"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Name) == 0 {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	logger.FromContext(r.Context()).Info("Updating ingredient", "id", id)

	ingredient, err := h.ingredientService.Update(id, req.Name, req.Category)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to update ingredient")
		if err.Error() == "ingredient not found" {
			http.Error(w, "Ingredient not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ingredient)
}

func (h *IngredientHandler) deleteIngredient(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ingredient ID", http.StatusBadRequest)
		return
	}

	logger.FromContext(r.Context()).Info("Deleting ingredient", "id", id)

	if err := h.ingredientService.Delete(id); err != nil {
		logger.LogError(r.Context(), err, "Failed to delete ingredient")
		if err.Error() == "ingredient not found" {
			http.Error(w, "Ingredient not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Ingredient deleted successfully",
	})
}

type TagHandler struct {
	tagService *models.TagService
}

func NewTagHandler(db *sql.DB) *TagHandler {
	return &TagHandler{
		tagService: models.NewTagService(db),
	}
}

func (h *TagHandler) HandleTags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		h.getTags(w, r, ctx)
	case http.MethodPost:
		h.createTag(w, r, ctx)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TagHandler) HandleTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		h.getTag(w, r, ctx)
	case http.MethodPut:
		h.updateTag(w, r, ctx)
	case http.MethodDelete:
		h.deleteTag(w, r, ctx)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TagHandler) getTags(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	logger.FromContext(r.Context()).Info("Getting tags list")

	tags, err := h.tagService.GetAll()
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to get tags")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

func (h *TagHandler) getTag(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	logger.FromContext(r.Context()).Info("Getting tag", "id", id)

	tag, err := h.tagService.GetByID(id)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to get tag")
		if err.Error() == "tag not found" {
			http.Error(w, "Tag not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tag)
}

func (h *TagHandler) createTag(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	var req struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Name) == 0 {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if len(req.Color) == 0 {
		req.Color = "#3B82F6" // Default blue color
	}

	logger.FromContext(r.Context()).Info("Creating new tag", "name", req.Name)

	tag, err := h.tagService.Create(req.Name, req.Color)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to create tag")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tag)
}

func (h *TagHandler) updateTag(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Name) == 0 {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if len(req.Color) == 0 {
		req.Color = "#3B82F6"
	}

	logger.FromContext(r.Context()).Info("Updating tag", "id", id)

	tag, err := h.tagService.Update(id, req.Name, req.Color)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to update tag")
		if err.Error() == "tag not found" {
			http.Error(w, "Tag not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tag)
}

func (h *TagHandler) deleteTag(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	logger.FromContext(r.Context()).Info("Deleting tag", "id", id)

	if err := h.tagService.Delete(id); err != nil {
		logger.LogError(r.Context(), err, "Failed to delete tag")
		if err.Error() == "tag not found" {
			http.Error(w, "Tag not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Tag deleted successfully",
	})
}
