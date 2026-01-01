package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"vexora-studio/internal/api"
	"vexora-studio/internal/dashboard"
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

	// Static Frontend (for testing)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	// API Endpoints
	mux.HandleFunc("POST /instagram", api.HandleCreateInstagramFeed)
	mux.HandleFunc("GET /instagram", api.HandleGetTodaysInstagramFeeds)
	mux.HandleFunc("GET /instagram/{identifier}", api.HandleGetInstagramFeeds)

	mux.HandleFunc("POST /twitter", api.HandleCreateTwitterFeed)
	mux.HandleFunc("GET /twitter", api.HandleGetTodaysTwitterFeeds)
	mux.HandleFunc("GET /twitter/{identifier}", api.HandleGetTwitterFeeds)

	mux.HandleFunc("POST /linkedin", api.HandleCreateLinkedinFeed)
	mux.HandleFunc("GET /linkedin", api.HandleGetTodaysLinkedinFeeds)
	mux.HandleFunc("GET /linkedin/{identifier}", api.HandleGetLinkedinFeeds)

	mux.HandleFunc("POST /newsletter", api.HandleCreateNewsletterFeed)
	mux.HandleFunc("GET /newsletter", api.HandleGetTodaysNewsletterFeeds)
	mux.HandleFunc("GET /newsletter/{identifier}", api.HandleGetNewsletterFeeds)

	// 4. Start Server
	port := ":8081"
	log.Printf("📸 Vexora Studio listening on %s", port)
	log.Println("🖥️  Dashboard available at http://localhost:8081/dashboard")
	dashboard.StartDashboard(":8082")

	server := &http.Server{
		Addr:    port,
		Handler: mux,
		// Generous timeouts for AI generation
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 300 * time.Second,
	}

	log.Println("Sending the mail")
	// smtp.Mail()

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
