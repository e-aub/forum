package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Comment struct {
	Comment_id int    `json:"comment_id"`
	Post_id    int    `json:"post_id"`
	User_id    int    `json:"user_id"`
	Content    string `json:"content"`
	Created_at string `json:"created_at"`
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewComment() *Comment {
	return &Comment{}
}

func getDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitializeDataBase() error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `
	CREATE TABLE IF NOT EXISTS comments (
		comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at TEXT NOT NULL
	);
	`
	_, err = db.Exec(query)
	return err
}

func (s *Server) HandleComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")

	w.Header().Set("Content-Type", "application/json")

	if err := InitializeDataBase(); err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Database initialization failed")
		return
	}

	switch r.Method {
	case "POST": //"http://localhost:8080/api/comments?post=${post.id}&user=${user.id}&comment=${comment.msj}"
		postID, _ := strconv.Atoi(r.URL.Query().Get("post"))
		userID, _ := strconv.Atoi(r.URL.Query().Get("user"))
		content := r.URL.Query().Get("comment")

		comment := NewComment()
		comment.Post_id = postID
		comment.User_id = userID
		comment.Content = content
		comment.Created_at = time.Now().Format(time.RFC3339)

		if err := comment.CreateComment(); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondWithJSON(w, http.StatusCreated, ApiResponse{
			Success: true,
			Message: "Comment created successfully",
			Data:    comment,
		})

	case "GET": //"http://localhost:8080/api/comments
		comments, err := GetCommentsByPost()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, ApiResponse{
			Success: true,
			Data:    comments,
		})

	case "DELETE": //"http://localhost:8080/api/comments?comment=${comment.id}&user=${user.id}"
		fmt.Println("detext")
		commentID, _ := strconv.Atoi(r.URL.Query().Get("comment"))
		userID, _ := strconv.Atoi(r.URL.Query().Get("user"))
		if err := DeleteComment(commentID, userID); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, ApiResponse{
			Success: true,
			Message: "Comment deleted successfully",
		})

	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ApiResponse{
		Success: false,
		Message: message,
	})
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

func (c *Comment) CreateComment() error {
	if c.Content == "" || c.Post_id == 0 || c.User_id == 0 {
		return errors.New("comment DATA issue")
	}

	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `
	INSERT INTO comments (post_id, user_id, content, created_at)
	VALUES (?, ?, ?, ?)
	`

	result, err := db.Exec(query, c.Post_id, c.User_id, c.Content, c.Created_at)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	c.Comment_id = int(id)

	return nil
}

func GetComment(commentID int) (Comment, error) {
	db, err := getDB()
	if err != nil {
		return Comment{}, err
	}
	defer db.Close()

	query := `
	SELECT comment_id, post_id, user_id, content, created_at 
	FROM comments
	WHERE comment_id = ?
	`

	var comment Comment
	err = db.QueryRow(query, commentID).Scan(
		&comment.Comment_id,
		&comment.Post_id,
		&comment.User_id,
		&comment.Content,
		&comment.Created_at,
	)
	if err == sql.ErrNoRows {
		return Comment{}, errors.New("comment not found")
	}
	if err != nil {
		return Comment{}, err
	}

	return comment, nil
}

func GetCommentsByPost() ([]Comment, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
	SELECT comment_id, post_id, user_id, content, created_at 
	FROM comments
	ORDER BY created_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.Comment_id,
			&comment.Post_id,
			&comment.User_id,
			&comment.Content,
			&comment.Created_at,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func DeleteComment(commentID int, userID int) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `
	SELECT user_id FROM comments 
	WHERE comment_id = ?
	`
	var dbUserID int
	err = db.QueryRow(query, commentID).Scan(&dbUserID)
	if err == sql.ErrNoRows {
		return errors.New("comment not found")
	}
	if err != nil {
		return err
	}

	if dbUserID != userID {
		return errors.New("unauthorized: only comment owner can delete")
	}

	query = `DELETE FROM comments WHERE comment_id = ?`
	result, err := db.Exec(query, commentID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("comment not found")
	}

	return nil
}
