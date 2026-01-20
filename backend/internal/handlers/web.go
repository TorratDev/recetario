package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"recipe-app/internal/logger"
	"recipe-app/internal/models"
)

type WebHandler struct {
	db            *sql.DB
	recipeService *models.RecipeService
	templates     *template.Template
}

func NewWebHandler(db *sql.DB) *WebHandler {
	// Parse templates
	templates := template.Must(template.ParseGlob("web/templates/*.html"))

	return &WebHandler{
		db:            db,
		recipeService: models.NewRecipeService(db),
		templates:     templates,
	}
}

func (h *WebHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	logger.FromContext(r.Context()).Info("Serving index page")

	if err := h.templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		logger.LogError(r.Context(), err, "Failed to render index template")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *WebHandler) HandleRecipes(w http.ResponseWriter, r *http.Request) {
	logger.FromContext(r.Context()).Info("Serving recipes page")

	if err := h.templates.ExecuteTemplate(w, "recipes.html", nil); err != nil {
		logger.LogError(r.Context(), err, "Failed to render recipes template")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *WebHandler) HandleNewRecipe(w http.ResponseWriter, r *http.Request) {
	logger.FromContext(r.Context()).Info("Serving new recipe page")

	// TODO: Create new-recipe.html template
	http.Error(w, "New recipe form not implemented yet", http.StatusNotImplemented)
}

func (h *WebHandler) HandleRecipeList(w http.ResponseWriter, r *http.Request) {
	logger.FromContext(r.Context()).Info("Serving recipe list via HTMX")

	// Parse query parameters (same as API)
	filter := models.RecipeFilter{}

	// For now, get default recipes
	recipes, err := h.recipeService.GetAll(filter)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to get recipes")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Render recipe cards
	tmpl := `{{range .}}
<div class="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-300">
    <div class="h-48 bg-gray-200">
        {{if .ImageURL}}
        <img src="{{.ImageURL}}" alt="{{.Title}}" class="w-full h-48 object-cover">
        {{else}}
        <div class="w-full h-48 bg-gray-200 flex items-center justify-center">
            <i class="fas fa-utensils text-4xl text-gray-400"></i>
        </div>
        {{end}}
    </div>
    
    <div class="p-4">
        <div class="flex justify-between items-start mb-2">
            <h3 class="text-lg font-semibold text-gray-900 truncate flex-1">{{.Title}}</h3>
            {{if .Difficulty}}
            <span class="ml-2 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                {{.Difficulty}}
            </span>
            {{end}}
        </div>
        
        {{if .Description}}
        <p class="text-gray-600 text-sm mb-3">{{.Description}}</p>
        {{end}}
        
        <div class="flex items-center justify-between text-sm text-gray-500 mb-3">
            {{if .PrepTime}}
            <span class="flex items-center">
                <i class="fas fa-clock mr-1"></i>{{.PrepTime}} min
            </span>
            {{end}}
            <span class="flex items-center">
                <i class="fas fa-users mr-1"></i>{{.Servings}} servings
            </span>
        </div>
        
        <div class="flex gap-2">
            <a href="/recipes/{{.ID}}" 
               class="flex-1 text-center bg-blue-600 text-white px-3 py-2 rounded-md hover:bg-blue-700 text-sm">
                <i class="fas fa-eye mr-1"></i>View
            </a>
        </div>
    </div>
</div>
{{else}}
<div class="col-span-full text-center py-12">
    <i class="fas fa-search text-6xl text-gray-300 mb-4"></i>
    <h3 class="text-xl font-semibold text-gray-600 mb-2">No recipes found</h3>
    <p class="text-gray-500">Try adjusting your search or filters</p>
</div>
{{end}}`

	t, err := template.New("recipes").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := t.Execute(w, recipes); err != nil {
		logger.LogError(r.Context(), err, "Failed to render recipe cards")
	}
}
