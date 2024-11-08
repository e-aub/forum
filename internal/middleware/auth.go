package auth

import (
	"database/sql"
	"net/http"
	"time"

	"forum/internal/database"
)

func IsUserRegistered(db *sql.DB, email, username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR username = ?);`
	err := db.QueryRow(query, email, username).Scan(&exists)
	return exists, err
}

// Register a new user in the database
func RegisterUser(db *sql.DB, username, email, password string) error {
	insertQuery := `INSERT INTO users (username, email, password) VALUES (?, ?, ?);`
	_, err := db.Exec(insertQuery, username, email, password)
	return err
}

func GetPasswordByUsername(db *sql.DB, username string) (string, error) {
	var password string
	query := `SELECT password FROM users WHERE username = ?;`
	err := db.QueryRow(query, username).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func ValidUser(w http.ResponseWriter, r *http.Request, db *sql.DB) (bool, int, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false, 0, nil
	}
	userid, err := database.Get_session(db, cookie.Value)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return false, 0, err
	}
	return true, userid, nil
}

func RemoveUser(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",              // Clear the value
		Expires:  time.Unix(0, 0), // Expire the cookie immediately
		Path:     "/",             // Match original path if specified
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode, // Match original SameSite attribute
	})
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return err
	}
	stmt, err := db.Prepare("DELETE FROM sessions WHERE session_id = ?")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cookie.Value)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err

	}
	return nil
}
