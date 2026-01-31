package handlers

import (
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
)

type WebHandler struct {
	templates *template.Template
}

type PageData struct {
	Title    string
	User     *User
	RecipeID string
}

func NewWebHandler() *WebHandler {
	templates, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		// Templates not found, create empty template for tests
		templates = template.New("")
	}

	return &WebHandler{
		templates: templates,
	}
}

func (h *WebHandler) renderTemplate(w http.ResponseWriter, templateName string, data PageData) {
	// Simple approach: create a new template set each time
	templates := template.Must(template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/header.html",
		"web/templates/footer.html",
		"web/templates/"+templateName,
	))

	err := templates.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *WebHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "RecipeApp - Discover, Create & Share Recipes",
		User:  h.getUserFromContext(r),
	}

	h.renderTemplate(w, "index.html", data)
}

func (h *WebHandler) HandleRecipes(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "All Recipes - RecipeApp",
		User:  h.getUserFromContext(r),
	}

	h.renderTemplate(w, "recipes.html", data)
}

func (h *WebHandler) HandleNewRecipe(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Create New Recipe - RecipeApp",
		User:  h.getUserFromContext(r),
	}

	h.renderTemplate(w, "new-recipe.html", data)
}

func (h *WebHandler) HandleRecipeDetail(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "id")
	data := PageData{
		Title:    "Recipe Detail - RecipeApp",
		User:     h.getUserFromContext(r),
		RecipeID: recipeID,
	}

	h.renderTemplate(w, "recipe-detail.html", data)
}

func (h *WebHandler) getUserFromContext(r *http.Request) *User {
	// This would get user from JWT token in request context
	// For now, return nil (not authenticated)
	return nil
}
