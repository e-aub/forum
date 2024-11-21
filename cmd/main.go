package main

import (
	"fmt"
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
	fmt.Println("db path:", dbPath)
	db := database.CreateDatabase(dbPath)
	defer db.Close()
	database.CreateTables(db)
	database.CleanupExpiredSessions(db)
	router := http.NewServeMux()
	///////////////////FOR FILE JS /////////////////////
	fs := http.FileServer(http.Dir("web/assets"))
	router.Handle("/assets/", http.StripPrefix("/assets/", fs))
	////////////////ROUTES////////////////////////////
	router.HandleFunc("/", handlers.Controlle_Home)
	router.Handle("/categories", auth.AuthMiddleware(db)(handlers.CategoriesHandler))
	router.Handle("/New_Post", auth.AuthMiddleware(db)(handlers.NewPostHandler))
	router.HandleFunc("/register", handlers.Register)

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.ValidUser(r, db)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		handlers.Login(w, r)
	})
	///////////////API////////////////////
	router.Handle("/api/posts", auth.AuthMiddleware(db)(handlers.PostsHandler))
	router.Handle("/api/comments", auth.AuthMiddleware(db)(handlers.CommentsApiHandler))

	router.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.Register_Api(w, r, db)
	})
	router.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.Login_Api(db, w, r)
	})
	router.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		err := auth.RemoveUser(w, r, db)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	router.Handle("/api/react/", auth.AuthMiddleware(db)(handlers.ReactHandler))
	log.Printf("Route server running on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Route server failed: %v", err)
	}
}
