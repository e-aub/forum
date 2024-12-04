package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	database "forum/internal/database"
	models "forum/internal/database/models"
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

func NewPostPageHandler(w http.ResponseWriter, r *http.Request, conn *database.Conn_db, userId int) {
	path := "./web/templates/"
	files := []string{
		path + "base.html",
		path + "pages/new_post.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
		return
	}

	categories, err := GetCategories(conn.Db)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
		return
	}
	feed := struct {
		Style      string
		Categories []models.Category
	}{
		Style:      "new_post.css",
		Categories: categories,
	}
	err = tmpl.ExecuteTemplate(w, "base", feed)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
		return
	}
}

func NewPostHandler(w http.ResponseWriter, r *http.Request, conn *database.Conn_db, userId int) {
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		utils.RespondWithError(w, utils.Err{Message: "Bad request", Unauthorized: false}, http.StatusBadRequest)
		return
	}

	post := &utils.Post{
		Title:      r.PostFormValue("title"),
		Content:    r.PostFormValue("content"),
		CreatedAt:  time.Now(),
		UserId:     userId,
		Categories: r.Form["category"],
	}

	if len(post.Title) >= 40 || len(post.Content) >= 300 {
		log.Printf("long format")
		utils.RespondWithError(w, utils.Err{Message: "bad request", Unauthorized: false}, http.StatusBadRequest)
		return
	}
	_, err := conn.InsertPost(post)
	if err != nil {
		utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func PostsHandler(w http.ResponseWriter, r *http.Request, conn *database.Conn_db, userId int) {
	w.Header().Set("content-type", "application/json")
	id := r.URL.Query().Get("post_id")
	if id != "" {
		postId, err := strconv.Atoi(id)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(nil)
			return
		}
		post, err := conn.ReadPost(postId)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			return
		}
		post.Categories, err = conn.GetPostCategories(post.PostId, userId)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			return
		}
		json, err := json.Marshal(post)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			return
		}
		_, err = w.Write(json)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			return
		}
		return
	}
	post, err := conn.ReadPost()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	json, err := json.Marshal(post.PostId)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	_, err = w.Write(json)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
}
