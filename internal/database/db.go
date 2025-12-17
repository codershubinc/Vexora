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
}

// SaveEntry stores the result
func SaveEntry(project, raw string, summary, polished, twitter, tags string) error {
	stmt, err := DB.Prepare(`
		INSERT INTO journal_entries(project_name, raw_notes, summary, polished_note, twitter_draft, tags) 
		VALUES(?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(project, raw, summary, polished, twitter, tags)
	return err
}

// JournalEntry represents a row in the database
type JournalEntry struct {
	ID           int
	ProjectName  string
	RawNotes     string
	Summary      string
	PolishedNote string
	TwitterDraft string
	Tags         string
	CreatedAt    string // Simplified for display
}

// GetEntries retrieves all journal entries ordered by date desc
func GetEntries() ([]JournalEntry, error) {
	rows, err := DB.Query(`
		SELECT id, project_name, raw_notes, summary, polished_note, twitter_draft, tags, created_at 
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
		// Note: SQLite stores dates as strings or numbers. We'll scan into string for simplicity.
		if err := rows.Scan(&e.ID, &e.ProjectName, &e.RawNotes, &e.Summary, &e.PolishedNote, &e.TwitterDraft, &e.Tags, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// GetEntryByID retrieves a single entry by ID
func GetEntryByID(id int) (*JournalEntry, error) {
	row := DB.QueryRow(`
		SELECT id, project_name, raw_notes, summary, polished_note, twitter_draft, tags, created_at 
		FROM journal_entries 
		WHERE id = ?
	`, id)

	var e JournalEntry
	if err := row.Scan(&e.ID, &e.ProjectName, &e.RawNotes, &e.Summary, &e.PolishedNote, &e.TwitterDraft, &e.Tags, &e.CreatedAt); err != nil {
		return nil, err
	}
	return &e, nil
}

// UpdateEntry updates the AI-generated fields for an entry
func UpdateEntry(id int, summary, polished, twitter, tags string) error {
	stmt, err := DB.Prepare(`
		UPDATE journal_entries 
		SET summary = ?, polished_note = ?, twitter_draft = ?, tags = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(summary, polished, twitter, tags, id)
	return err
}
func GetEntriesByProject(projectName string) ([]JournalEntry, error) {
	rows, err := DB.Query(`
		SELECT id, project_name, raw_notes, summary, polished_note, twitter_draft, tags, created_at 
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
		if err := rows.Scan(&e.ID, &e.ProjectName, &e.RawNotes, &e.Summary, &e.PolishedNote, &e.TwitterDraft, &e.Tags, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}
