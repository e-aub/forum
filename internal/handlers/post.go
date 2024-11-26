package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	database "forum/internal/database"
	"forum/internal/utils"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		utils.RespondWithError(w, utils.Err{Message: "404 page not found", Unauthorized: false}, http.StatusNotFound)
		return
	}
	if r.Method == http.MethodGet {
		path := "./web/templates/"
		files := []string{
			path + "base.html",
			path + "pages/posts.html",
		}
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Println("Error loading template:", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		feed := struct {
			Style string
			Posts bool
		}{
			Style: "post.css",
			Posts: false,
		}
		err = tmpl.ExecuteTemplate(w, "base", feed)
		if err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
		}
		return
	}
}

func NewPostPageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	tmpl, err := template.ParseFiles("web/templates/newPost.html")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
		return
	}
	categories, err := GetCategories(db, false)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, categories)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
		return
	}
}

func NewPostHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		utils.RespondWithError(w, utils.Err{Message: "Bad request", Unauthorized: false}, http.StatusBadRequest)
		return
	}

	post := &utils.Post{
		Title:      r.PostFormValue("title"),
		Content:    r.PostFormValue("content"),
		Created_At: time.Now(),
		UserId:     userId,
	}
	categories := r.Form["category"]
	_, err := database.InsertPost(post, db, categories)
	if err != nil {
		utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func PostsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	w.Header().Set("content-type", "application/json")
	id := r.URL.Query().Get("post_id")
	if id != "" {
		postId, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(nil)
			return
		}
		post, err := database.ReadPost(db, userId, postId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			return
		}
		post.Categories, err = database.GetPostCategories(db, post.PostId, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			return
		}
		json, err := json.Marshal(post)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			return
		}
		_, err = w.Write(json)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			return
		}
		return
	}
	lastindex, err := database.GetLastPostId(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	json, err := json.Marshal(lastindex)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	_, err = w.Write(json)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
}
