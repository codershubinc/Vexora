package api

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"vexora-studio/internal/database"
	"vexora-studio/internal/llm"
)

// HandleMiscPage renders the form for manual entry
func HandleMiscPage(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("templates", "misc.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("❌ Template Error: %v", err)
		http.Error(w, "Template Error", 500)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("❌ Render Error: %v", err)
	}
}

// HandleMiscSubmit processes the manual entry
func HandleMiscSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	rawText := r.FormValue("raw_text")
	if rawText == "" {
		http.Error(w, "Content is required", 400)
		return
	}

	projectName := "Miscellaneous"
	log.Printf("📝 Processing Miscellaneous Entry...")

	startTime := time.Now()

	// Call LLM
	journal, err := llm.MessageLlama(projectName, rawText)
	if err != nil {
		log.Printf("❌ AI Error: %v", err)
		http.Error(w, "AI Processing Failed", 500)
		return
	}

	duration := time.Since(startTime)
	log.Printf("⏱️ Processing took %v", duration)

	// Flatten tags slice to string
	tagsStr := strings.Join(journal.Tags, ",")

	// Save to DB
	err = database.SaveEntry(
		projectName,
		rawText,
		journal.Summary,
		journal.PolishedNote,
		journal.TwitterDraft,
		journal.LinkedinDraft,
		journal.InstagramCaption,
		tagsStr,
		duration.String(),
	)

	if err != nil {
		log.Printf("❌ DB Save Error: %v", err)
		http.Error(w, "Database Save Failed", 500)
		return
	}

	log.Printf("✅ Miscellaneous Entry Saved Successfully")

	// Redirect to the project view for Miscellaneous
	http.Redirect(w, r, "/project/Miscellaneous", http.StatusSeeOther)
}
