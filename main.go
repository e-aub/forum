package main 

import "fmt"
import "net/http"
import "forum/handlers"

func main(){

	server := http.Server{
		Addr : "127.0.0.1:8080",
		Handler : handlers.Mux,
	}
	
	handlers.Mux.HandleFunc("/", handlers.Home)

	fmt.Println("serve has been launched at localhost:8080")
	fmt.Println("http://localhost:8080")

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("\nfatal:\n\tserver has been closed. port specified already on use")
		return
	}
	
}
