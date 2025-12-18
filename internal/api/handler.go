package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	"vexora-studio/internal/database"
	"vexora-studio/internal/llm"
)

type Payload struct {
	ProjectName string `json:"project_name"`
	Diff        string `json:"diff"` // This actually contains the MD notes now
}

func HandleIngest(w http.ResponseWriter, r *http.Request) {
	var p Payload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Bad JSON", 400)
		return
	}

	log.Printf("📥 Received notes for %s. Processing...", p.ProjectName)

	startTime := time.Now()

	// Call Llama (Synchronous for now is fine, or go async if you prefer)
	journal, err := llm.MessageLlama(p.ProjectName, p.Diff)
	if err != nil {
		log.Printf("❌ AI Error: %v", err)
		http.Error(w, "AI Processing Failed", 500)
		return
	}

	duration := time.Since(startTime)

	// Flatten tags slice to string
	tagsStr := strings.Join(journal.Tags, ",")

	// Save to DB
	err = database.SaveEntry(
		p.ProjectName,
		p.Diff,
		journal.Summary,
		journal.PolishedNote,
		journal.TwitterDraft,
		journal.LinkedinDraft,
		journal.InstagramCaption,
		tagsStr,
		duration.String(),
	)

	if err != nil {
		log.Printf("❌ DB Error: %v", err)
		http.Error(w, "Database Error", 500)
		return
	}

	log.Printf("✅ Journal Saved: %s", journal.Summary)
	w.WriteHeader(200)
	w.Write([]byte("Journal Entry Created"))
}
