// Command seed populates a SQLite database with a large number of synthetic
// recipes, for manually exercising the app (or the running server) at a
// scale beyond the handful of rows main.go seeds on every startup. Unlike
// the seeding built into internal/repositories/search_bench_test.go, this
// writes a real file you can point RecipeApp's server at and keep around,
// rather than a throwaway benchmark temp DB.
//
// Usage:
//
//	go run ./cmd/seed -db recipeapp-bench.db -count 5000
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"

	_ "modernc.org/sqlite"

	"recipe-app/internal/database"
	"recipe-app/internal/models"
	"recipe-app/internal/repositories"
)

var (
	adjectives     = []string{"Spicy", "Creamy", "Roasted", "Grilled", "Baked", "Fresh", "Rustic", "Smoky", "Zesty", "Slow-Cooked"}
	mains          = []string{"Chicken", "Beef", "Salmon", "Tofu", "Lentil", "Mushroom", "Shrimp", "Pork", "Chickpea", "Eggplant"}
	dishes         = []string{"Curry", "Stew", "Soup", "Salad", "Tacos", "Risotto", "Stir Fry", "Pasta", "Casserole", "Pie"}
	cuisines       = []string{"Italian", "Mexican", "Indian", "Thai", "Greek", "Japanese", "French", "Spanish", "Moroccan", "Korean"}
	categories     = []string{"Main", "Starter", "Dessert", "Breakfast", "Side", "Snack"}
	difficulty     = []string{"easy", "medium", "hard"}
	ingredientPool = []string{
		"onion", "garlic", "olive oil", "tomato", "basil", "cumin", "paprika", "coconut milk",
		"soy sauce", "ginger", "lime", "cilantro", "black pepper", "butter", "flour", "sugar",
		"chili flakes", "parmesan", "yogurt", "lemon",
	}
)

func main() {
	dbPath := flag.String("db", "recipeapp-bench.db", "path to the SQLite database file to create/populate")
	count := flag.Int("count", 5000, "number of synthetic recipes to insert")
	seed := flag.Int64("seed", 42, "random seed, for reproducible datasets")
	flag.Parse()

	db, err := sql.Open("sqlite", *dbPath)
	if err != nil {
		log.Fatalf("open %s: %v", *dbPath, err)
	}
	defer db.Close()

	if err := database.ApplySchema(db); err != nil {
		log.Fatalf("apply schema: %v", err)
	}

	if _, err := db.Exec(
		`INSERT OR IGNORE INTO users (id, email, username, first_name, last_name, password_hash, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		"seed-user", "seed@recipeapp.local", "seed", "Seed", "User",
		"$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // bcrypt hash of "password"
	); err != nil {
		log.Fatalf("create seed user: %v", err)
	}

	repo := repositories.NewRecipeRepository(db)
	rng := rand.New(rand.NewSource(*seed))

	for i := 0; i < *count; i++ {
		recipe := randomRecipe(rng, i)
		if err := repo.CreateRecipe(recipe); err != nil {
			log.Fatalf("create recipe #%d: %v", i, err)
		}
		if (i+1)%500 == 0 {
			fmt.Printf("seeded %d/%d recipes\n", i+1, *count)
		}
	}

	fmt.Printf("done: %d recipes written to %s\n", *count, *dbPath)
}

func randomRecipe(rng *rand.Rand, i int) *models.Recipe {
	adj := adjectives[rng.Intn(len(adjectives))]
	main := mains[rng.Intn(len(mains))]
	dish := dishes[rng.Intn(len(dishes))]
	cuisine := cuisines[rng.Intn(len(cuisines))]

	ingredients := make([]models.Ingredient, 0, 5)
	for j := 0; j < 5; j++ {
		ingredients = append(ingredients, models.Ingredient{
			Name:   ingredientPool[rng.Intn(len(ingredientPool))],
			Amount: fmt.Sprintf("%d", rng.Intn(5)+1),
			Unit:   "unit",
		})
	}

	return &models.Recipe{
		UserID:      "seed-user",
		Title:       fmt.Sprintf("%s %s %s #%d", adj, main, dish, i),
		Description: fmt.Sprintf("A %s %s dish made with %s, perfect as a %s course.", adj, cuisine, main, dish),
		PrepTime:    rng.Intn(30) + 5,
		CookTime:    rng.Intn(60) + 5,
		Servings:    rng.Intn(6) + 1,
		Difficulty:  difficulty[rng.Intn(len(difficulty))],
		Category:    categories[rng.Intn(len(categories))],
		Cuisine:     cuisine,
		Ingredients: ingredients,
	}
}
