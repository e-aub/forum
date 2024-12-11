package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"forum/internal/database"
	"forum/internal/utils"
)

func GetCommentsHandler(w http.ResponseWriter, r *http.Request, file *sql.DB, userId int) {
	postID, _ := strconv.Atoi(r.URL.Query().Get("post"))
	comments, err := database.GetComments(postID, file, userId)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, `{"error":"internal server error"}`)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, comments)
}

func AddCommentHandler(w http.ResponseWriter, r *http.Request, file *sql.DB, userId int) {
	//"http://localhost:8080/comments?post=${post.id}&comment=${comment.msj}"
	postID, _ := strconv.Atoi(r.URL.Query().Get("post"))
	content := r.URL.Query().Get("comment")
	user_name, err := database.GetUserName(userId, file)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, `{"error":"Bad request"}`)
		return
	}
	comment := utils.Comment{}
	comment.Post_id = postID
	comment.User_name = user_name
	comment.User_id = userId
	comment.Content = content
	comment.Created_at = time.Now().Format(time.RFC3339)

	if err := database.CreateComment(&comment, file); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, `{"error":"Bad request"}`)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, comment)
}
