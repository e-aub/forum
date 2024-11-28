package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/internal/database"
	models "forum/internal/database/models"
	"html/template"
	"net/http"
	"os"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	switch r.Method {
	case "GET":
		category := r.URL.Query().Get("category")
		if category != "" {
			postIds, err := database.GetCategoryContentIds(db, category, userId)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			path := "./web/templates/"
			files := []string{
				path + "base.html",
				path + "pages/posts.html",
			}
			template, err := template.ParseFiles(files...)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			jsonIds, err := json.Marshal(postIds)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			feed := struct {
				Style string
				Posts string
			}{
				Style: "post.css",
				Posts: string(jsonIds),
			}
			template.ExecuteTemplate(w, "base", feed)
			return
		}
		withCreatedAndDeleted := false
		if userId != 0 {
			withCreatedAndDeleted = true
		}
		categories, err := GetCategories(db, withCreatedAndDeleted)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		path := "./web/templates/"
		files := []string{
			path + "base.html",
			path + "pages/categories.html",
		}
		template, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		feed := struct {
			Style      string
			Categories []models.Category
		}{
			Style:      "categories.css",
			Categories: categories,
		}
		template.ExecuteTemplate(w, "base", feed)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetCategories(db *sql.DB, withLikedAndCreated bool) ([]models.Category, error) {
	var result []models.Category
	var err error
	var rows *sql.Rows
	if withLikedAndCreated {
		rows, err = db.Query(`SELECT id, name, description FROM categories`)
	} else {
		rows, err = db.Query(`SELECT id, name, description FROM categories WHERE id != 1 AND id != 2`)
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var row models.Category
		if err := rows.Scan(&row.Id, &row.Name, &row.Description); err == nil {
			result = append(result, row)
		} else {
			return nil, err
		}
	}
	return result, nil
}
