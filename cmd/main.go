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

	mainMux.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		_, userId, err := auth.ValidUser(w, r, db)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		handlers.CategoriesHandler(w, r, db, userId)
	})

	mainMux.HandleFunc("/New_Post", func(w http.ResponseWriter, r *http.Request) {
		ok, userID, err := auth.ValidUser(w, r, db)
		if err != nil {
			// here adding template for eroor
			return
		}
		if ok {
			handlers.NewPostHandler(w, r, userID, db)
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
	mainMux.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		_, userID, err := auth.ValidUser(w, r, db)
		if err != nil {
			return
		}

		handlers.Controlle_Api(w, r, db, userID)
	})
	mainMux.HandleFunc("/api/comments", func(w http.ResponseWriter, r *http.Request) {
		ok, userID, err := auth.ValidUser(w, r, db)
		if err != nil {
			return
		}
		handlers.Controlle_Api_Comment(w, r, userID, ok, db)
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
