package schema

var LinkedinFeedDBSchema = `
CREATE TABLE IF NOT EXISTS linkedin_feeds (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	feed TEXT,
	project_name TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`
