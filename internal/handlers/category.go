package handlers

import (
	"database/sql"
	"errors"
	models "forum/internal/database/models"
	"html/template"
	"net/http"
	"strings"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "GET":
		category, err := CategoriesUrlParser(r.URL.Path)
		if err != nil {
			http.Error(w, "NOT FOUND", http.StatusNotFound)
			return
		}
		if category != "" {
			return
		}
		categories, err := GetCategories(db)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		template, err := template.ParseFiles("web/templates/categories.html")
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		template.Execute(w, categories)
	default:
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
	}
}

func GetCategories(db *sql.DB) ([]models.Category, error) {
	var result []models.Category
	rows, err := db.Query(`SELECT id, name, description FROM categories`)
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
func CategoriesUrlParser(url string) (string, error) {
	parts := strings.Split(strings.Trim(url, "/"), "/")
	if len(parts) == 2 {
		return parts[1], nil
	} else if len(parts) == 1 && parts[0] == "categories" {
		return "", nil
	}
	return "", errors.New("NOT FOUND")
}
