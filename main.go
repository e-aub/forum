package main

import (
	"log"
	"net/http"

	"forum/internal/handlers"
)

func main() {
	server := http.NewServeMux()

	// Serve static files from the "templates" directory under the "/templates/" path


	// Define other routes
	server.HandleFunc("/", handlers.Controlle_Home)
	server.HandleFunc("/api", handlers.Controlle_Api)

	// Start the server
	log.Println("Server running on http://localhost:8000")
	if err := http.ListenAndServe(":8000", server); err != nil {
		log.Fatal(err)
	}
}
