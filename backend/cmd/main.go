package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	appmiddleware "recipe-app/internal/appmiddleware"
	"recipe-app/internal/handlers"
	"recipe-app/internal/logger"
)

func main() {
	log := logger.New()

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

	r.Get("/", handlers.NewWebHandler().HandleIndex)

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", handlers.NewAuthHandler(authService).HandleRegister)
			r.Post("/login", handlers.NewAuthHandler(authService).HandleLogin)
			r.Post("/refresh", handlers.NewAuthHandler(authService).HandleRefresh)
		})

		r.Route("/recipes", func(r chi.Router) {
			r.With(authService.OptionalAuthMiddleware).Get("/", handlers.NewAPIHandler().HandleRecipes)
			r.With(authService.AuthMiddleware).Post("/", handlers.NewAPIHandler().HandleCreateRecipe)
			r.Route("/{id}", func(r chi.Router) {
				r.With(authService.OptionalAuthMiddleware).Get("/", handlers.NewAPIHandler().HandleRecipe)
				r.With(authService.AuthMiddleware).Put("/", handlers.NewAPIHandler().HandleUpdateRecipe)
				r.With(authService.AuthMiddleware).Delete("/", handlers.NewAPIHandler().HandleDeleteRecipe)
			})
		})

		r.With(authService.AuthMiddleware).Get("/users/profile", handlers.NewUserHandler().HandleProfile)
		r.With(authService.AuthMiddleware).Put("/users/profile", handlers.NewUserHandler().HandleUpdateProfile)
	})

	r.Route("/recipes", func(r chi.Router) {
		r.Get("/", handlers.NewWebHandler().HandleRecipes)
		r.With(authService.AuthMiddleware).Get("/new", handlers.NewWebHandler().HandleNewRecipe)
	})

	log.Info("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Error("Server failed to start", "error", err)
	}
}
