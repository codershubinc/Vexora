package schema

var InstagramFeedDBSchema = `
CREATE TABLE IF NOT EXISTS instagram_feeds (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    feed TEXT, 
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	project_name TEXT
);`
