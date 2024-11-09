package utils

import "time"

var Colors = map[string]string{"green": "\033[42m", "red": "\033[41m", "reset": "\033[0m"}

type Posts struct {
	PostId     int
	UserId     int
	UserName   string
	Title      string
	Category   string
	Content    string
	Created_At time.Time
}

type Comment struct {
	Comment_id int    `json:"comment_id"`
	Post_id    int    `json:"post_id"`
	User_id    int    `json:"user_id"`
	User_name  string `json:"user_name"`
	Content    string `json:"content"`
	Created_at string `json:"created_at"`
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
