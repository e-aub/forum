package main

import (
	"fmt"
	"net/http"

	"forum/internal/server"
)

func main() {
	mux := server.RegisterRoutes()
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
