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

type CollectionHandler struct {
	collectionService *models.CollectionService
}

func NewCollectionHandler(db *sql.DB) *CollectionHandler {
	return &CollectionHandler{
		collectionService: models.NewCollectionService(db),
	}
}

func (h *CollectionHandler) HandleCollections(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		h.getCollections(w, r, ctx)
	case http.MethodPost:
		h.createCollection(w, r, ctx)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CollectionHandler) HandleCollection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		h.getCollection(w, r, ctx)
	case http.MethodPut:
		h.updateCollection(w, r, ctx)
	case http.MethodDelete:
		h.deleteCollection(w, r, ctx)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CollectionHandler) HandleCollectionRecipes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		h.getCollectionRecipes(w, r, ctx)
	case http.MethodPost:
		h.addRecipeToCollection(w, r, ctx)
	case http.MethodDelete:
		h.removeRecipeFromCollection(w, r, ctx)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CollectionHandler) getCollections(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	uid, ok := appmiddleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	logger.FromContext(r.Context()).Info("Getting collections", "user_id", userID)

	collections, err := h.collectionService.GetAll(userID)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to get collections")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collections)
}

func (h *CollectionHandler) getCollection(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid collection ID", http.StatusBadRequest)
		return
	}

	uid, ok := appmiddleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	includeRecipes := r.URL.Query().Get("include_recipes") == "true"

	logger.FromContext(r.Context()).Info("Getting collection", "id", id, "include_recipes", includeRecipes)

	collection, err := h.collectionService.GetByID(id, userID, includeRecipes)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to get collection")
		if err.Error() == "collection not found" {
			http.Error(w, "Collection not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collection)
}

func (h *CollectionHandler) createCollection(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	uid, ok := appmiddleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req models.CreateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	validator := validation.NewValidator()
	validator.Required("name", req.Name)
	validator.MinLength("name", req.Name, 2)
	validator.MaxLength("name", req.Name, 255)

	if validator.HasErrors() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors":  validator.GetErrors(),
			"message": "Validation failed",
		})
		return
	}

	logger.FromContext(r.Context()).Info("Creating collection", "name", req.Name, "user_id", userID)

	collection, err := h.collectionService.Create(userID, req)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to create collection")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(collection)
}

func (h *CollectionHandler) updateCollection(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid collection ID", http.StatusBadRequest)
		return
	}

	uid, ok := appmiddleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req models.UpdateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	validator := validation.NewValidator()
	if req.Name != nil {
		validator.MinLength("name", *req.Name, 2)
		validator.MaxLength("name", *req.Name, 255)
	}
	if req.Description != nil {
		validator.MaxLength("description", *req.Description, 1000)
	}

	if validator.HasErrors() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors":  validator.GetErrors(),
			"message": "Validation failed",
		})
		return
	}

	logger.FromContext(r.Context()).Info("Updating collection", "id", id, "user_id", uid)

	collection, err := h.collectionService.Update(id, uid, req)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to update collection")
		if err.Error() == "collection not found" {
			http.Error(w, "Collection not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collection)
}

func (h *CollectionHandler) deleteCollection(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid collection ID", http.StatusBadRequest)
		return
	}

	uid, ok := appmiddleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	logger.FromContext(r.Context()).Info("Deleting collection", "id", id, "user_id", uid)

	if err := h.collectionService.Delete(id, uid); err != nil {
		logger.LogError(r.Context(), err, "Failed to delete collection")
		if err.Error() == "collection not found" {
			http.Error(w, "Collection not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Collection deleted successfully",
	})
}

func (h *CollectionHandler) getCollectionRecipes(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	idStr := chi.URLParam(r, "id")
	collectionID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid collection ID", http.StatusBadRequest)
		return
	}

	uid, ok := appmiddleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	logger.FromContext(r.Context()).Info("Getting collection recipes", "collection_id", collectionID)

	collection, err := h.collectionService.GetByID(collectionID, uid, true)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to get collection recipes")
		if err.Error() == "collection not found" {
			http.Error(w, "Collection not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collection)
}

func (h *CollectionHandler) addRecipeToCollection(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	uid, ok := appmiddleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req struct {
		RecipeID     int `json:"recipe_id"`
		CollectionID int `json:"collection_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	validator := validation.NewValidator()
	validator.Required("recipe_id", strconv.Itoa(req.RecipeID))
	validator.Required("collection_id", strconv.Itoa(req.CollectionID))

	if validator.HasErrors() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors":  validator.GetErrors(),
			"message": "Validation failed",
		})
		return
	}

	logger.FromContext(r.Context()).Info("Adding recipe to collection", "recipe_id", req.RecipeID, "collection_id", req.CollectionID)

	if err := h.collectionService.AddRecipe(req.CollectionID, req.RecipeID); err != nil {
		logger.LogError(r.Context(), err, "Failed to add recipe to collection")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Recipe added to collection successfully",
	})
}

func (h *CollectionHandler) removeRecipeFromCollection(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	idStr := chi.URLParam(r, "id")
	collectionID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid collection ID", http.StatusBadRequest)
		return
	}

	uid, ok := appmiddleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req struct {
		RecipeID int `json:"recipe_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	validator := validation.NewValidator()
	validator.Required("recipe_id", strconv.Itoa(req.RecipeID))

	if validator.HasErrors() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors":  validator.GetErrors(),
			"message": "Validation failed",
		})
		return
	}

	logger.FromContext(r.Context()).Info("Removing recipe from collection", "collection_id", collectionID, "recipe_id", req.RecipeID)

	if err := h.collectionService.RemoveRecipe(collectionID, req.RecipeID); err != nil {
		logger.LogError(r.Context(), err, "Failed to remove recipe from collection")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Recipe removed from collection successfully",
	})
}
