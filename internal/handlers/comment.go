package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	db "forum/internal/database"
	util "forum/internal/utils"
)

func Controlle_Api_Comment(w http.ResponseWriter, r *http.Request, user_id float64, valided bool) {
	if r.URL.Path != "/api/comments" {
		http.Error(w, "not found", 404)
	}

	switch r.Method {
	case "GET": //"http://localhost:8080/api/comments?post=${post.id}
		postID, _ := strconv.Atoi(r.URL.Query().Get("post"))
		comments, err := db.GetComments(postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusOK, comments)
	case "POST": //"http://localhost:8080/api/comments?post=${post.id}&comment=${comment.msj}"
		if valided {
			postID, _ := strconv.Atoi(r.URL.Query().Get("post"))
			content := r.URL.Query().Get("comment")
			user_name, err := db.GetUserName(int(user_id))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			comment := util.Creat_New_Comment()
			comment.Post_id = postID
			comment.User_name = user_name
			comment.User_id = int(user_id)
			comment.Content = content
			comment.Created_at = time.Now().Format(time.RFC3339)

			if err := db.CreateComment(comment); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			respondWithJSON(w, http.StatusCreated, comment)
		} else {
			respondWithJSON(w, http.StatusUnauthorized, "You are in guest mode try to login")
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(code)
	w.Write(response)
}
