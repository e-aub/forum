package handlers

type Posts struct {
	UserId     int
	PostId     int
	Title      string
	Content    string
	Created_at string
	UpDated_at string
}

func (*Posts) Creat_New_Post() *Posts {
	return &Posts{}
}

func (p *Posts) Update_Post(newcontent string, time string, Title string) {
	p.Content = newcontent
	p.UpDated_at = time
	p.Title = Title
}
