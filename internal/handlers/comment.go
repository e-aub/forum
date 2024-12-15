package handlers

import (
	"database/sql"
	"encoding/json"
	"html"
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
	user_name, err := database.GetUserName(userId, file)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, `{"error":"Bad request"}`)
		return
	}
	comment := utils.Comment{}
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}
	if len(comment.Content) > 150 {
		http.Error(w, "length of comment over 150 character", http.StatusBadRequest)
		return
	}
	comment.Content = html.EscapeString(comment.Content)
	comment.User_name = user_name
	comment.User_id = userId
	comment.Created_at = time.Now().Format(time.RFC3339)

	if err := database.CreateComment(&comment, file); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, `{"error":"Bad request"}`)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, comment)
}
