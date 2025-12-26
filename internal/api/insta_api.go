package api

import (
	"encoding/json"
	"log"
	"net/http"
	"vexora-studio/internal/database"
	"vexora-studio/internal/llm"
)

func HandleCreateInstagramFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	body := r.FormValue("body")
	if body == "" {
		http.Error(w, "Content is required", 400)
		return
	}

	rawContent := r.FormValue("raw_content")
	projectName := r.FormValue("project_name")

	if rawContent == "" {
		http.Error(w, "Raw content is required", 400)
		return
	}

	data, err := llm.GenerateContent(llm.TypeInstagram, rawContent)
	if err != nil {
		log.Printf("❌ Instagram Generation Failed: %v", err)
		http.Error(w, "Content Generation Failed", 500)
		return
	}

	err = database.InsertInstagramFeed(data, projectName)
	if err != nil {
		log.Printf("❌ Instagram DB Insert Failed: %v", err)
		http.Error(w, "Database Insertion Failed", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))

}

func HandleGetTodaysInstagramFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := database.GetTodaysInstagramFeeds()
	if err != nil {
		log.Printf("❌ DB Error: %v", err)
		http.Error(w, "Database Error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feeds)
}

func HandleGetInstagramFeeds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", 405)
		return
	}
	identifier := r.PathValue("identifier")

	// Try to fetch by ID first if it's numeric
	feed, err := database.GetInstagramFeedByID(identifier)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"feed": feed})
		return
	}

	// Otherwise fetch by project name
	feeds, err := database.GetInstagramFeedsByProject(identifier)
	if err != nil {
		http.Error(w, "Database Retrieval Failed", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(feeds); err != nil {
		http.Error(w, "JSON Encoding Failed", 500)
	}
}
