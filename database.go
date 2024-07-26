package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error
	fmt.Println("Initializing database...")
	db, err = sql.Open("sqlite3", "todo.db")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database initialized")

	fmt.Println("Creating table...")
	createTable := `
	CREATE TABLE IF NOT EXISTS todos(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT FALSE
	)`
	fmt.Println("Table created")

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}
