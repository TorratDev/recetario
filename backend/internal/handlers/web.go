package handlers

import (
	"html"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"recipe-app/internal/appmiddleware"
	"recipe-app/internal/logger"
	"recipe-app/internal/models"
)

type WebHandler struct {
	templates *template.Template
	users     UserStore
	recipes   RecipeStore
}

type PageData struct {
	Title    string
	User     *User
	RecipeID string
}

func NewWebHandler(users UserStore, recipes RecipeStore) *WebHandler {
	templates, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		// Templates not found, create empty template for tests
		templates = template.New("")
	}

	return &WebHandler{
		templates: templates,
		users:     users,
		recipes:   recipes,
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
	user := h.getUserFromContext(r)
	if user == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := PageData{
		Title: "Create New Recipe - RecipeApp",
		User:  user,
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

func (h *WebHandler) HandleEditRecipe(w http.ResponseWriter, r *http.Request) {
	user := h.getUserFromContext(r)
	if user == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	recipeID := chi.URLParam(r, "id")
	data := PageData{
		Title:    "Edit Recipe - RecipeApp",
		User:     user,
		RecipeID: recipeID,
	}

	h.renderTemplate(w, "edit-recipe.html", data)
}

func (h *WebHandler) HandleCollections(w http.ResponseWriter, r *http.Request) {
	user := h.getUserFromContext(r)
	if user == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := PageData{
		Title: "My Collections - RecipeApp",
		User:  user,
	}

	h.renderTemplate(w, "collections.html", data)
}

func (h *WebHandler) HandleRecipeSearch(w http.ResponseWriter, r *http.Request) {
	// Extract query parameters for search
	q := r.URL.Query()
	filter := &models.RecipeFilter{
		Query:      q.Get("q"),
		Difficulty: q.Get("difficulty"),
		MinCookTime: parseIntDefault(q.Get("min_cook_time"), 0),
		MaxCookTime: parseIntDefault(q.Get("max_cook_time"), 0),
	}

	limit := parseIntDefault(q.Get("limit"), 20)
	offset := parseIntDefault(q.Get("offset"), 0)

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// Search recipes
	recipes, _, err := h.recipes.SearchRecipes(filter, limit, offset)
	if err != nil {
		http.Error(w, "Error searching recipes", http.StatusInternalServerError)
		return
	}

	// Return HTML template of recipe cards
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Simple HTML rendering
	if len(recipes) == 0 {
		w.Write([]byte(`<div style="grid-column:1/-1; text-align:center; padding:3rem; color: #6b7280;"><p>No se encontraron recetas</p></div>`))
		return
	}

	var html strings.Builder
	for _, recipe := range recipes {
		if recipe == nil {
			continue
		}
		html.WriteString(h.recipeCardHTML(recipe))
	}

	w.Write([]byte(html.String()))
}

func (h *WebHandler) recipeCardHTML(recipe *models.Recipe) string {
	var html strings.Builder
	html.WriteString(`<div class="recipe-card">`)

	// Image
	if recipe.ImageURL != "" {
		html.WriteString(`<img src="` + escapeAttr(recipe.ImageURL) + `" alt="` + escapeAttr(recipe.Title) + `" class="recipe-image">`)
	} else {
		html.WriteString(`<div class="recipe-image-placeholder">No image</div>`)
	}

	html.WriteString(`<div class="recipe-content">`)
	html.WriteString(`<a href="/recipes/` + recipe.ID + `" class="recipe-title">` + escapeHTML(recipe.Title) + `</a>`)

	if recipe.Description != "" {
		html.WriteString(`<p class="recipe-description">` + escapeHTML(recipe.Description) + `</p>`)
	}

	html.WriteString(`<div class="recipe-meta">`)
	html.WriteString(`<span class="meta-item">⏱️ ` + strconv.Itoa(recipe.CookTime) + `min</span>`)
	html.WriteString(`<span class="meta-item">👥 ` + strconv.Itoa(recipe.Servings) + `</span>`)
	html.WriteString(`<span class="meta-item difficulty-` + recipe.Difficulty + `">` + recipe.Difficulty + `</span>`)
	html.WriteString(`</div>`)

	html.WriteString(`</div>`)
	html.WriteString(`</div>`)

	return html.String()
}

func (h *WebHandler) getUserFromContext(r *http.Request) *User {
	claims, ok := appmiddleware.GetUserClaims(r.Context())
	if !ok || claims == nil || claims.UserID == "" || h.users == nil {
		return nil
	}

	u, err := h.users.GetUserByID(claims.UserID)
	if err != nil {
		logger.FromContext(r.Context()).Error("Failed to load user from context", "error", err)
		return nil
	}

	return &User{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func escapeHTML(s string) string {
	return html.EscapeString(s)
}

func escapeAttr(s string) string {
	return html.EscapeString(s)
}
