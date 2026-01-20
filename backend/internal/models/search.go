package models

import (
	"database/sql"
	"fmt"
	"strings"
)

type SearchService struct {
	db *sql.DB
}

func NewSearchService(db *sql.DB) *SearchService {
	return &SearchService{db: db}
}

type SearchFilters struct {
	Query       string   `json:"query"`
	UserID      *int     `json:"user_id,omitempty"`
	Difficulty  string   `json:"difficulty,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Categories  []int    `json:"categories,omitempty"`
	MinPrepTime *int     `json:"min_prep_time,omitempty"`
	MaxPrepTime *int     `json:"max_prep_time,omitempty"`
	MinCookTime *int     `json:"min_cook_time,omitempty"`
	MaxCookTime *int     `json:"max_cook_time,omitempty"`
	MinServings *int     `json:"min_servings,omitempty"`
	MaxServings *int     `json:"max_servings,omitempty"`
	IsPublic    *bool    `json:"is_public,omitempty"`
	SortBy      string   `json:"sort_by,omitempty"`    // relevance, created_at, title, prep_time, cook_time
	SortOrder   string   `json:"sort_order,omitempty"` // ASC, DESC
	Limit       *int     `json:"limit,omitempty"`
	Offset      *int     `json:"offset,omitempty"`
}

type SearchResult struct {
	Recipes     []Recipe `json:"recipes"`
	TotalCount  int      `json:"total_count"`
	CurrentPage int      `json:"current_page"`
	PageSize    int      `json:"page_size"`
	HasNext     bool     `json:"has_next"`
	HasPrev     bool     `json:"has_prev"`
}

func (ss *SearchService) SearchRecipes(filters SearchFilters) (*SearchResult, error) {
	// Build WHERE clause
	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	// Text search using PostgreSQL full-text search
	if filters.Query != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("(to_tsvector('english', r.title || ' ' || COALESCE(r.description, '') || ' ' || COALESCE(r.instructions, '')) @@ plainto_tsquery($%d))", argIndex))
		args = append(args, filters.Query)
		argIndex++
	}

	// User filter
	if filters.UserID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("r.user_id = $%d", argIndex))
		args = append(args, *filters.UserID)
		argIndex++
	}

	// Difficulty filter
	if filters.Difficulty != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("r.difficulty = $%d", argIndex))
		args = append(args, filters.Difficulty)
		argIndex++
	}

	// Tags filter (ANY tag in the list)
	if len(filters.Tags) > 0 {
		tagPlaceholders := make([]string, len(filters.Tags))
		for i := range filters.Tags {
			tagPlaceholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, filters.Tags[i])
			argIndex++
		}
		whereConditions = append(whereConditions, fmt.Sprintf("EXISTS (SELECT 1 FROM recipe_tags rt JOIN tags t ON rt.tag_id = t.id WHERE rt.recipe_id = r.id AND t.name IN (%s))", strings.Join(tagPlaceholders, ", ")))
	}

	// Categories filter
	if len(filters.Categories) > 0 {
		catPlaceholders := make([]string, len(filters.Categories))
		for i := range filters.Categories {
			catPlaceholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, filters.Categories[i])
			argIndex++
		}
		whereConditions = append(whereConditions, fmt.Sprintf("EXISTS (SELECT 1 FROM recipe_categories rc WHERE rc.recipe_id = r.id AND rc.category_id IN (%s))", strings.Join(catPlaceholders, ", ")))
	}

	// Time filters
	if filters.MinPrepTime != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("(r.prep_time >= $%d OR r.prep_time IS NULL)", argIndex))
		args = append(args, *filters.MinPrepTime)
		argIndex++
	}

	if filters.MaxPrepTime != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("(r.prep_time <= $%d OR r.prep_time IS NULL)", argIndex))
		args = append(args, *filters.MaxPrepTime)
		argIndex++
	}

	if filters.MinCookTime != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("(r.cook_time >= $%d OR r.cook_time IS NULL)", argIndex))
		args = append(args, *filters.MinCookTime)
		argIndex++
	}

	if filters.MaxCookTime != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("(r.cook_time <= $%d OR r.cook_time IS NULL)", argIndex))
		args = append(args, *filters.MaxCookTime)
		argIndex++
	}

	// Servings filter
	if filters.MinServings != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("r.servings >= $%d", argIndex))
		args = append(args, *filters.MinServings)
		argIndex++
	}

	if filters.MaxServings != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("r.servings <= $%d", argIndex))
		args = append(args, *filters.MaxServings)
		argIndex++
	}

	// Public/private filter
	if filters.IsPublic != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("r.is_public = $%d", argIndex))
		args = append(args, *filters.IsPublic)
		argIndex++
	}

	// Combine WHERE conditions
	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// ORDER BY clause with relevance scoring
	orderBy := "ORDER BY r.created_at DESC"
	if filters.Query != "" {
		if filters.SortBy == "relevance" {
			orderBy = fmt.Sprintf("ORDER BY ts_rank(to_tsvector('english', r.title || ' ' || COALESCE(r.description, '') || ' ' || COALESCE(r.instructions, '')), plainto_tsquery($%d)) DESC, r.created_at DESC", argIndex)
		} else if filters.SortBy == "title" {
			orderBy = fmt.Sprintf("ORDER BY r.title %s, r.created_at DESC", filters.SortOrder)
		} else if filters.SortBy == "prep_time" {
			orderBy = fmt.Sprintf("ORDER BY COALESCE(r.prep_time, 0) %s, r.created_at DESC", filters.SortOrder)
		} else if filters.SortBy == "cook_time" {
			orderBy = fmt.Sprintf("ORDER BY COALESCE(r.cook_time, 0) %s, r.created_at DESC", filters.SortOrder)
		} else if filters.SortBy == "created_at" {
			orderBy = fmt.Sprintf("ORDER BY r.created_at %s", filters.SortOrder)
		}
	} else {
		// No query, just sort by requested field
		if filters.SortBy == "title" {
			orderBy = fmt.Sprintf("ORDER BY r.title %s", filters.SortOrder)
		} else if filters.SortBy == "prep_time" {
			orderBy = fmt.Sprintf("ORDER BY COALESCE(r.prep_time, 0) %s", filters.SortOrder)
		} else if filters.SortBy == "cook_time" {
			orderBy = fmt.Sprintf("ORDER BY COALESCE(r.cook_time, 0) %s", filters.SortOrder)
		}
	}

	// Build the main query
	query := fmt.Sprintf(`
		SELECT r.id, r.user_id, r.title, r.description, r.instructions, 
			   r.prep_time, r.cook_time, r.servings, r.difficulty, 
			   r.image_url, r.is_public, r.created_at, r.updated_at
		FROM recipes r
		%s
		%s
	`, whereClause, orderBy)

	// Add LIMIT and OFFSET
	limitClause := ""
	if filters.Limit != nil {
		limitClause = fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, *filters.Limit)
		argIndex++
	}

	offsetClause := ""
	if filters.Offset != nil {
		offsetClause = fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, *filters.Offset)
		argIndex++
	}

	// Execute main query
	fullQuery := query + limitClause + offsetClause
	rows, err := ss.db.Query(fullQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search query: %w", err)
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var recipe Recipe
		err := rows.Scan(
			&recipe.ID, &recipe.UserID, &recipe.Title, &recipe.Description,
			&recipe.Instructions, &recipe.PrepTime, &recipe.CookTime,
			&recipe.Servings, &recipe.Difficulty, &recipe.ImageURL,
			&recipe.IsPublic, &recipe.CreatedAt, &recipe.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %w", err)
		}
		recipes = append(recipes, recipe)
	}

	// Get total count for pagination
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM recipes r
		%s
	`, whereClause)

	var totalCount int
	err = ss.db.QueryRow(countQuery, args[:len(args)-len(limitClause)/len(" LIMIT ")-len(offsetClause)/len(" OFFSET ")]...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Calculate pagination info
	pageSize := 20
	if filters.Limit != nil {
		pageSize = *filters.Limit
	}

	currentPage := 1
	if filters.Offset != nil && pageSize > 0 {
		currentPage = (*filters.Offset / pageSize) + 1
	}

	hasNext := (currentPage * pageSize) < totalCount
	hasPrev := currentPage > 1

	return &SearchResult{
		Recipes:     recipes,
		TotalCount:  totalCount,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
	}, nil
}

func (ss *SearchService) GetSearchSuggestions(query string, limit int) ([]string, error) {
	if len(query) < 2 {
		return []string{}, nil
	}

	sql := `
		SELECT DISTINCT title 
		FROM recipes 
		WHERE title ILIKE $1 
		ORDER BY length(title), title
		LIMIT $2
	`

	rows, err := ss.db.Query(sql, query+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggestions: %w", err)
	}
	defer rows.Close()

	var suggestions []string
	for rows.Next() {
		var title string
		err := rows.Scan(&title)
		if err != nil {
			return nil, fmt.Errorf("failed to scan suggestion: %w", err)
		}
		suggestions = append(suggestions, title)
	}

	return suggestions, nil
}

func (ss *SearchService) GetPopularTags(limit int) ([]Tag, error) {
	sql := `
		SELECT t.id, t.name, t.color, t.created_at, COUNT(rt.recipe_id) as recipe_count
		FROM tags t
		JOIN recipe_tags rt ON t.id = rt.tag_id
		GROUP BY t.id, t.name, t.color, t.created_at
		ORDER BY recipe_count DESC, t.name
		LIMIT $1
	`

	rows, err := ss.db.Query(sql, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular tags: %w", err)
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Color, &tag.CreatedAt, new(int)) // recipe_count ignored for now
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
