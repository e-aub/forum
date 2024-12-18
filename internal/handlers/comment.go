package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/internal/database"
	"forum/internal/utils"
)

func GetCommentsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	postID, err := strconv.Atoi(r.URL.Query().Get("post"))
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "Bad Request"})
		return
	}

	comments, err := database.GetComments(postID, db, userId)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, comments)
}

func AddCommentHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	userName, err := database.GetUserName(userId, db)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	comment := utils.Comment{}

	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}
	comment.Content = strings.TrimSpace(comment.Content)
	if len(comment.Content) < 1 || len(comment.Content) > 2000 {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "Comment must be between 3 and 2000 characters"})
		return
	}

	comment.User_name = userName
	comment.User_id = userId
	comment.Created_at = time.Now().Format(time.RFC3339)

	if err := database.CreateComment(&comment, db); err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, comment)
}
