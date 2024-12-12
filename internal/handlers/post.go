package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	tmpl "forum/web"

	database "forum/internal/database"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusNotFound, tmpl.Err{Message: "page not found"})
		return
	} else if r.Method != "GET" {
		tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusNotFound, tmpl.Err{Message: "page not found"})
		return
	}
	tmpl.ExecuteTemplate(w, []string{"posts", "sideBar"}, http.StatusOK, nil)
}

func PostsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
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
		post, err := database.ReadPost(db, userId, postId)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			return
		}
		post.Categories, err = database.GetPostCategories(db, post.PostId, userId)
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
	lastindex, err := database.GetLastPostId(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	json, err := json.Marshal(lastindex)
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
