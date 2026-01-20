package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Recipe struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
	Instructions string    `json:"instructions" db:"instructions"`
	PrepTime     *int      `json:"prep_time" db:"prep_time"`
	CookTime     *int      `json:"cook_time" db:"cook_time"`
	Servings     int       `json:"servings" db:"servings"`
	Difficulty   string    `json:"difficulty" db:"difficulty"`
	ImageURL     *string   `json:"image_url" db:"image_url"`
	IsPublic     bool      `json:"is_public" db:"is_public"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`

	// Joined fields
	Ingredients []RecipeIngredient `json:"ingredients,omitempty"`
	Tags        []Tag              `json:"tags,omitempty"`
	Categories  []Category         `json:"categories,omitempty"`
	User        *User              `json:"user,omitempty"`
}

type RecipeIngredient struct {
	ID           int         `json:"id" db:"id"`
	RecipeID     int         `json:"recipe_id" db:"recipe_id"`
	IngredientID int         `json:"ingredient_id" db:"ingredient_id"`
	Quantity     *float64    `json:"quantity" db:"quantity"`
	Unit         *string     `json:"unit" db:"unit"`
	Notes        *string     `json:"notes" db:"notes"`
	Ingredient   *Ingredient `json:"ingredient,omitempty"`
}

type Ingredient struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Category  *string   `json:"category" db:"category"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Tag struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Color     string    `json:"color" db:"color"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Category struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	ParentID  *int      `json:"parent_id" db:"parent_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type RecipeService struct {
	db *sql.DB
}

func NewRecipeService(db *sql.DB) *RecipeService {
	return &RecipeService{db: db}
}

type RecipeFilter struct {
	UserID     *int     `json:"user_id,omitempty"`
	Title      *string  `json:"title,omitempty"`
	Difficulty *string  `json:"difficulty,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	IsPublic   *bool    `json:"is_public,omitempty"`
	Limit      *int     `json:"limit,omitempty"`
	Offset     *int     `json:"offset,omitempty"`
	SortBy     string   `json:"sort_by,omitempty"`    // created_at, title, prep_time
	SortOrder  string   `json:"sort_order,omitempty"` // ASC, DESC
}

func (rs *RecipeService) Create(recipe *Recipe) (*Recipe, error) {
	query := `
		INSERT INTO recipes (user_id, title, description, instructions, prep_time, cook_time, 
			servings, difficulty, image_url, is_public)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	err := rs.db.QueryRow(query,
		recipe.UserID, recipe.Title, recipe.Description, recipe.Instructions,
		recipe.PrepTime, recipe.CookTime, recipe.Servings, recipe.Difficulty,
		recipe.ImageURL, recipe.IsPublic).Scan(&recipe.ID, &recipe.CreatedAt, &recipe.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create recipe: %w", err)
	}

	return recipe, nil
}

func (rs *RecipeService) GetByID(id int, includeRelations bool) (*Recipe, error) {
	query := `
		SELECT id, user_id, title, description, instructions, prep_time, cook_time, 
			servings, difficulty, image_url, is_public, created_at, updated_at
		FROM recipes 
		WHERE id = $1
	`

	recipe := &Recipe{}
	err := rs.db.QueryRow(query, id).Scan(
		&recipe.ID, &recipe.UserID, &recipe.Title, &recipe.Description,
		&recipe.Instructions, &recipe.PrepTime, &recipe.CookTime,
		&recipe.Servings, &recipe.Difficulty, &recipe.ImageURL,
		&recipe.IsPublic, &recipe.CreatedAt, &recipe.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}

	if includeRelations {
		recipe.Ingredients, _ = rs.getIngredients(recipe.ID)
		recipe.Tags, _ = rs.getTags(recipe.ID)
		recipe.Categories, _ = rs.getCategories(recipe.ID)
	}

	return recipe, nil
}

func (rs *RecipeService) GetAll(filter RecipeFilter) ([]Recipe, error) {
	query := `
		SELECT id, user_id, title, description, instructions, prep_time, cook_time, 
			servings, difficulty, image_url, is_public, created_at, updated_at
		FROM recipes
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	// Build WHERE clauses
	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *filter.UserID)
		argIndex++
	}

	if filter.Title != nil {
		query += fmt.Sprintf(" AND title ILIKE $%d", argIndex)
		args = append(args, "%"+*filter.Title+"%")
		argIndex++
	}

	if filter.Difficulty != nil {
		query += fmt.Sprintf(" AND difficulty = $%d", argIndex)
		args = append(args, *filter.Difficulty)
		argIndex++
	}

	if filter.IsPublic != nil {
		query += fmt.Sprintf(" AND is_public = $%d", argIndex)
		args = append(args, *filter.IsPublic)
		argIndex++
	}

	// Add ORDER BY
	sortBy := "created_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}

	sortOrder := "DESC"
	if filter.SortOrder == "ASC" {
		sortOrder = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Add LIMIT and OFFSET
	if filter.Limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, *filter.Limit)
		argIndex++
	}

	if filter.Offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, *filter.Offset)
	}

	rows, err := rs.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query recipes: %w", err)
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var recipe Recipe
		err := rows.Scan(
			&recipe.ID, &recipe.UserID, &recipe.Title, &recipe.Description,
			&recipe.Instructions, &recipe.PrepTime, &recipe.CookTime,
			&recipe.Servings, &recipe.Difficulty, &recipe.ImageURL,
			&recipe.IsPublic, &recipe.CreatedAt, &recipe.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %w", err)
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (rs *RecipeService) Update(id int, recipe *Recipe) (*Recipe, error) {
	query := `
		UPDATE recipes 
		SET title = $1, description = $2, instructions = $3, prep_time = $4, 
			cook_time = $5, servings = $6, difficulty = $7, image_url = $8, 
			is_public = $9, updated_at = NOW()
		WHERE id = $10
		RETURNING updated_at
	`

	err := rs.db.QueryRow(query,
		recipe.Title, recipe.Description, recipe.Instructions, recipe.PrepTime,
		recipe.CookTime, recipe.Servings, recipe.Difficulty, recipe.ImageURL,
		recipe.IsPublic, id).Scan(&recipe.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update recipe: %w", err)
	}

	recipe.ID = id
	return recipe, nil
}

func (rs *RecipeService) Delete(id int) error {
	query := "DELETE FROM recipes WHERE id = $1"

	result, err := rs.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete recipe: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("recipe not found")
	}

	return nil
}

func (rs *RecipeService) getIngredients(recipeID int) ([]RecipeIngredient, error) {
	query := `
		SELECT ri.id, ri.recipe_id, ri.ingredient_id, ri.quantity, ri.unit, ri.notes,
			   i.id, i.name, i.category
		FROM recipe_ingredients ri
		LEFT JOIN ingredients i ON ri.ingredient_id = i.id
		WHERE ri.recipe_id = $1
	`

	rows, err := rs.db.Query(query, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []RecipeIngredient
	for rows.Next() {
		var ri RecipeIngredient
		var ing Ingredient

		err := rows.Scan(
			&ri.ID, &ri.RecipeID, &ri.IngredientID, &ri.Quantity, &ri.Unit, &ri.Notes,
			&ing.ID, &ing.Name, &ing.Category)

		if err != nil {
			return nil, err
		}

		ri.Ingredient = &ing
		ingredients = append(ingredients, ri)
	}

	return ingredients, nil
}

func (rs *RecipeService) getTags(recipeID int) ([]Tag, error) {
	query := `
		SELECT t.id, t.name, t.color, t.created_at
		FROM tags t
		JOIN recipe_tags rt ON t.id = rt.tag_id
		WHERE rt.recipe_id = $1
	`

	rows, err := rs.db.Query(query, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Color, &tag.CreatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (rs *RecipeService) getCategories(recipeID int) ([]Category, error) {
	query := `
		SELECT c.id, c.user_id, c.name, c.parent_id, c.created_at
		FROM categories c
		JOIN recipe_categories rc ON c.id = rc.category_id
		WHERE rc.recipe_id = $1
	`

	rows, err := rs.db.Query(query, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.UserID, &category.Name,
			&category.ParentID, &category.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
