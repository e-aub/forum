package database

import (
	"database/sql"
	"log"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

func CreateUsersTable(db *sql.DB) {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );
    `
	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}
func Init() *sql.DB {
	db, err := sql.Open("sqlite3", "internal/database/test.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func IsUserRegistered(db *sql.DB, email, username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR username = ?);`
	err := db.QueryRow(query, email, username).Scan(&exists)
	return exists, err
}

// Register a new user in the database
func RegisterUser(db *sql.DB, username, email, password string) error {
	insertQuery := `INSERT INTO users (username, email, password) VALUES (?, ?, ?);`
	_, err := db.Exec(insertQuery, username, email, password)
	return err
}
