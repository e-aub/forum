package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"forum/internal/database"
	"forum/internal/types"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var user types.User
	db := database.Init()
	database.CreateUsersTable(db)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	exists, err := database.IsUserRegistered(db, user.Email, user.Username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	if err := database.RegisterUser(db, user.Username, user.Email, user.Password); err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
