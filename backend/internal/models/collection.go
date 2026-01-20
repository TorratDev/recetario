package models

import (
	"database/sql"
	"fmt"
)

type Collection struct {
	ID          int     `json:"id" db:"id"`
	UserID      int     `json:"user_id" db:"user_id"`
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
	CreatedAt   string  `json:"created_at" db:"created_at"`
	UpdatedAt   string  `json:"updated_at" db:"updated_at"`
}

type CollectionService struct {
	db *sql.DB
}

func NewCollectionService(db *sql.DB) *CollectionService {
	return &CollectionService{db: db}
}

type CreateCollectionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateCollectionRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type CollectionWithRecipes struct {
	ID          int      `json:"id"`
	UserID      int      `json:"user_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CreatedAt   string   `json:"created_at"`
	Recipes     []Recipe `json:"recipes"`
	RecipeCount int      `json:"recipe_count"`
}

func (cs *CollectionService) Create(userID int, req CreateCollectionRequest) (*Collection, error) {
	query := `
		INSERT INTO collections (user_id, name, description)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, created_at
	`

	collection := &Collection{
		UserID:      userID,
		Name:        req.Name,
		Description: &req.Description,
	}

	err := cs.db.QueryRow(query, userID, req.Name, req.Description).Scan(
		&collection.ID, &collection.Name, &collection.Description, &collection.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	return collection, nil
}

func (cs *CollectionService) GetByID(id int, userID int, includeRecipes bool) (*CollectionWithRecipes, error) {
	// Get collection basic info
	query := `
		SELECT id, user_id, name, description, created_at
		FROM collections
		WHERE id = $1 AND user_id = $2
	`

	collection := &CollectionWithRecipes{}
	err := cs.db.QueryRow(query, id, userID).Scan(
		&collection.ID, &collection.UserID, &collection.Name, &collection.Description, &collection.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	if includeRecipes {
		// Get recipes in this collection
		recipeQuery := `
			SELECT r.id, r.user_id, r.title, r.description, r.instructions, 
				   r.prep_time, r.cook_time, r.servings, r.difficulty, 
				   r.image_url, r.is_public, r.created_at, r.updated_at
			FROM recipes r
			JOIN recipe_collections rc ON r.id = rc.recipe_id
			WHERE rc.collection_id = $1
			ORDER BY r.created_at DESC
		`

		rows, err := cs.db.Query(recipeQuery, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get collection recipes: %w", err)
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

		collection.Recipes = recipes
		collection.RecipeCount = len(recipes)
	}

	return collection, nil
}

func (cs *CollectionService) GetAll(userID int) ([]Collection, error) {
	query := `
		SELECT id, user_id, name, description, created_at
		FROM collections
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := cs.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}
	defer rows.Close()

	var collections []Collection
	for rows.Next() {
		var collection Collection
		err := rows.Scan(
			&collection.ID, &collection.UserID, &collection.Name,
			&collection.Description, &collection.CreatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan collection: %w", err)
		}
		collections = append(collections, collection)
	}

	return collections, nil
}

func (cs *CollectionService) Update(id int, userID int, req UpdateCollectionRequest) (*Collection, error) {
	query := `
		UPDATE collections 
		SET name = COALESCE($1, name),
			description = COALESCE($2, description),
			updated_at = NOW()
		WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, name, description, created_at, updated_at
	`

	collection := &Collection{
		ID:     id,
		UserID: userID,
	}

	err := cs.db.QueryRow(query, req.Name, req.Description, id, userID).Scan(
		&collection.ID, &collection.UserID, &collection.Name, &collection.Description,
		&collection.CreatedAt, &collection.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update collection: %w", err)
	}

	return collection, nil
}

func (cs *CollectionService) Delete(id int, userID int) error {
	query := "DELETE FROM collections WHERE id = $1 AND user_id = $2"

	result, err := cs.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("collection not found")
	}

	return nil
}

func (cs *CollectionService) AddRecipe(collectionID int, recipeID int) error {
	query := `
		INSERT INTO recipe_collections (recipe_id, collection_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`

	_, err := cs.db.Exec(query, recipeID, collectionID)
	if err != nil {
		return fmt.Errorf("failed to add recipe to collection: %w", err)
	}

	return nil
}

func (cs *CollectionService) RemoveRecipe(collectionID int, recipeID int) error {
	query := `
		DELETE FROM recipe_collections 
		WHERE recipe_id = $1 AND collection_id = $2
	`

	result, err := cs.db.Exec(query, recipeID, collectionID)
	if err != nil {
		return fmt.Errorf("failed to remove recipe from collection: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	// Don't return error if recipe wasn't in collection
	_ = rowsAffected

	return nil
}

func (cs *CollectionService) GetRecipeCollections(recipeID int) ([]Collection, error) {
	query := `
		SELECT c.id, c.user_id, c.name, c.description, c.created_at
		FROM collections c
		JOIN recipe_collections rc ON c.id = rc.collection_id
		WHERE rc.recipe_id = $1
		ORDER BY c.name
	`

	rows, err := cs.db.Query(query, recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe collections: %w", err)
	}
	defer rows.Close()

	var collections []Collection
	for rows.Next() {
		var collection Collection
		err := rows.Scan(
			&collection.ID, &collection.UserID, &collection.Name,
			&collection.Description, &collection.CreatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan collection: %w", err)
		}
		collections = append(collections, collection)
	}

	return collections, nil
}
