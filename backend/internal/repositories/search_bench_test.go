package repositories

import (
	"database/sql"
	"fmt"
	"math/rand"
	"path/filepath"
	"testing"

	"recipe-app/internal/database"
	"recipe-app/internal/models"
)

// Word pools used to generate recipes whose title/description/ingredient text
// looks enough like real data for LIKE '%term%' and (future) FTS5 matching to
// behave representatively — not just N copies of one string.
var (
	benchAdjectives  = []string{"Spicy", "Creamy", "Roasted", "Grilled", "Baked", "Fresh", "Rustic", "Smoky", "Zesty", "Slow-Cooked"}
	benchMains       = []string{"Chicken", "Beef", "Salmon", "Tofu", "Lentil", "Mushroom", "Shrimp", "Pork", "Chickpea", "Eggplant"}
	benchDishes      = []string{"Curry", "Stew", "Soup", "Salad", "Tacos", "Risotto", "Stir Fry", "Pasta", "Casserole", "Pie"}
	benchCuisines    = []string{"Italian", "Mexican", "Indian", "Thai", "Greek", "Japanese", "French", "Spanish", "Moroccan", "Korean"}
	benchCategories  = []string{"Main", "Starter", "Dessert", "Breakfast", "Side", "Snack"}
	benchDifficulty  = []string{"easy", "medium", "hard"}
	benchIngredients = []string{
		"onion", "garlic", "olive oil", "tomato", "basil", "cumin", "paprika", "coconut milk",
		"soy sauce", "ginger", "lime", "cilantro", "black pepper", "butter", "flour", "sugar",
		"chili flakes", "parmesan", "yogurt", "lemon",
	}
)

// seedBenchRecipes inserts n synthetic recipes (each with ~5 ingredients) via
// the real repository (not raw SQL) so the seeded data goes through the same
// code path production writes do. Runs before b.ResetTimer() in every caller.
func seedBenchRecipes(b *testing.B, repo *RecipeRepository, n int) {
	b.Helper()
	rng := rand.New(rand.NewSource(42))

	for i := 0; i < n; i++ {
		adj := benchAdjectives[rng.Intn(len(benchAdjectives))]
		main := benchMains[rng.Intn(len(benchMains))]
		dish := benchDishes[rng.Intn(len(benchDishes))]
		cuisine := benchCuisines[rng.Intn(len(benchCuisines))]

		ingredients := make([]models.Ingredient, 0, 5)
		for j := 0; j < 5; j++ {
			ingredients = append(ingredients, models.Ingredient{
				Name:   benchIngredients[rng.Intn(len(benchIngredients))],
				Amount: fmt.Sprintf("%d", rng.Intn(5)+1),
				Unit:   "unit",
			})
		}

		recipe := &models.Recipe{
			UserID:      "bench-user",
			Title:       fmt.Sprintf("%s %s %s #%d", adj, main, dish, i),
			Description: fmt.Sprintf("A %s %s dish made with %s, perfect as a %s course.", adj, cuisine, main, dish),
			PrepTime:    rng.Intn(30) + 5,
			CookTime:    rng.Intn(60) + 5,
			Servings:    rng.Intn(6) + 1,
			Difficulty:  benchDifficulty[rng.Intn(len(benchDifficulty))],
			Category:    benchCategories[rng.Intn(len(benchCategories))],
			Cuisine:     cuisine,
			Ingredients: ingredients,
		}
		if err := repo.CreateRecipe(recipe); err != nil {
			b.Fatalf("seed CreateRecipe #%d: %v", i, err)
		}
	}
}

func newBenchDB(b *testing.B) *sql.DB {
	b.Helper()
	dir := b.TempDir()
	db, err := sql.Open("sqlite", filepath.Join(dir, "bench.db"))
	if err != nil {
		b.Fatalf("open sqlite: %v", err)
	}
	b.Cleanup(func() { db.Close() })
	if err := database.ApplySchema(db); err != nil {
		b.Fatalf("apply schema: %v", err)
	}
	return db
}

// BenchmarkSearchRecipes_TextQuery measures the current LIKE '%term%' path
// (title/description/ingredient-name) at increasing table sizes — this is
// the baseline docs/search/BENCHMARKS.md compares an FTS5 implementation
// against.
func BenchmarkSearchRecipes_TextQuery(b *testing.B) {
	for _, n := range []int{100, 1000, 5000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			db := newBenchDB(b)
			repo := NewRecipeRepository(db)
			seedBenchRecipes(b, repo, n)

			filter := &models.RecipeFilter{Query: "Chicken"}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, _, err := repo.SearchRecipes(filter, 20, 0); err != nil {
					b.Fatalf("SearchRecipes: %v", err)
				}
			}
		})
	}
}

// BenchmarkSearchRecipes_FilterOnly measures the exact-match filter path
// (no text query, no indexes on difficulty/category/cuisine today) at the
// same table sizes, isolating the cost of the LIKE clause from the cost of
// the plain equality scans.
func BenchmarkSearchRecipes_FilterOnly(b *testing.B) {
	for _, n := range []int{100, 1000, 5000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			db := newBenchDB(b)
			repo := NewRecipeRepository(db)
			seedBenchRecipes(b, repo, n)

			filter := &models.RecipeFilter{Difficulty: "medium", Cuisine: "Italian"}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, _, err := repo.SearchRecipes(filter, 20, 0); err != nil {
					b.Fatalf("SearchRecipes: %v", err)
				}
			}
		})
	}
}

// BenchmarkSuggestTitles measures the prefix-LIKE autocomplete path, which
// FTS_DESIGN.md notes stays on LIKE (a plain index serves left-anchored LIKE
// in SQLite) rather than moving to FTS5.
func BenchmarkSuggestTitles(b *testing.B) {
	for _, n := range []int{100, 1000, 5000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			db := newBenchDB(b)
			repo := NewRecipeRepository(db)
			seedBenchRecipes(b, repo, n)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err := repo.SuggestTitles("Spicy", 10); err != nil {
					b.Fatalf("SuggestTitles: %v", err)
				}
			}
		})
	}
}
