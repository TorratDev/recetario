package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"recipe-app/internal/logger"
)

type APIHandler struct {
	templates *template.Template
}

func NewAPIHandler() *APIHandler {
	templates, err := template.ParseFiles("web/templates/recipe-cards.html", "web/templates/recipe-detail-content.html")
	if err != nil {
		// Templates not found, create empty template for tests
		templates = template.New("")
	}
	return &APIHandler{
		templates: templates,
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
	recipes := []map[string]interface{}{
		{"id": "1", "title": "Spaghetti Bolognese", "description": "Classic Italian pasta dish with rich meat sauce", "cook_time": 30, "difficulty": "medium"},
		{"id": "2", "title": "Chicken Curry", "description": "Spicy and aromatic Indian curry with tender chicken", "cook_time": 45, "difficulty": "hard"},
		{"id": "3", "title": "Caesar Salad", "description": "Fresh romaine lettuce with creamy Caesar dressing", "cook_time": 15, "difficulty": "easy"},
		{"id": "4", "title": "Beef Tacos", "description": "Mexican-style tacos with seasoned ground beef", "cook_time": 25, "difficulty": "medium"},
		{"id": "5", "title": "Chocolate Cake", "description": "Rich and moist chocolate cake with fudge frosting", "cook_time": 60, "difficulty": "hard"},
		{"id": "6", "title": "Greek Salad", "description": "Mediterranean salad with feta cheese and olives", "cook_time": 10, "difficulty": "easy"},
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		tmpl := h.templates.Lookup("recipe-cards.html")
		if tmpl != nil {
			data := map[string]interface{}{"recipes": recipes}
			err := tmpl.Execute(w, data)
			if err != nil {
				http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
			}
			return
		} else {
			// Debug: template not found
			http.Error(w, "Template recipe-cards.html not found", http.StatusInternalServerError)
			return
		}
	}

	// Default JSON response
	w.Header().Set("Content-Type", "application/json")
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
	recipe := map[string]interface{}{
		"id":          "1",
		"title":       "Spaghetti Bolognese",
		"description": "Classic Italian pasta dish with rich meat sauce",
		"prep_time":   30,
		"cook_time":   45,
		"servings":    4,
		"difficulty":  "medium",
		"ingredients": []map[string]interface{}{
			{"name": "spaghetti", "quantity": "400g", "unit": "grams"},
			{"name": "ground beef", "quantity": "500g", "unit": "grams"},
			{"name": "tomato sauce", "quantity": "800ml", "unit": "milliliters"},
			{"name": "onion", "quantity": "1", "unit": "large"},
			{"name": "garlic", "quantity": "3", "unit": "cloves"},
			{"name": "olive oil", "quantity": "2", "unit": "tablespoons"},
		},
		"instructions": []string{
			"Bring a large pot of salted water to boil and cook spaghetti according to package directions.",
			"Heat olive oil in a large pan over medium heat. Add chopped onion and cook until translucent.",
			"Add minced garlic and cook for another minute until fragrant.",
			"Add ground beef and cook until browned, breaking it up with a wooden spoon.",
			"Pour in tomato sauce and simmer for 15-20 minutes, stirring occasionally.",
			"Season with salt, pepper, and Italian herbs to taste.",
			"Drain pasta and toss with the bolognese sauce. Serve hot with grated Parmesan cheese.",
		},
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		tmpl := h.templates.Lookup("recipe-detail-content.html")
		if tmpl != nil {
			data := map[string]interface{}{"recipe": recipe}
			tmpl.Execute(w, data)
			return
		}
	}

	// Default JSON response
	w.Header().Set("Content-Type", "application/json")
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
