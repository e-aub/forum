package main

import (
	"log"
	"net/http"
	"os"
	"time"

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
	connection_db := database.CreateDatabase(dbPath)
	defer connection_db.Db.Close()

	// Create tables if not exist
	connection_db.CreateTables()

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
			handlers.RegisterHandler(w, r, connection_db)
		default:
			utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.LoginPageHandler(w, r)
		case "POST":
			handlers.LoginHandler(w, r, connection_db)
		default:
			utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		err := auth.RemoveUser(w, r, connection_db.Db)
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
		userId, _ := auth.ValidUser(r, connection_db)
		handlers.PostsHandler(w, r, connection_db, userId)
	})

	router.HandleFunc("/new_post", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			auth.AuthMiddleware(connection_db, handlers.NewPostPageHandler).ServeHTTP(w, r)
		case "POST":
			auth.AuthMiddleware(connection_db, handlers.NewPostHandler).ServeHTTP(w, r)
		// We need to add UPDATE and DELETE methods to handle theses operations on posts
		default:
			utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
		}
	})

	// We need to add UPDATE and DELETE methods to handle theses operations on comments
	router.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			userId, _ := auth.ValidUser(r, connection_db)
			handlers.GetCommentsHandler(w, r, connection_db, userId)
		case "POST":
			auth.AuthMiddleware(connection_db, handlers.AddCommentHandler).ServeHTTP(w, r)
		default:
			utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/react", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "GET" {
			userId, _ := auth.ValidUser(r, connection_db)
			handlers.GetReactionsHandler(w, r, connection_db.Db, userId)
		} else if method == "PUT" {
			auth.AuthMiddleware(connection_db, handlers.InsertOrUpdateReactionHandler).ServeHTTP(w, r)
		} else if method == "DELETE" {
			auth.AuthMiddleware(connection_db, handlers.DeleteReactionHandler).ServeHTTP(w, r)
		} else {
			utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
			return
		}
	})

	router.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.CategoriesHandler(w, r, connection_db)
		default:
			utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/me/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			auth.AuthMiddleware(connection_db, handlers.MeHandler).ServeHTTP(w, r)
		default:
			utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
		}
	})

	go func() {
		for {
			database.CleanupExpiredSessions(connection_db.Db)
			time.Sleep(2 * time.Hour)
		}
	}()

	log.Printf("Route server running on http://localhost:%s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, router))
}
