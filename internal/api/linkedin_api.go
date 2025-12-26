package api

import (
	"encoding/json"
	"log"
	"net/http"
	"vexora-studio/internal/database"
	"vexora-studio/internal/llm"
)

func HandleCreateLinkedinFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	rawContent := r.FormValue("raw_content")
	projectName := r.FormValue("project_name")

	if rawContent == "" {
		http.Error(w, "Raw content is required", 400)
		return
	}

	data, err := llm.GenerateContent(llm.TypeLinkedIn, rawContent)
	if err != nil {
		log.Printf("❌ LinkedIn Generation Failed: %v", err)
		http.Error(w, "Content Generation Failed", 500)
		return
	}

	err = database.InsertLinkedinFeed(data, projectName)
	if err != nil {
		log.Printf("❌ LinkedIn DB Insert Failed: %v", err)
		http.Error(w, "Database Insertion Failed", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}

func HandleGetTodaysLinkedinFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := database.GetTodaysLinkedinFeeds()
	if err != nil {
		log.Printf("❌ DB Error: %v", err)
		http.Error(w, "Database Error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feeds)
}

func HandleGetLinkedinFeeds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", 405)
		return
	}
	identifier := r.PathValue("identifier")

	// Try to fetch by ID first
	feed, err := database.GetLinkedinFeedByID(identifier)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"feed": feed})
		return
	}

	feeds, err := database.GetLinkedinFeedsByProject(identifier)
	if err != nil {
		http.Error(w, "Database Retrieval Failed", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(feeds); err != nil {
		http.Error(w, "JSON Encoding Failed", 500)
	}
}
