package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	database "forum/internal/database"

	util "forum/internal/utils"
)

func Controlle_Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	// fmt.Println(r.Cookies())

	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("web/templates/posts.html")
		if err != nil {
			log.Println("Error loading template:", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, false); err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
		}
		return
	}
}

func NewPostHandler(w http.ResponseWriter, r *http.Request, userId int, db *sql.DB) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("web/templates/creat_Post.html")
		if err != nil {
			log.Printf("Error parsing template: %v", err)
			http.Error(w, "Internal Server Errorr", http.StatusInternalServerError)
			return
		}
		categories, err := GetCategories(db, false)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "get categories", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, categories)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		post := &util.Posts{
			Title:      r.PostFormValue("title"),
			Content:    r.PostFormValue("content"),
			Created_At: time.Now(),
			UserId:     userId,
		}
		// category := r.PostFormValue("category")
		categories := r.Form["category"]
		_, err := database.Insert_Post(post, db, categories)
		if err != nil {
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func Controlle_Api(w http.ResponseWriter, r *http.Request, file *sql.DB, userId int) {
	if r.URL.Path != "/api/posts" {
		http.Error(w, "not found", 404)
	}
	if r.Method != "GET" {
		http.Error(w, "method Not allowed", http.StatusMethodNotAllowed)
	}
	id := r.FormValue("id")
	if id != "" {
		idint, _ := strconv.Atoi(id)
		post := database.Read_Post(idint, file)
		post.Categories, _ = database.GetPostCategories(post.PostId, file, userId)
		json, err := json.Marshal(post)
		if err != nil {
			log.Fatal(err)
		}
		_, _ = w.Write(json)
		return
	}
	lastindex := database.Get_Last(file)
	json, err := json.Marshal(lastindex)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = w.Write(json)
}
