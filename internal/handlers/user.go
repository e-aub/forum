package handlers

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"forum/internal/database"
	middleware "forum/internal/middleware"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

var requestData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	path := "./web/templates/"
	files := []string{
		path + "base.html",
		path + "pages/register.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Println("Error loading template:", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
	feed := struct {
		Style string
	}{
		Style: "register.css",
	}
	err = tmpl.ExecuteTemplate(w, "base", feed)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
	}
	return
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
	userID, err := database.GetUserIDByUsername(db, username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Create a session and set a cookie
	sessionID, err := GenerateSessionID()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	expiration := time.Now().Add(1 * time.Hour)
	err = database.InsertSession(db, sessionID, userID, expiration)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Path:    "/",
		Value:   sessionID,
		Expires: expiration,
	})
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	path := "./web/templates/"
	files := []string{
		path + "base.html",
		path + "pages/login.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Println("Error loading template:", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
	feed := struct {
		Style string
	}{
		Style: "login.css",
	}
	err = tmpl.ExecuteTemplate(w, "base", feed)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	username := requestData.Username
	email := requestData.Email
	password := requestData.Password

	ok, err := middleware.IsUserRegistered(db, email, username)
	if !ok {
		http.Error(w, "Incorect Username", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "internaInternal Server Error", http.StatusInternalServerError)
		return
	}
	hachedPassword, err := middleware.GetPasswordByUsername(db, username)
	if err != nil {
		http.Error(w, "internaInternal Server Error", http.StatusInternalServerError)
		return
	}
	if !CheckPasswordHash(password, hachedPassword) {
		http.Error(w, "Incorrect Password", http.StatusConflict)
		return
	}
	userID, err := database.GetUserIDByUsername(db, username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Create a session and set a cookie
	sessionID, err := GenerateSessionID()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	expiration := time.Now().Add(1 * time.Hour)
	err = database.InsertSession(db, sessionID, userID, expiration)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "session_token",
		Path:  "/",
		Value: sessionID,
	})

	w.WriteHeader(http.StatusOK)
}

func GenerateSessionID() (string, error) {
	sessionID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return sessionID.String(), nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
