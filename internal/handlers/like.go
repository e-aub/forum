package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"forum/internal/database"
)

func Controlle_Api_likes(w http.ResponseWriter, r *http.Request, user_id int, valided bool, file *sql.DB) {
	if r.URL.Path != "/api/likes" {
		http.Error(w, "not found", 404)
	}

	switch r.Method {
	case "GET": //"http://localhost:8080/api/likes?postId=${}"
		postID, _ := strconv.Atoi(r.URL.Query().Get("postId"))
		likes, err := database.GetLikes(postID, file)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusOK, likes)
	case "POST": //"http://localhost:8080/api/likes?postId=${}&type=${}"
		if valided {
			postID, _ := strconv.Atoi(r.URL.Query().Get("postId"))
			type_like := r.URL.Query().Get("type")
			if err := database.CreateLike(postID, user_id, type_like, file); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			respondWithJSON(w, http.StatusCreated, "like posted succesfuly")
		} else {
			respondWithJSON(w, http.StatusUnauthorized, "You are in guest mode try to login")
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
