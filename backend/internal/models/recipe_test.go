package models

import (
	"testing"
	"time"
)

func TestRecipe_Validation(t *testing.T) {
	tests := []struct {
		name        string
		recipe      Recipe
		expectError bool
	}{
		{
			name: "Valid recipe",
			recipe: Recipe{
				ID:           "1",
				Title:        "Test Recipe",
				Description:  "A test recipe",
				PrepTime:     30,
				CookTime:     45,
				Servings:     4,
				Difficulty:   "medium",
				Category:     "dinner",
				Cuisine:      "italian",
				Ingredients:  []Ingredient{{ID: "1", Name: "flour", Amount: "2", Unit: "cups"}},
				Instructions: []Instruction{{ID: "1", Text: "Mix ingredients", Position: 1}},
				Tags:         []string{"easy", "quick"},
				ImageURL:     "http://example.com/image.jpg",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			expectError: false,
		},
		{
			name: "Empty title",
			recipe: Recipe{
				ID:          "1",
				Title:       "",
				Description: "A test recipe",
				PrepTime:    30,
				CookTime:    45,
				Servings:    4,
				Difficulty:  "medium",
			},
			expectError: true,
		},
		{
			name: "Negative prep time",
			recipe: Recipe{
				ID:          "1",
				Title:       "Test Recipe",
				Description: "A test recipe",
				PrepTime:    -5,
				CookTime:    45,
				Servings:    4,
				Difficulty:  "medium",
			},
			expectError: true,
		},
		{
			name: "Invalid difficulty",
			recipe: Recipe{
				ID:          "1",
				Title:       "Test Recipe",
				Description: "A test recipe",
				PrepTime:    30,
				CookTime:    45,
				Servings:    4,
				Difficulty:  "invalid",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.recipe.Validate()
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestRecipeFilter_Validation(t *testing.T) {
	tests := []struct {
		name        string
		filter      RecipeFilter
		expectError bool
	}{
		{
			name: "Valid filter",
			filter: RecipeFilter{
				Category:    "dinner",
				Cuisine:     "italian",
				Difficulty:  "medium",
				Tags:        []string{"easy", "quick"},
				MaxPrepTime: 60,
				MaxCookTime: 120,
				MinServings: 2,
				MaxServings: 6,
			},
			expectError: false,
		},
		{
			name: "Invalid difficulty",
			filter: RecipeFilter{
				Difficulty: "invalid",
			},
			expectError: true,
		},
		{
			name: "Min servings greater than max servings",
			filter: RecipeFilter{
				MinServings: 8,
				MaxServings: 4,
			},
			expectError: true,
		},
		{
			name: "Negative max prep time",
			filter: RecipeFilter{
				MaxPrepTime: -10,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.filter.Validate()
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestIngredient_Validation(t *testing.T) {
	tests := []struct {
		name        string
		ingredient  Ingredient
		expectError bool
	}{
		{
			name: "Valid ingredient",
			ingredient: Ingredient{
				ID:       "1",
				RecipeID: "1",
				Name:     "flour",
				Amount:   "2",
				Unit:     "cups",
				Position: 1,
			},
			expectError: false,
		},
		{
			name: "Empty name",
			ingredient: Ingredient{
				ID:       "1",
				RecipeID: "1",
				Name:     "",
				Amount:   "2",
				Unit:     "cups",
				Position: 1,
			},
			expectError: true,
		},
		{
			name: "Negative position",
			ingredient: Ingredient{
				ID:       "1",
				RecipeID: "1",
				Name:     "flour",
				Amount:   "2",
				Unit:     "cups",
				Position: -1,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ingredient.Validate()
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestInstruction_Validation(t *testing.T) {
	tests := []struct {
		name        string
		instruction Instruction
		expectError bool
	}{
		{
			name: "Valid instruction",
			instruction: Instruction{
				ID:          "1",
				RecipeID:    "1",
				Text:        "Mix ingredients",
				Position:    1,
				Duration:    5,
				Temperature: 350,
			},
			expectError: false,
		},
		{
			name: "Empty text",
			instruction: Instruction{
				ID:       "1",
				RecipeID: "1",
				Text:     "",
				Position: 1,
			},
			expectError: true,
		},
		{
			name: "Negative position",
			instruction: Instruction{
				ID:       "1",
				RecipeID: "1",
				Text:     "Mix ingredients",
				Position: -1,
			},
			expectError: true,
		},
		{
			name: "Negative duration",
			instruction: Instruction{
				ID:       "1",
				RecipeID: "1",
				Text:     "Mix ingredients",
				Position: 1,
				Duration: -5,
			},
			expectError: true,
		},
		{
			name: "Negative temperature",
			instruction: Instruction{
				ID:          "1",
				RecipeID:    "1",
				Text:        "Mix ingredients",
				Position:    1,
				Temperature: -10,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.instruction.Validate()
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestSearchResult_Pagination(t *testing.T) {
	tests := []struct {
		name       string
		result     SearchResult
		totalPages int
	}{
		{
			name: "First page",
			result: SearchResult{
				Recipes: []Recipe{{}, {}, {}},
				Total:   10,
				Page:    1,
				PerPage: 3,
			},
			totalPages: 4,
		},
		{
			name: "Last page with partial results",
			result: SearchResult{
				Recipes: []Recipe{{}},
				Total:   10,
				Page:    4,
				PerPage: 3,
			},
			totalPages: 4,
		},
		{
			name: "Single page",
			result: SearchResult{
				Recipes: []Recipe{{}, {}},
				Total:   2,
				Page:    1,
				PerPage: 10,
			},
			totalPages: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totalPages := tt.result.TotalPages()
			if totalPages != tt.totalPages {
				t.Errorf("Expected total pages %d, got %d", tt.totalPages, totalPages)
			}
		})
	}
}
