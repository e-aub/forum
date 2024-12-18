package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"time"

	"forum/internal/database"
	"forum/internal/utils"
)

type querys struct {
	postid int
	limit  int
	from   int
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request, file *sql.DB, userId int) {
	var data querys
	err := getDataQuery(&data, r)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, fmt.Sprintf(`{"error":"status bad request %v"}`, err))
		return
	}

	comments, err := database.GetComments(data.postid, file, userId, data.limit, data.from)
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
	if len(comment.Content) > 2000 {
		http.Error(w, "length of comment over 2000 character", http.StatusBadRequest)
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

func getDataQuery(field *querys, r *http.Request) error {
	allKeys := []string{"post", "from", "limit"}
	for _, key := range allKeys {
		data, err := strconv.Atoi(r.URL.Query().Get(key))
		if err != nil {
			return errors.New("faild to get " + key + " value")
		}
		switch key {
		case "post":
			field.postid = data
		case "from":
			field.from = data
		case "limit":
			field.limit = data
		}
	}
	return nil
}
