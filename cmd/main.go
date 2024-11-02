package main

import (
	"log"
	"net/http"
	"os"

	"forum/internal/database"
	"forum/internal/handlers"

	_ "github.com/mattn/go-sqlite3"
)

func enableCors(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler(w, r)
	}
}

func main() {
	dbPath := os.Getenv("DB_PATH")
	db := database.CreateDatabase(dbPath)
	defer db.Close()
	database.CreateTables(db)

	mainMux := http.NewServeMux()

	fs := http.FileServer(http.Dir("assets"))
	mainMux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mainMux.HandleFunc("/", handlers.Controlle_Home)
	mainMux.HandleFunc("/New_Post", handlers.NewPostHandler)
	mainMux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.Register(db, w, r)
	})
	mainMux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.Login(db, w, r)
	})

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api", enableCors(handlers.Controlle_Api))

	go func() {
		log.Println("API server running on http://localhost:8000")
		if err := http.ListenAndServe(":8000", apiMux); err != nil {
			log.Fatalf("API server failed: %v", err)
		}
	}()

	log.Println("Route server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mainMux); err != nil {
		log.Fatalf("Route server failed: %v", err)
	}
}
