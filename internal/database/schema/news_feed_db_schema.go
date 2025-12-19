package schema

var NewsletterDBSchema = `
CREATE TABLE IF NOT EXISTS newsletters (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	feed  TEXT,
	project_name TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP

);`
