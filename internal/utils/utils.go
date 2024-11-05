package utils

import "time"

var Colors = map[string]string{"green": "\033[42m", "red": "\033[41m", "reset": "\033[0m"}

type Posts struct {
	PostId     float64
	UserId     float64
	Title      string
	Content    string
	Created_At time.Time
}

func Creat_New_Post() *Posts {
	return &Posts{}
}

func (p *Posts) Update_Post(title string, content string, time time.Time) {
	p.Title = title
	p.Content = content
}
