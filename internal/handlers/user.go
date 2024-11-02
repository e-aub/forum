package handlers

import (
	"database/sql"
	"encoding/json"
	middleware "forum/internal/middleware"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var requestData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
	t, err := template.ParseFiles("web/templates/register.html")
	if err != nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, nil); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}

}
func Register_Api(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}
	username := requestData.Username
	email := requestData.Email
	password := requestData.Password
	ok, err := middleware.IsUserRegistered(db, email, username)
	if err != nil {
		http.Error(w, "internaInternal Server Error", http.StatusInternalServerError)
		return
	}
	if ok {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	if len(username) < 8 || len(password) < 8 || len(username) > 30 || len(password) > 64 {
		http.Error(w, "invalid username/password", http.StatusNotAcceptable)
		return
	}

	hachedPassword, err := HashPassword(password)
	if err != nil {
		http.Error(w, "Invalid password", http.StatusNotAcceptable)
		return
	}
	err = middleware.RegisterUser(db, username, email, hachedPassword)
	if err != nil {
		http.Error(w, "internaInternal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
	t, err := template.ParseFiles("web/templates/login.html")
	if err != nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, nil); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func Login_Api(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	username := requestData.Username
	email := requestData.Email
	password := requestData.Password

	ok, err := middleware.IsUserRegistered(db, email, username)
	if err != nil {
		http.Error(w, "internaInternal Server Error", http.StatusInternalServerError)
		return
	}
	hachedPassword, err := middleware.GetPasswordByUsername(db, username)
	if err != nil {
		http.Error(w, "internaInternal Server Error", http.StatusInternalServerError)
		return
	}
	if !ok || !CheckPasswordHash(password, hachedPassword) {
		http.Error(w, "Incorrect password or Username", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
