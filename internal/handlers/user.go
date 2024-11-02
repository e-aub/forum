package handlers

import (
	"database/sql"
	"fmt"
	middleware "forum/internal/middleware"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Register(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/templates/register.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		if err := t.Execute(w, nil); err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	email := r.FormValue("email")
	if len(username) < 8 || len(password) < 8 {
		http.Error(w, "invalid username/password", http.StatusNotAcceptable)
		return
	}
	ok, err := middleware.IsUserRegistered(db, email, username)
	if err != nil {
		http.Error(w, "internaInternal Server Error", http.StatusInternalServerError)
		return
	}
	if ok {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	hachedPassword, err := HashPassword(password)
	if err != nil {
		http.Error(w, "Invalid password", http.StatusNotAcceptable)

	}
	err = middleware.RegisterUser(db, username, email, hachedPassword)
	if err != nil {
		http.Error(w, "internaInternal Server Error", http.StatusInternalServerError)
	}
	fmt.Fprintln(w, "user registred")
}

func Login(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/templates/login.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		if err := t.Execute(w, nil); err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
	username := r.FormValue("username")
	email := r.FormValue("email")

	password := r.FormValue("password")
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
		http.Error(w, "User not exists", http.StatusConflict)
		return
	}
	fmt.Fprintln(w, "logiin succefully")

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
