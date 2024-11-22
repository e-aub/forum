package utils

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"
)

var Colors = map[string]string{"green": "\033[42m", "red": "\033[41m", "reset": "\033[0m"}

type Err struct {
	Message      string
	Unauthorized bool
}
type Posts struct {
	PostId       int
	UserId       int
	UserName     string
	Title        string
	Categories   []string
	Content      string
	LikeCount    int
	DislikeCount int
	Created_At   time.Time
	Clicked      bool
	DisClicked   bool
}

type Comment struct {
	Comment_id   int    `json:"comment_id"`
	Post_id      int    `json:"post_id"`
	User_id      int    `json:"user_id"`
	User_name    string `json:"user_name"`
	Content      string `json:"content"`
	Created_at   string `json:"created_at"`
	LikeCount    int    `json:"like_count"`
	DislikeCount int    `json:"dislike_count"`
	Clicked      bool   `json:"clicked"`
	DisClicked   bool   `json:"disclicked"`
}

func Creat_New_Post() *Posts {
	return &Posts{}
}

func Creat_New_Comment() *Comment {
	return &Comment{}
}

func (p *Posts) Update_Post(title string, content string, time time.Time) {
	p.Title = title
	p.Content = content
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, Err Err, statuscode int) {
	tmpl, err := template.ParseFiles("web/templates/error.html")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statuscode)
	tmpl.Execute(w, Err)
}


