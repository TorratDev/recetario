package models

import (
	"database/sql"
	"fmt"
)

type IngredientService struct {
	db *sql.DB
}

func NewIngredientService(db *sql.DB) *IngredientService {
	return &IngredientService{db: db}
}

func (is *IngredientService) Create(name, category string) (*Ingredient, error) {
	query := `
		INSERT INTO ingredients (name, category)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	ingredient := &Ingredient{
		Name:     name,
		Category: &category,
	}

	err := is.db.QueryRow(query, name, category).Scan(&ingredient.ID, &ingredient.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create ingredient: %w", err)
	}

	return ingredient, nil
}

func (is *IngredientService) GetByID(id int) (*Ingredient, error) {
	query := `
		SELECT id, name, category, created_at
		FROM ingredients
		WHERE id = $1
	`

	ingredient := &Ingredient{}
	err := is.db.QueryRow(query, id).Scan(&ingredient.ID, &ingredient.Name, &ingredient.Category, &ingredient.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get ingredient: %w", err)
	}

	return ingredient, nil
}

func (is *IngredientService) GetAll() ([]Ingredient, error) {
	query := `
		SELECT id, name, category, created_at
		FROM ingredients
		ORDER BY name
	`

	rows, err := is.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query ingredients: %w", err)
	}
	defer rows.Close()

	var ingredients []Ingredient
	for rows.Next() {
		var ingredient Ingredient
		err := rows.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Category, &ingredient.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ingredient: %w", err)
		}
		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}

func (is *IngredientService) Search(query string) ([]Ingredient, error) {
	sql := `
		SELECT id, name, category, created_at
		FROM ingredients
		WHERE name ILIKE $1
		ORDER BY name
		LIMIT 10
	`

	rows, err := is.db.Query(sql, "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search ingredients: %w", err)
	}
	defer rows.Close()

	var ingredients []Ingredient
	for rows.Next() {
		var ingredient Ingredient
		err := rows.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Category, &ingredient.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ingredient: %w", err)
		}
		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}

func (is *IngredientService) Update(id int, name, category string) (*Ingredient, error) {
	query := `
		UPDATE ingredients 
		SET name = $1, category = $2
		WHERE id = $3
	`

	_, err := is.db.Exec(query, name, category, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update ingredient: %w", err)
	}

	return is.GetByID(id)
}

func (is *IngredientService) Delete(id int) error {
	query := "DELETE FROM ingredients WHERE id = $1"

	result, err := is.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete ingredient: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ingredient not found")
	}

	return nil
}

type TagService struct {
	db *sql.DB
}

func NewTagService(db *sql.DB) *TagService {
	return &TagService{db: db}
}

func (ts *TagService) Create(name, color string) (*Tag, error) {
	query := `
		INSERT INTO tags (name, color)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	tag := &Tag{
		Name:  name,
		Color: color,
	}

	err := ts.db.QueryRow(query, name, color).Scan(&tag.ID, &tag.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return tag, nil
}

func (ts *TagService) GetByID(id int) (*Tag, error) {
	query := `
		SELECT id, name, color, created_at
		FROM tags
		WHERE id = $1
	`

	tag := &Tag{}
	err := ts.db.QueryRow(query, id).Scan(&tag.ID, &tag.Name, &tag.Color, &tag.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	return tag, nil
}

func (ts *TagService) GetAll() ([]Tag, error) {
	query := `
		SELECT id, name, color, created_at
		FROM tags
		ORDER BY name
	`

	rows, err := ts.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Color, &tag.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (ts *TagService) Update(id int, name, color string) (*Tag, error) {
	query := `
		UPDATE tags 
		SET name = $1, color = $2
		WHERE id = $3
	`

	_, err := ts.db.Exec(query, name, color, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	return ts.GetByID(id)
}

func (ts *TagService) Delete(id int) error {
	query := "DELETE FROM tags WHERE id = $1"

	result, err := ts.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tag not found")
	}

	return nil
}
