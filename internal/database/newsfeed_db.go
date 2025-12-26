package database

func InsertNewsletterFeed(feed string, projectName string) error {
	query := `INSERT INTO newsletters (feed , project_name) VALUES (?, ?);`
	_, err := DB.Exec(query, feed, projectName)
	if err != nil {
		return err
	}
	return nil
}

func GetNewsletterFeeds(projectName string) ([]string, error) {
	query := `SELECT feed FROM newsletters WHERE project_name = ?;`
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

func GetNewsletterByID(id string) (string, error) {
	query := `SELECT feed FROM newsletters WHERE id = ?;`
	var feed string
	err := DB.QueryRow(query, id).Scan(&feed)
	if err != nil {
		return "", err
	}
	return feed, nil
}

func GetTodaysNewsletterFeeds() ([]string, error) {
	query := `SELECT feed FROM newsletters WHERE DATE(created_at) = DATE('now');`
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
