package handlers

import (
	"database/sql"
	"fmt"
	"forum/internal/database"
	models "forum/internal/database/models"
	"html/template"
	"net/http"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "GET":
		fmt.Println(r.URL.String())
		category := r.URL.Query().Get("category")
		println(category)
		if category != "" {
			postIds, err := database.GetCategoryContentIds(db, category)
			if err != nil {
				http.Error(w, "internal server error1", http.StatusInternalServerError)
				return
			}
			posts, err := database.GetCategoryContent(db, category)
			if err != nil {
				http.Error(w, "internal server error1", http.StatusInternalServerError)
				return
			}

			w.Header().Add("content-type", "application/json")
			fmt.Fprintf(w, "%v\n%v", posts, postIds)
			return
		}
		categories, err := GetCategories(db, true)
		if err != nil {
			http.Error(w, "internal server error2", http.StatusInternalServerError)
			return
		}
		template, err := template.ParseFiles("web/templates/categories.html")
		if err != nil {
			http.Error(w, "internal server error3", http.StatusInternalServerError)
			return
		}
		template.Execute(w, categories)
	default:
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
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
