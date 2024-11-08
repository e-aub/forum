package main

import (
	"log"
	"net/http"
	"os"

	"forum/internal/database"
	"forum/internal/handlers"
	auth "forum/internal/middleware"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	port := os.Getenv("PORT")
	db := database.CreateDatabase(dbPath)
	defer db.Close()
	database.CreateTables(db)
	database.CleanupExpiredSessions(db)
	mainMux := http.NewServeMux()
	///////////////////FOR FILE JS /////////////////////
	fs := http.FileServer(http.Dir("web/assets"))
	mainMux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	////////////////ROUTES////////////////////////////
	mainMux.HandleFunc("/", handlers.Controlle_Home)
	mainMux.HandleFunc("/New_Post", func(w http.ResponseWriter, r *http.Request) {
		ok, userID, err := auth.ValidUser(w, r, db)
		if err != nil {
			// here adding template for eroor
			return
		}
		if ok {
			handlers.NewPostHandler(w, r, userID)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
	mainMux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.Register(w, r)
	})
	mainMux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		ok, _, err := auth.ValidUser(w, r, db)
		if err != nil {
			// here adding template for eroor
			return
		}
		if ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		handlers.Login(w, r)
	})
	///////////////API////////////////////
	mainMux.HandleFunc("/api/posts", handlers.Controlle_Api)
	mainMux.HandleFunc("/api/comments", func(w http.ResponseWriter, r *http.Request) {
		ok, userID, err := auth.ValidUser(w, r, db)
		if err != nil {
			return
		}
		handlers.Controlle_Api_Comment(w, r, userID, ok)
		// http.Redirect(w, r, "/login", http.StatusSeeOther)

	})

	mainMux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.Register_Api(db, w, r)
	})
	mainMux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.Login_Api(db, w, r)

	})
	mainMux.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		err := auth.RemoveUser(w, r, db)
		if err != nil {
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Printf("Route server running on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mainMux); err != nil {
		log.Fatalf("Route server failed: %v", err)
	}
}
