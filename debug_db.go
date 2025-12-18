package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./data/vexora.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT count(*) FROM journal_entries")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		rows.Scan(&count)
	}
	fmt.Printf("Total entries: %d\n", count)

	rows, err = db.Query("SELECT project_name FROM journal_entries")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Projects:")
	for rows.Next() {
		var name string
		rows.Scan(&name)
		fmt.Println("- " + name)
	}
}
