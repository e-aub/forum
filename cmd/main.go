package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	fmt.Println(dbPath)
	// port := os.Getenv("PORT")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// Verify the connection
	if err = db.Ping(); err != nil {
		fmt.Printf("\033[41mError accessing database: %s\033[0m\n", err.Error())
		return
	} else {
		fmt.Println("\033[42mDatabase created/opened successfully!\033[0m")
	}
	time.Sleep(time.Minute)
}
