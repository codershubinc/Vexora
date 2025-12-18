package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init(dbPath string) {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// NEW SCHEMA: Matches the JSON output from Llama
	schema := `
	CREATE TABLE IF NOT EXISTS journal_entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_name TEXT,
		raw_notes TEXT,
		
		-- AI Generated Fields
		summary TEXT,
		polished_note TEXT,
		twitter_draft TEXT,
		tags TEXT, -- Store as comma-separated string

		
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = DB.Exec(schema)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Migration: Add new columns if they don't exist
	// We ignore errors here because if the column exists, it will error, which is fine.
	DB.Exec("ALTER TABLE journal_entries ADD COLUMN linkedin_draft TEXT")
	DB.Exec("ALTER TABLE journal_entries ADD COLUMN instagram_caption TEXT")
	DB.Exec("ALTER TABLE journal_entries ADD COLUMN processing_time TEXT")
}

// SaveEntry stores the result
func SaveEntry(project, raw string, summary, polished, twitter, linkedin, instagram, tags, processingTime string) error {
	stmt, err := DB.Prepare(`
		INSERT INTO journal_entries(project_name, raw_notes, summary, polished_note, twitter_draft, linkedin_draft, instagram_caption, tags, processing_time) 
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(project, raw, summary, polished, twitter, linkedin, instagram, tags, processingTime)
	return err
}

// JournalEntry represents a row in the database
type JournalEntry struct {
	ID               int
	ProjectName      string
	RawNotes         string
	Summary          string
	PolishedNote     string
	TwitterDraft     string
	LinkedinDraft    string
	InstagramCaption string
	Tags             string
	ProcessingTime   string
	CreatedAt        string // Simplified for display
}

// GetEntries retrieves all journal entries ordered by date desc
func GetEntries() ([]JournalEntry, error) {
	rows, err := DB.Query(`
		SELECT id, project_name, raw_notes, summary, polished_note, twitter_draft, linkedin_draft, instagram_caption, tags, processing_time, created_at 
		FROM journal_entries 
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []JournalEntry
	for rows.Next() {
		var e JournalEntry
		// Handle NULLs for new columns if necessary, but Scan handles them if we use sql.NullString or just string (if not null).
		// Since we added them as TEXT, they might be NULL for old entries.
		// We should probably use sql.NullString or just accept that Scan might fail if we don't handle it.
		// For simplicity, let's assume they are empty strings if null (coalesce in query is safer).
		var linkedin, instagram, procTime sql.NullString

		if err := rows.Scan(&e.ID, &e.ProjectName, &e.RawNotes, &e.Summary, &e.PolishedNote, &e.TwitterDraft, &linkedin, &instagram, &e.Tags, &procTime, &e.CreatedAt); err != nil {
			return nil, err
		}
		e.LinkedinDraft = linkedin.String
		e.InstagramCaption = instagram.String
		e.ProcessingTime = procTime.String
		entries = append(entries, e)
	}
	return entries, nil
}

// GetEntryByID retrieves a single entry by ID
func GetEntryByID(id int) (*JournalEntry, error) {
	row := DB.QueryRow(`
		SELECT id, project_name, raw_notes, summary, polished_note, twitter_draft, linkedin_draft, instagram_caption, tags, processing_time, created_at 
		FROM journal_entries 
		WHERE id = ?
	`, id)

	var e JournalEntry
	var linkedin, instagram, procTime sql.NullString
	if err := row.Scan(&e.ID, &e.ProjectName, &e.RawNotes, &e.Summary, &e.PolishedNote, &e.TwitterDraft, &linkedin, &instagram, &e.Tags, &procTime, &e.CreatedAt); err != nil {
		return nil, err
	}
	e.LinkedinDraft = linkedin.String
	e.InstagramCaption = instagram.String
	e.ProcessingTime = procTime.String
	return &e, nil
}

// UpdateEntry updates the AI-generated fields for an entry
func UpdateEntry(id int, summary, polished, twitter, linkedin, instagram, tags, processingTime string) error {
	stmt, err := DB.Prepare(`
		UPDATE journal_entries 
		SET summary = ?, polished_note = ?, twitter_draft = ?, linkedin_draft = ?, instagram_caption = ?, tags = ?, processing_time = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(summary, polished, twitter, linkedin, instagram, tags, processingTime, id)
	return err
}
func GetEntriesByProject(projectName string) ([]JournalEntry, error) {
	rows, err := DB.Query(`
		SELECT id, project_name, raw_notes, summary, polished_note, twitter_draft, linkedin_draft, instagram_caption, tags, processing_time, created_at 
		FROM journal_entries 
		WHERE project_name = ?
		ORDER BY created_at DESC
	`, projectName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []JournalEntry
	for rows.Next() {
		var e JournalEntry
		var linkedin, instagram, procTime sql.NullString
		if err := rows.Scan(&e.ID, &e.ProjectName, &e.RawNotes, &e.Summary, &e.PolishedNote, &e.TwitterDraft, &linkedin, &instagram, &e.Tags, &procTime, &e.CreatedAt); err != nil {
			return nil, err
		}
		e.LinkedinDraft = linkedin.String
		e.InstagramCaption = instagram.String
		e.ProcessingTime = procTime.String
		entries = append(entries, e)
	}
	return entries, nil
}
