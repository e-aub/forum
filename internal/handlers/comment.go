package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"forum/internal/database"
	"forum/internal/utils"
	util "forum/internal/utils"
)

func CommentsApiHandler(w http.ResponseWriter, r *http.Request, file *sql.DB, userId int) {
	if r.URL.Path != "/api/comments" {
		http.Error(w, "not found", 404)
	}

	switch r.Method {
	case "GET": //"http://localhost:8080/api/comments?post=${post.id}
		postID, _ := strconv.Atoi(r.URL.Query().Get("post"))
		comments, err := database.GetComments(postID, file, userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, comments)
	case "POST": //"http://localhost:8080/api/comments?post=${post.id}&comment=${comment.msj}"
		if userId > 0 {
			postID, _ := strconv.Atoi(r.URL.Query().Get("post"))
			content := r.URL.Query().Get("comment")
			user_name, err := database.GetUserName(userId, file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			comment := util.Creat_New_Comment()
			comment.Post_id = postID
			comment.User_name = user_name
			comment.User_id = userId
			comment.Content = content
			comment.Created_at = time.Now().Format(time.RFC3339)

			if err := database.CreateComment(comment, file); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			utils.RespondWithJSON(w, http.StatusCreated, comment)
		} else {
			utils.RespondWithJSON(w, http.StatusUnauthorized, "You are in guest mode try to login")
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
