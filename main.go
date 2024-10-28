package main

import (
	"fmt"
	"time"

	"forum/handlers"
)

func main() {
	var post handlers.Posts
	post.Content = "wa l update"
	post.Title = "khdmiiiii"
	post.UserId = 10
	post.PostId = 6
	post.Updated_At = time.April.String()
	handlers.Creat_TAble_Posts()
	// handlers.Insert_Post(&post)
	// handlers.Update_Post(&post)
	var post2 handlers.Posts
	post2.PostId = 7
	handlers.Delete_Post(&post2)
	 handlers.Read_Post(1)
}
