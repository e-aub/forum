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
	"forum/internal/utils"
)

func Controlle_Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		utils.RespondWithError(w, utils.Err{Message: "404 page not found", Unauthorized: false}, http.StatusNotFound)
		return
	}
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("web/templates/posts.html")
		if err != nil {
			log.Println("Error loading template:", err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, false); err != nil {
			log.Println("Error executing template:", err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
		}
		return
	}
}

func NewPostHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	if userId <= 0 {
		utils.RespondWithError(w, utils.Err{Message: "Unauthorized: Please login and try again", Unauthorized: true}, http.StatusUnauthorized)
		return
	}
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("web/templates/newPost.html")
		if err != nil {
			log.Printf("Error parsing template: %v", err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
		categories, err := GetCategories(db, false)
		if err != nil {
			fmt.Println(err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, categories)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
			utils.RespondWithError(w, utils.Err{Message: "Bad request", Unauthorized: false}, http.StatusBadRequest)
			return
		}

		post := &utils.Posts{
			Title:      r.PostFormValue("title"),
			Content:    r.PostFormValue("content"),
			Created_At: time.Now(),
			UserId:     userId,
		}
		// category := r.PostFormValue("category")
		categories := r.Form["category"]
		_, err := database.InsertPost(post, db, categories)
		if err != nil {
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
}

func PostsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	if r.Method != "GET" {
		utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("content-type", "application/json")
	id := r.FormValue("id")
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
