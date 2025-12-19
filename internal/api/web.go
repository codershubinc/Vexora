package api

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"vexora-studio/internal/database"
)

type PageData struct {
	Feeds       []string
	CurrentType string
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	feedType := r.URL.Query().Get("type")
	if feedType == "" {
		feedType = "instagram" // Default
	}

	var feeds []string
	var err error

	switch feedType {
	case "twitter":
		feeds, err = database.GetTodaysTwitterFeeds()
	case "linkedin":
		feeds, err = database.GetTodaysLinkedinFeeds()
	case "newsletter":
		feeds, err = database.GetTodaysNewsletterFeeds()
	case "instagram":
		fallthrough
	default:
		feeds, err = database.GetTodaysInstagramFeeds()
		feedType = "instagram"
	}

	if err != nil {
		log.Printf("❌ DB Error: %v", err)
		http.Error(w, "Database Error", 500)
		return
	}
	renderTemplate(w, PageData{Feeds: feeds, CurrentType: feedType})
}

func renderTemplate(w http.ResponseWriter, data PageData) {
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
