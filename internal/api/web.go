package api

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"vexora-studio/internal/database"
)

type PageData struct {
	Entries        []database.JournalEntry
	CurrentProject string
}

func HandleWebIndex(w http.ResponseWriter, r *http.Request) {
	entries, err := database.GetEntries()
	if err != nil {
		log.Printf("❌ DB Error: %v", err)
		http.Error(w, "Database Error", 500)
		return
	}
	renderTemplate(w, PageData{Entries: entries})
}

func HandleWebProject(w http.ResponseWriter, r *http.Request) {
	project := r.PathValue("project")
	entries, err := database.GetEntriesByProject(project)
	if err != nil {
		log.Printf("❌ DB Error: %v", err)
		http.Error(w, "Database Error", 500)
		return
	}
	renderTemplate(w, PageData{Entries: entries, CurrentProject: project})
}

func renderTemplate(w http.ResponseWriter, data PageData) {
	log.Printf("Rendering template with %d entries. Project: %s", len(data.Entries), data.CurrentProject)
	tmplPath := filepath.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("❌ Template Error: %v", err)
		http.Error(w, "Template Error", 500)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("❌ Render Error: %v", err)
	}
}
