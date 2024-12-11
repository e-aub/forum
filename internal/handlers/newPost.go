package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"forum/internal/database"
	"forum/internal/utils"
	tmpl "forum/web"
)

func NewPostPageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	categories, err := GetCategories(db)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		tmpl.ExecuteTemplate(w, "error", http.StatusInternalServerError, tmpl.Err{Message: "internal server error"})
		return
	}
	tmpl.ExecuteTemplate(w, "new_post", http.StatusOK, categories)
}

func NewPostHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		tmpl.ExecuteTemplate(w, "error", http.StatusBadRequest, tmpl.Err{Message: "Bad request"})
		return
	}

	post := &utils.Post{
		Title:     r.PostFormValue("title"),
		Content:   r.PostFormValue("content"),
		CreatedAt: time.Now(),
		UserId:    userId,
	}
	categories := r.Form["category"]

	if len(post.Title) >= 40 || len(post.Content) >= 300 {
		log.Printf("long format")
		tmpl.ExecuteTemplate(w, "error", http.StatusBadRequest, tmpl.Err{Message: "bad request"})
		return
	}
	_, err := database.InsertPost(post, db, categories)
	if err != nil {
		tmpl.ExecuteTemplate(w, "error", http.StatusInternalServerError, tmpl.Err{Message: "internal server error"})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
