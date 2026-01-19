package main

import (
	"log"
	"net/http"
	"recipe-app/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	
	// API routes
	apiHandler := handlers.NewAPIHandler()
	mux.HandleFunc("/api/recipes", apiHandler.HandleRecipes)
	mux.HandleFunc("/api/recipes/", apiHandler.HandleRecipe)
	
	// Web routes (HTMX)
	webHandler := handlers.NewWebHandler()
	mux.HandleFunc("/", webHandler.HandleIndex)
	mux.HandleFunc("/recipes", webHandler.HandleRecipes)
	mux.HandleFunc("/recipes/new", webHandler.HandleNewRecipe)
	
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}