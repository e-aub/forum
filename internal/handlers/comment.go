package handlers

import (
	"net/http"
	"strconv"
	"time"

	"forum/internal/database"
	"forum/internal/utils"
)

func GetCommentsHandler(w http.ResponseWriter, r *http.Request, conn *database.Conn_db, userId int) {
	postID, _ := strconv.Atoi(r.URL.Query().Get("post"))
	comments, err := conn.GetComments(postID, userId)
	if err != nil {
		utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, comments)
}

func AddCommentHandler(w http.ResponseWriter, r *http.Request, conn *database.Conn_db, userId int) {
	//"http://localhost:8080/comments?post=${post.id}&comment=${comment.msj}"
	postID, _ := strconv.Atoi(r.URL.Query().Get("post"))
	content := r.URL.Query().Get("comment")
	post, err := conn.ReadPost(userId)
	if err != nil {
		utils.RespondWithError(w, utils.Err{Message: "Bad request", Unauthorized: false}, http.StatusBadRequest)
		return
	}
	comment := utils.Comment{}
	comment.Post_id = postID
	comment.User_name = post.UserName
	comment.User_id = userId
	comment.Content = content
	comment.Created_at = time.Now().Format(time.RFC3339)

	if err := conn.CreateComment(&comment); err != nil {
		utils.RespondWithError(w, utils.Err{Message: "Bad request", Unauthorized: false}, http.StatusBadRequest)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, comment)
}
