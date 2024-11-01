package main

import (
	"log"
	"net/http"

	"forum/handlers/controllers"
)

func main() {
	server := http.NewServeMux()

	// Serve static files from the "templates" directory under the "/templates/" path
	fs := http.FileServer(http.Dir("templates"))
	server.Handle("/templates/", http.StripPrefix("/templates/", fs))

	// Define other routes
	server.HandleFunc("/", controllers.Controlle_Home)
	server.HandleFunc("/api", controllers.Controlle_Api)

	// Start the server
	log.Println("Server running on http://localhost:8000")
	if err := http.ListenAndServe(":8000", server); err != nil {
		log.Fatal(err)
	}
}
