package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"forum/internal/types"

)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/users", s.RegisterHandler)

	mux.HandleFunc("/health", s.healthHandler)

	return s.corsMiddleware(mux)
}

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method )

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var user types.User
	fmt.Println(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	exists, err := s.db.IsUserRegistered( user.Email, user.Username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	if err := s.db.RegisterUser( user.Username, user.Email, user.Password); err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        // If it's an OPTIONS request, return without calling the next handler
        if r.Method == http.MethodOptions {
            return
        }

        // Call the next handler
        next.ServeHTTP(w, r)
    })
}