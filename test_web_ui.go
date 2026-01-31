package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("Testing web UI...")
	testWebUI()
}

func testWebUI() {
	// Start server in background
	go func() {
		os.Chdir("backend")
		os.Setenv("JWT_SECRET", "test-secret")
		fmt.Println("Starting server for test...")
		if err := http.ListenAndServe(":8081", nil); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Test main page
	resp, err := http.Get("http://localhost:8081/")
	if err != nil {
		log.Fatalf("Failed to get homepage: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Expected status 200, got %d", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response: %v", err)
		return
	}

	// Check if HTML content is present
	html := string(body)
	if !contains(html, "<!DOCTYPE html>") {
		fmt.Println("❌ Response should contain HTML DOCTYPE")
		return
	}

	if !contains(html, "RecipeApp") {
		fmt.Println("❌ Response should contain RecipeApp branding")
		return
	}

	// Test recipes page
	resp2, err := http.Get("http://localhost:8081/recipes")
	if err != nil {
		fmt.Printf("Failed to get recipes page: %v", err)
		return
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != 200 {
		fmt.Printf("❌ Expected status 200 for recipes page, got %d", resp2.StatusCode)
		return
	}

	// Test API endpoint
	resp3, err := http.Get("http://localhost:8081/api/recipes")
	if err != nil {
		fmt.Printf("Failed to get API recipes: %v", err)
		return
	}
	defer resp3.Body.Close()

	if resp3.StatusCode != 200 {
		fmt.Printf("❌ Expected status 200 for API recipes, got %d", resp3.StatusCode)
		return
	}

	fmt.Println("✅ Web UI tests passed!")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
