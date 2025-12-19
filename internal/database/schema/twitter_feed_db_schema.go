package schema

var TwitterFeedDBSchema = `
CREATE TABLE IF NOT EXISTS twitter_feeds (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    feed TEXT, 
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	project_name TEXT
);`
