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
//testing
func main() {
	dbPath := os.Getenv("DB_PATH")
	port := os.Getenv("PORT")
	db := database.CreateDatabase(dbPath)
	defer db.Close()
	database.CreateTables(db)
	mainMux := http.NewServeMux()
	///////////////////FOR FILE JS /////////////////////
	fs := http.FileServer(http.Dir("web/assets"))
	mainMux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	////////////////ROUTES////////////////////////////
	mainMux.HandleFunc("/", handlers.Controlle_Home)
	mainMux.HandleFunc("/New_Post", func(w http.ResponseWriter, r *http.Request) {
		userId, is_user := auth.ValidUser(w, r, db)
		if !is_user {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		handlers.NewPostHandler(w, r, userId, db)
	})
	mainMux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.Register(w, r)
	})
	mainMux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.Login(w, r)
	})

	///////////////API////////////////////
	mainMux.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.Controlle_Api(w, r, db)
	})
	mainMux.HandleFunc("/api/comments", func(w http.ResponseWriter, r *http.Request) {
		user_id, is_user := auth.ValidUser(w, r, db)
		handlers.Controlle_Api_Comment(w, r, user_id, is_user, db)
	})

	mainMux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.Register_Api(db, w, r)
	})
	mainMux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.Login_Api(db, w, r)
	})
	mainMux.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		auth.RemoveUser(w, r, db)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Printf("Route server running on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mainMux); err != nil {
		log.Fatalf("Route server failed: %v", err)
	}
}
