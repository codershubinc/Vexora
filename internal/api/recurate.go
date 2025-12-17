package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"vexora-studio/internal/database"
	"vexora-studio/internal/llm"
)

func HandleRecurate(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	// 1. Fetch Original Entry
	entry, err := database.GetEntryByID(id)
	if err != nil {
		log.Printf("❌ DB Error: %v", err)
		http.Error(w, "Entry not found", 404)
		return
	}

	log.Printf("🔄 Recurating entry %d for project %s...", id, entry.ProjectName)

	// 2. Call LLM Again
	journal, err := llm.MessageLlama(entry.ProjectName, entry.RawNotes)
	if err != nil {
		log.Printf("❌ AI Error: %v", err)
		http.Error(w, "AI Processing Failed", 500)
		return
	}

	// Flatten tags slice to string
	tagsStr := strings.Join(journal.Tags, ",")

	// 3. Update DB
	err = database.UpdateEntry(
		id,
		journal.Summary,
		journal.PolishedNote,
		journal.TwitterDraft,
		tagsStr,
	)

	if err != nil {
		log.Printf("❌ DB Update Error: %v", err)
		http.Error(w, "Database Update Failed", 500)
		return
	}

	log.Printf("✅ Entry %d Recurated Successfully", id)
	
	// Redirect back to where they came from (or home)
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func HandleFork(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	// 1. Fetch Original Entry
	entry, err := database.GetEntryByID(id)
	if err != nil {
		log.Printf("❌ DB Error: %v", err)
		http.Error(w, "Entry not found", 404)
		return
	}

	log.Printf("🔀 Forking entry %d for project %s...", id, entry.ProjectName)

	// 2. Call LLM Again (Generate new version)
	journal, err := llm.MessageLlama(entry.ProjectName, entry.RawNotes)
	if err != nil {
		log.Printf("❌ AI Error: %v", err)
		http.Error(w, "AI Processing Failed", 500)
		return
	}

	// Flatten tags slice to string
	tagsStr := strings.Join(journal.Tags, ",")

	// 3. Save as NEW Entry
	err = database.SaveEntry(
		entry.ProjectName,
		entry.RawNotes,
		journal.Summary,
		journal.PolishedNote,
		journal.TwitterDraft,
		tagsStr,
	)

	if err != nil {
		log.Printf("❌ DB Save Error: %v", err)
		http.Error(w, "Database Save Failed", 500)
		return
	}

	log.Printf("✅ Entry %d Forked Successfully", id)

	// Redirect back to where they came from (or home)
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}
