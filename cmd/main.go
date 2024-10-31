package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Colors = map[string]string{"green": "\033[42m", "red": "\033[41m", "reset": "\033[0m"}
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	port := os.Getenv("PORT")
	fmt.Println(port)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// Verify the connection
	if err = db.Ping(); err != nil {
		log.Fatalf("%sError accessing database: %s%s\n", Colors["red"], err.Error(), Colors["reset"])
	} else {
		fmt.Printf("%sDatabase created/opened successfully!%s\n", Colors["green"], Colors["reset"])
	}

	_, err = db.Exec(`PRAGMA foreign_keys=ON;`)
	if err != nil {
		log.Fatalf("%sError enabling foreign keys: %s%s\n", Colors["red"], err.Error(), Colors["reset"])
	}
	time.Sleep(time.Minute)
}
