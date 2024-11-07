package main

import (
	"database/sql"
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
	mainMux := http.NewServeMux()
	///////////////////FOR FILE JS /////////////////////
	fs := http.FileServer(http.Dir("web/assets"))
	mainMux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	////////////////ROUTES////////////////////////////
	mainMux.HandleFunc("/", handlers.Controlle_Home)
	mainMux.HandleFunc("/New_Post", func(w http.ResponseWriter, r *http.Request) {
		userId, err := auth.ValidUser(w, r, db)
		if err != nil {
			if err == http.ErrNoCookie {
				handlers.Login(w, r)
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		handlers.NewPostHandler(w, r, userId)
	})
	mainMux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.Register(w, r)
	})
	mainMux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				handlers.Login(w, r)
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		_, err = database.Get_session(cookie.Value)
		if err != nil {
			if err == sql.ErrNoRows {
				handlers.Login(w, r)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	///////////////API////////////////////
	mainMux.HandleFunc("/api/posts", handlers.Controlle_Api)
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
