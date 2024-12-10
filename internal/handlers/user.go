package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"forum/internal/database"
	middleware "forum/internal/middleware"
	utils "forum/internal/utils"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

func MeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	switch r.URL.Path {
	case "/me/liked_posts":
		query := `SELECT post_id FROM reactions WHERE user_id = ? AND type_id = 'like'`
		rows, err := utils.QueryRows(db, query, userId)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.RespondWithError(w, utils.Err{Message: "Internal server Error"}, http.StatusInternalServerError)
			return
		}
		postIds := []int{}
		for rows.Next() {
			var postId int
			err := rows.Scan(&postId)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				utils.RespondWithError(w, utils.Err{Message: "Internal server Error"}, http.StatusInternalServerError)
			}
			postIds = append(postIds, postId)
		}
		path := "./web/templates/"
		files := []string{
			path + "base.html",
			path + "pages/posts.html",
		}
		template, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
		jsonIds, err := json.Marshal(postIds)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
		feed := struct {
			Style string
			Posts string
		}{
			Style: "post.css",
			Posts: string(jsonIds),
		}
		err = template.ExecuteTemplate(w, "base", feed)
		// err = template.Execute(w, string(jsonIds))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
	case "/me/created_posts":
		query := `SELECT id FROM posts WHERE user_id = ?`
		rows, err := utils.QueryRows(db, query, userId)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.RespondWithError(w, utils.Err{Message: "Internal server Error"}, http.StatusInternalServerError)
			return
		}
		postIds := []int{}
		for rows.Next() {
			var postId int
			err := rows.Scan(&postId)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				utils.RespondWithError(w, utils.Err{Message: "Internal server Error"}, http.StatusInternalServerError)
			}
			postIds = append(postIds, postId)
		}
		path := "./web/templates/"
		files := []string{
			path + "base.html",
			path + "pages/posts.html",
		}
		template, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
		jsonIds, err := json.Marshal(postIds)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
		feed := struct {
			Style string
			Posts string
		}{
			Style: "post.css",
			Posts: string(jsonIds),
		}
		err = template.ExecuteTemplate(w, "base", feed)
		// err = template.Execute(w, string(jsonIds))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
	default:
		utils.RespondWithError(w, utils.Err{Message: "404 Not found"}, http.StatusNotFound)
	}
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
	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
	}
	// return
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var userData utils.User
	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}
	// fmt.Println(userData.UserName)

	ok, err := middleware.IsUserRegistered(db, &userData)
	if err != nil {
		http.Error(w, "internaInternal Server Error", http.StatusInternalServerError)
		return
	}
	if ok {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	if len(userData.UserName) < 8 || len(userData.Password) < 8 || len(userData.UserName) > 30 || len(userData.Password) > 64 {
		http.Error(w, "invalid username/password", http.StatusNotAcceptable)
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
	var userData utils.User
	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	password := userData.Password
	ok, err := middleware.ValidCredential(db, &userData)
	// fmt.Println(userData.UserId)

	if !ok {
		http.Error(w, "Incorect Username or password", http.StatusConflict)
		return
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !CheckPasswordHash(&password, &userData.Password) {
		http.Error(w, "Incorrect Password", http.StatusConflict)
		return
	}
	ok, err = middleware.GetActiveSession(db, &userData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if ok {
		err = middleware.DeleteSession(db, &userData)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

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

func GenerateSessionID() (string, error) {
	sessionID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return sessionID.String(), nil
}

func HashPassword(password *string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), 14)
	*password = string(bytes)
	return err
}

func CheckPasswordHash(password, hash *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*hash), []byte(*password))
	return err == nil
}
