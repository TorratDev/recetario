package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"recipe-app/internal/logger"
	"recipe-app/internal/models"
)

type SearchHandler struct {
	searchService *models.SearchService
}

func NewSearchHandler(db *sql.DB) *SearchHandler {
	return &SearchHandler{
		searchService: models.NewSearchService(db),
	}
}

func (h *SearchHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		h.searchRecipes(w, r, ctx)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *SearchHandler) HandleSuggestions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		h.getSuggestions(w, r, ctx)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *SearchHandler) HandlePopularTags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		h.getPopularTags(w, r, ctx)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *SearchHandler) searchRecipes(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	logger.FromContext(r.Context()).Info("Processing recipe search")

	// Parse query parameters
	filters := models.SearchFilters{
		Query:      r.URL.Query().Get("q"),
		UserID:     h.parseOptionalInt(r.URL.Query().Get("user_id")),
		Difficulty: r.URL.Query().Get("difficulty"),
		SortBy:     r.URL.Query().Get("sort_by"),
		SortOrder:  r.URL.Query().Get("sort_order"),
		Limit:      h.parseOptionalInt(r.URL.Query().Get("limit")),
		Offset:     h.parseOptionalInt(r.URL.Query().Get("offset")),
	}

	// Parse tags array
	if tagsParam := r.URL.Query().Get("tags"); tagsParam != "" {
		// Simple parsing of comma-separated tags
		// In production, you might want more sophisticated parsing
		filters.Tags = []string{tagsParam}
	}

	// Parse categories array
	if categoriesParam := r.URL.Query().Get("categories"); categoriesParam != "" {
		// Simple parsing
		cats := h.parseIntList(categoriesParam)
		filters.Categories = cats
	}

	// Parse time ranges
	filters.MinPrepTime = h.parseOptionalInt(r.URL.Query().Get("min_prep_time"))
	filters.MaxPrepTime = h.parseOptionalInt(r.URL.Query().Get("max_prep_time"))
	filters.MinCookTime = h.parseOptionalInt(r.URL.Query().Get("min_cook_time"))
	filters.MaxCookTime = h.parseOptionalInt(r.URL.Query().Get("max_cook_time"))

	// Parse servings range
	filters.MinServings = h.parseOptionalInt(r.URL.Query().Get("min_servings"))
	filters.MaxServings = h.parseOptionalInt(r.URL.Query().Get("max_servings"))

	// Parse public/private filter
	if isPublicStr := r.URL.Query().Get("is_public"); isPublicStr != "" {
		isPublic, err := strconv.ParseBool(isPublicStr)
		if err == nil {
			filters.IsPublic = &isPublic
		}
	}

	// Set default values
	if filters.SortBy == "" {
		filters.SortBy = "created_at"
	}
	if filters.SortOrder == "" {
		filters.SortOrder = "DESC"
	}
	if filters.Limit == nil {
		limit := 20
		filters.Limit = &limit
	}

	// Perform search
	result, err := h.searchService.SearchRecipes(filters)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to search recipes")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *SearchHandler) getSuggestions(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	query := r.URL.Query().Get("q")
	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]string{})
		return
	}

	logger.FromContext(r.Context()).Info("Getting search suggestions", "query", query)

	suggestions, err := h.searchService.GetSearchSuggestions(query, 10)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to get suggestions")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

func (h *SearchHandler) getPopularTags(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	limit := h.parseOptionalInt(r.URL.Query().Get("limit"))
	if limit == nil {
		limitVal := 20
		limit = &limitVal
	}

	logger.FromContext(r.Context()).Info("Getting popular tags", "limit", *limit)

	tags, err := h.searchService.GetPopularTags(*limit)
	if err != nil {
		logger.LogError(r.Context(), err, "Failed to get popular tags")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

// Helper methods
func (h *SearchHandler) parseOptionalInt(s string) *int {
	if s == "" {
		return nil
	}
	if i, err := strconv.Atoi(s); err == nil {
		return &i
	}
	return nil
}

func (h *SearchHandler) parseIntList(s string) []int {
	if s == "" {
		return []int{}
	}

	var result []int
	for _, part := range strings.Split(s, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			if i, err := strconv.Atoi(trimmed); err == nil {
				result = append(result, i)
			}
		}
	}
	return result
}
