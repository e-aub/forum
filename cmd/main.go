package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	// port := os.Getenv("PORT")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	db.Ping()
	fmt.Println("ddd")
	time.Sleep(time.Minute)
}
