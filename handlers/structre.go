package handlers

type Posts struct {
	PostId     float64
	UserId     float64
	Title      string
	Content    string
	Created_At string
	Updated_At string
}

func Creat_New_Post() *Posts {
	return &Posts{}
}

func (p *Posts) Update_Post(title string, content string, time string) {
	p.Title = title
	p.Content = content
	p.Updated_At = time
}
