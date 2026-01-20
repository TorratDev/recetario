package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"

	appmiddleware "recipe-app/internal/appmiddleware"
	"recipe-app/internal/handlers"
	"recipe-app/internal/logger"
)

func main() {
	log := logger.New()

	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://recipe_user:recipe_password@localhost:5432/recipe_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Error("Failed to ping database", "error", err)
		return
	}

	r := chi.NewRouter()

	rateLimiter := appmiddleware.NewRateLimiter(100, time.Minute)
	authService := appmiddleware.NewAuthService(os.Getenv("JWT_SECRET"))

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.Recoverer)
	r.Use(appmiddleware.RequestLogger)
	r.Use(appmiddleware.ErrorHandler)
	r.Use(appmiddleware.CORS(appmiddleware.DefaultCORSConfig()))
	r.Use(appmiddleware.RateLimit(rateLimiter))
	r.Use(appmiddleware.SecurityHeaders)

	r.Get("/", handlers.NewWebHandler(db).HandleIndex)

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", handlers.NewAuthHandler(authService).HandleRegister)
			r.Post("/login", handlers.NewAuthHandler(authService).HandleLogin)
			r.Post("/refresh", handlers.NewAuthHandler(authService).HandleRefresh)
		})

		r.Route("/recipes", func(r chi.Router) {
			r.With(authService.OptionalAuthMiddleware).Get("/", handlers.NewAPIHandler(db).HandleRecipes)
			r.With(authService.AuthMiddleware).Post("/", handlers.NewAPIHandler(db).HandleCreateRecipe)
			r.Route("/{id}", func(r chi.Router) {
				r.With(authService.OptionalAuthMiddleware).Get("/", handlers.NewAPIHandler(db).HandleRecipe)
				r.With(authService.AuthMiddleware).Put("/", handlers.NewAPIHandler(db).HandleUpdateRecipe)
				r.With(authService.AuthMiddleware).Delete("/", handlers.NewAPIHandler(db).HandleDeleteRecipe)
			})
		})

		r.With(authService.AuthMiddleware).Get("/users/profile", handlers.NewUserHandler().HandleProfile)
		r.With(authService.AuthMiddleware).Put("/users/profile", handlers.NewUserHandler().HandleUpdateProfile)
	})

	// Web routes (HTMX)
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.NewWebHandler(db).HandleIndex)
		r.Get("/recipe-list", handlers.NewWebHandler(db).HandleRecipeList) // HTMX endpoint
	})

	r.Route("/recipes", func(r chi.Router) {
		r.Get("/", handlers.NewWebHandler(db).HandleRecipes)
		r.With(authService.AuthMiddleware).Get("/new", handlers.NewWebHandler(db).HandleNewRecipe)
	})

	// Ingredient routes
	r.Route("/api/ingredients", func(r chi.Router) {
		r.With(authService.AuthMiddleware).Get("/", handlers.NewIngredientHandler(db).HandleIngredients)
		r.With(authService.AuthMiddleware).Post("/", handlers.NewIngredientHandler(db).HandleIngredients)
		r.Route("/{id}", func(r chi.Router) {
			r.With(authService.AuthMiddleware).Get("/", handlers.NewIngredientHandler(db).HandleIngredient)
			r.With(authService.AuthMiddleware).Put("/", handlers.NewIngredientHandler(db).HandleIngredient)
			r.With(authService.AuthMiddleware).Delete("/", handlers.NewIngredientHandler(db).HandleIngredient)
		})
	})

	// Tag routes
	r.Route("/api/tags", func(r chi.Router) {
		r.With(authService.AuthMiddleware).Get("/", handlers.NewTagHandler(db).HandleTags)
		r.With(authService.AuthMiddleware).Post("/", handlers.NewTagHandler(db).HandleTags)
		r.Route("/{id}", func(r chi.Router) {
			r.With(authService.AuthMiddleware).Get("/", handlers.NewTagHandler(db).HandleTag)
			r.With(authService.AuthMiddleware).Put("/", handlers.NewTagHandler(db).HandleTag)
			r.With(authService.AuthMiddleware).Delete("/", handlers.NewTagHandler(db).HandleTag)
		})
	})

	// File upload routes
	r.With(authService.AuthMiddleware).Post("/api/upload", handlers.NewFileHandler().HandleUpload)
	r.With(authService.AuthMiddleware).Post("/api/upload/multiple", handlers.NewFileHandler().HandleMultiUpload)
	// Search routes
	r.With(authService.OptionalAuthMiddleware).Get("/api/search", handlers.NewSearchHandler(db).HandleSearch)
	r.With(authService.OptionalAuthMiddleware).Get("/api/search/suggestions", handlers.NewSearchHandler(db).HandleSuggestions)
	r.With(authService.OptionalAuthMiddleware).Get("/api/search/tags/popular", handlers.NewSearchHandler(db).HandlePopularTags)

	r.With(authService.AuthMiddleware).Delete("/api/upload/{filename}", handlers.NewFileHandler().HandleDelete)
	r.Get("/uploads/{filename}", handlers.NewFileHandler().ServeFile)

	log.Info("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Error("Server failed to start", "error", err)
	}
}
