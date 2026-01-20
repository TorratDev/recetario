package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"recipe-app/internal/migration"
)

func main() {
	// Database connection string - in production, this should come from environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://recipe_user:recipe_password@localhost:5432/recipe_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	migrator := migration.NewMigrator(db)

	// Load migrations from files
	if err := migrator.LoadMigrations("migrations"); err != nil {
		log.Fatalf("Failed to load migrations: %v", err)
	}

	// Run migrations
	if err := migrator.Up(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("Migrations completed successfully!")
}
