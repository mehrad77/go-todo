package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// init the SQLite database and creates tables if they don't exist
func Initialize() {
	var err error
	DB, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Create tables
	createTables()
}

func createTables() {
	userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
    );
    `
	_, err := DB.Exec(userTable)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	todoTable := `
    CREATE TABLE IF NOT EXISTS todos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        title TEXT NOT NULL,
        completed BOOLEAN NOT NULL,
        FOREIGN KEY(user_id) REFERENCES users(id)
    );
    `
	_, err = DB.Exec(todoTable)
	if err != nil {
		log.Fatalf("Error creating todos table: %v", err)
	}
}
