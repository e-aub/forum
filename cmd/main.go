package main

import (
	"log"
	"net/http"
	"os"

	"forum/internal/database"
	"forum/internal/handlers"
	auth "forum/internal/middleware"
	"forum/internal/utils"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	port := "8080"
	print("port:", port)
	// Create Database file
	db := database.CreateDatabase(dbPath)
	defer db.Close()

	// Create tables if not exist
	database.CreateTables(db)

	database.CleanupExpiredSessions(db)

	// Create a multipluxer
	router := http.NewServeMux()

	// File Server (need some improvement by rearrange css and js files by separating them)
	fs := http.FileServer(http.Dir("web/assets"))
	router.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// HomePage handler
	router.HandleFunc("/", handlers.HomePageHandler)

	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.RegisterPageHandler(w, r)
		case "POST":
			handlers.RegisterHandler(w, r, db)
		default:
			utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.LoginPageHandler(w, r)
		case "POST":
			handlers.LoginHandler(w, r, db)
		default:
			utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		err := auth.RemoveUser(w, r, db)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	router.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
			return
		}
		userId, _ := auth.ValidUser(r, db)
		handlers.PostsHandler(w, r, db, userId)
	})

	router.HandleFunc("/new_post", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			auth.AuthMiddleware(db, handlers.NewPostPageHandler).ServeHTTP(w, r)
		case "POST":
			auth.AuthMiddleware(db, handlers.NewPostHandler).ServeHTTP(w, r)
		// We need to add UPDATE and DELETE methods to handle theses operations on posts
		default:
			utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
		}
	})

	// We need to add UPDATE and DELETE methods to handle theses operations on comments
	router.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			userId, _ := auth.ValidUser(r, db)
			handlers.GetCommentsHandler(w, r, db, userId)
		case "POST":
			auth.AuthMiddleware(db, handlers.AddCommentHandler).ServeHTTP(w, r)
		default:
			utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/react", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "GET" {
			userId, _ := auth.ValidUser(r, db)
			handlers.GetReactionsHandler(w, r, db, userId)
		} else if method == "PUT" {
			auth.AuthMiddleware(db, handlers.InsertOrUpdateReactionHandler).ServeHTTP(w, r)
		} else if method == "DELETE" {
			auth.AuthMiddleware(db, handlers.DeleteReactionHandler).ServeHTTP(w, r)
		} else {
			utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
			return
		}
	})

	router.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.CategoriesHandler(w, r, db)
		default:
			utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/me/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			auth.AuthMiddleware(db, handlers.MeHandler).ServeHTTP(w, r)
		default:
			utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
		}
	})

	log.Printf("Route server running on http://localhost:%s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, router))
}
