package main

import (
	"fmt"
	"forum/internal/database"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	port := os.Getenv("PORT")
	fmt.Println(port)
	db := database.CreateDatabase(dbPath)
	defer db.Close()
	database.CreateTables(db)
}
