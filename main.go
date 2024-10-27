package main

import (
	"time"

	"forum/handlers"
)

func main() {
	// server := http.Server{
	// 	Addr : "127.0.0.1:8080",
	// 	Handler : handlers.Mux,
	// }

	// handlers.Mux.HandleFunc("/", handlers.Home)

	// fmt.Println("serve has been launched at localhost:8080")
	// fmt.Println("http://localhost:8080")

	// err := server.ListenAndServe()
	// if err != nil {
	// 	fmt.Println("\nfatal:\n\tserver has been closed. port specified already on use")
	// 	return
	// }
	var post handlers.Posts
	post.Content = "we need to test it"
	post.Title = "testing"
	post.UserId = 10
	post.Created_At = time.Now().GoString()
	handlers.Creat_TAble_Posts()
	handlers.Insert_Post(&post)
}
