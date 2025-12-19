package database

import (
	"database/sql"
	"log"
	"os"
	"time"
	"vexora-studio/internal/database/schema"
)

var DB *sql.DB

func Init(path string) error {
	if err := os.Mkdir("data", 0o755); err != nil && !os.IsExist(err) {
		return err
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(20 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Println("Database ping failed, %w", err)
		return err
	}
	DB = db

	// creating the tables in db

	if _, err := DB.Exec(schema.InstagramFeedDBSchema); err != nil {
		return err
	}

	if _, err := DB.Exec(schema.NewsletterDBSchema); err != nil {
		return err
	}

	if _, err := DB.Exec(schema.LinkedinFeedDBSchema); err != nil {
		return err
	}

	if _, err := DB.Exec(schema.TwitterFeedDBSchema); err != nil {
		return err
	}
	return nil

}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func GetDB() *sql.DB {
	return DB
}
