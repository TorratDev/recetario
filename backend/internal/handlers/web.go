package handlers

import (
	"net/http"
)

type WebHandler struct{}

func NewWebHandler() *WebHandler {
	return &WebHandler{}
}

func (h *WebHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/index.html")
}

func (h *WebHandler) HandleRecipes(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/recipes.html")
}

func (h *WebHandler) HandleNewRecipe(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/new-recipe.html")
}