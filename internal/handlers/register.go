package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"forum/internal/database"
	middleware "forum/internal/middleware"
	"forum/internal/utils"
	tmpl "forum/web"
)

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "register", http.StatusOK, nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var userData utils.User
	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if len(userData.UserName) < 5 || len(userData.Password) < 8 || len(userData.UserName) > 30 || len(userData.Password) > 64 || !isValidEmail(&userData.Email) {
		http.Error(w, "invalid username/password", http.StatusNotAcceptable)
		return
	}
	ok, err := middleware.IsUserRegistered(db, &userData)
	if err != nil {
		http.Error(w, "internaInternal Server Error", http.StatusInternalServerError)
		return
	}
	if ok {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	err = HashPassword(&userData.Password)
	if err != nil {
		http.Error(w, "Invalid password", http.StatusNotAcceptable)
		return
	}
	err = middleware.RegisterUser(db, &userData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Create a session and set a cookie
	userData.SessionId, err = GenerateSessionID()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userData.Expiration = time.Now().Add(1 * time.Hour)
	err = database.InsertSession(db, &userData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Path:    "/",
		Value:   userData.SessionId,
		Expires: userData.Expiration,
	})
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}
