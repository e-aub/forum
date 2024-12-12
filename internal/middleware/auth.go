package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"forum/internal/database"
	"forum/internal/utils"
	tmpl "forum/web"
)

type customHandler func(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int)

func AuthMiddleware(db *sql.DB, next customHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, _ := ValidUser(r, db)
		if userId <= 0 {
			fmt.Println(r.Header.Get("Content-Type"))
			if r.Header.Get("Content-Type") == "application/json" {
				utils.RespondWithJSON(w, http.StatusUnauthorized, `{"error":"Unauthorized"}`)
				return
			}
			tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusUnauthorized, tmpl.Err{Message: "You are unauthorized, please log in", Unauthorized: true})
			return
		}
		next(w, r, db, userId)
	})
}

func IsUserRegistered(db *sql.DB, userData *utils.User) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR username = ?);`
	err := db.QueryRow(query, userData.Email, userData.UserName).Scan(&exists)
	return exists, err
}

func RegisterUser(db *sql.DB, userData *utils.User) error {
	insertQuery := `INSERT INTO users (username, email, password) VALUES (?, ?, ?);`
	result, err := db.Exec(insertQuery, userData.UserName, userData.Email, userData.Password)
	if err != nil {
		return err
	}
	userData.UserId, err = result.LastInsertId()
	return err
}

func GetActiveSession(db *sql.DB, userData *utils.User) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM sessions WHERE user_id = ?  AND expires_at > ?);`
	err := db.QueryRow(query, userData.UserId, userData.Expiration).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func DeleteSession(db *sql.DB, userData *utils.User) error {
	query := `DELETE FROM sessions WHERE user_id =  ?;`
	_, err := db.Exec(query, userData.UserId)
	return err
}

func ValidCredential(db *sql.DB, userData *utils.User) (bool, error) {
	query := `SELECT id, password FROM users WHERE username = ?;`
	err := db.QueryRow(query, userData.UserName).Scan(&userData.UserId, &userData.Password)
	if err != nil {
		return false, err
	}
	return true, err
}

func ValidUser(r *http.Request, db *sql.DB) (int, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err != http.ErrNoCookie {
			return 0, err
		}
		return 0, nil
	}
	userid, err := database.Get_session(cookie.Value, db)
	if err != nil {
		return 0, err
	}
	return userid, nil
}

func RemoveUser(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",              // Clear the value
		Expires: time.Unix(0, 0), // Expire the cookie immediately
		// Path:     "/",             // Match original path if specified
		// HttpOnly: true,
		// Secure:   true,
		// SameSite: http.SameSiteStrictMode, // Match original SameSite attribute
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
