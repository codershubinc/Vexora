package database

func InsertInstagramFeed(feed string, projectName string) error {
	query := `INSERT INTO instagram_feeds (feed, project_name) VALUES (?, ?);`
	_, err := DB.Exec(query, feed, projectName)
	if err != nil {
		return err
	}
	return nil
}

func GetInstagramFeedsByProject(projectName string) ([]string, error) {
	query := `SELECT feed FROM instagram_feeds WHERE project_name = ?;`
	rows, err := DB.Query(query, projectName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []string
	for rows.Next() {
		var feed string
		if err := rows.Scan(&feed); err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return feeds, nil
}

func GetTodaysInstagramFeeds() ([]string, error) {
	query := `SELECT feed FROM instagram_feeds WHERE DATE(created_at) = DATE('now');`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []string
	for rows.Next() {
		var feed string
		if err := rows.Scan(&feed); err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return feeds, nil
}
