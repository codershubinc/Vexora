package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"vexora-studio/internal/api"
	"vexora-studio/internal/database"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 1. Setup Data Directory
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("❌ Failed to create data directory: %v", err)
	}

	// 2. Initialize Database
	if err := database.Init("./data/vexora.db"); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}

	// 3. Setup Router
	mux := http.NewServeMux()

	// Web Interface
	mux.HandleFunc("GET /", api.HandleHome)

	// API Endpoints
	mux.HandleFunc("POST /instagram", api.HandleCreateInstagramFeed)
	mux.HandleFunc("POST /twitter", api.HandleCreateTwitterFeed)
	mux.HandleFunc("POST /linkedin", api.HandleCreateLinkedinFeed)
	mux.HandleFunc("POST /newsletter", api.HandleCreateNewsletterFeed)

	// 4. Start Server
	port := ":8081"
	log.Printf("📸 Vexora Studio listening on %s", port)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
		// Generous timeouts for AI generation
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 300 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
