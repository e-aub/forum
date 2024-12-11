package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"forum/internal/database"
	models "forum/internal/database/models"
	"forum/internal/utils"
	tmpl "forum/web"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.URL.Query().Has("category") {
		category := r.URL.Query().Get("category")

		result, err := utils.QueryRow(db, `SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)`, category)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			tmpl.ExecuteTemplate(w, "error", http.StatusInternalServerError, tmpl.Err{Message: "internal server error"})
			return
		}
		var exists bool
		if err := result.Scan(&exists); err != nil {
			fmt.Fprintln(os.Stderr, err)
			tmpl.ExecuteTemplate(w, "error", http.StatusInternalServerError, tmpl.Err{Message: "internal server error"})
			return
		}
		if !exists {
			tmpl.ExecuteTemplate(w, "error", http.StatusNotFound, tmpl.Err{Message: "404 page not found"})
			return
		}
		postIds, err := database.GetCategoryContentIds(db, category)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			tmpl.ExecuteTemplate(w, "error", http.StatusInternalServerError, tmpl.Err{Message: "internal server error"})
			return
		}
		jsonIds, err := json.Marshal(postIds)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			tmpl.ExecuteTemplate(w, "error", http.StatusInternalServerError, tmpl.Err{Message: "internal server error"})
			return
		}
		tmpl.ExecuteTemplate(w, "posts", http.StatusOK, string(jsonIds))
	} else {
		categories, err := GetCategories(db)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			tmpl.ExecuteTemplate(w, "error", http.StatusInternalServerError, tmpl.Err{Message: "internal server error"})
			return
		}
		tmpl.ExecuteTemplate(w, "categories", http.StatusOK, categories)
	}
}

func GetCategories(db *sql.DB) ([]models.Category, error) {
	var result []models.Category
	var err error
	var rows *sql.Rows
	rows, err = utils.QueryRows(db, `SELECT id, name, description FROM categories`)
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
