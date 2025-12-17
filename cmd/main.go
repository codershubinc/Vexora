package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"vexora-studio/internal/api"
	"vexora-studio/internal/database"
	"vexora-studio/internal/middleware"
)

func main() {
	// 1. Setup Data Directory
	// Ensure the 'data' folder exists so SQLite doesn't crash
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("❌ Failed to create data directory: %v", err)
	}

	// 2. Initialize Database
	// This will create the 'journal_entries' table if it doesn't exist
	database.Init("./data/vexora.db")

	// 3. Configure Authentication
	// CRITICAL: This string must match the 'Secret' in your Sentinel config
	secretKey := "my-test-secret"
	authMiddleware := middleware.ValidateSignature(secretKey)

	// 4. Setup Router
	mux := http.NewServeMux()

	// Wrap the ingest handler with the auth middleware
	// This ensures no one can post fake journals without the key
	mux.Handle("POST /ingest", authMiddleware(http.HandlerFunc(api.HandleIngest)))

	// Web Interface (Public)
	mux.HandleFunc("GET /", api.HandleWebIndex)
	mux.HandleFunc("GET /project/{project}", api.HandleWebProject)
	mux.HandleFunc("POST /recurate/{id}", api.HandleRecurate)
	mux.HandleFunc("POST /fork/{id}", api.HandleFork)

	// 5. Start Server
	port := ":8080"
	log.Printf("📚 Vexora Studio (Journal Edition) listening on %s", port)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
		// Generous timeouts because local AI (Llama) can be slow
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 300 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
