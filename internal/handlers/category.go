package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/internal/database"
	models "forum/internal/database/models"
	"forum/internal/utils"
	"html/template"
	"net/http"
	"os"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	if r.URL.Query().Has("category") {
		category := r.URL.Query().Get("category")

		result, err := utils.QueryRow(db, `SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)`, category)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
		var exists bool
		if err := result.Scan(&exists); err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
		if !exists {
			utils.RespondWithError(w, utils.Err{Message: "404 page not found", Unauthorized: false}, http.StatusNotFound)
			return
		}
		postIds, err := database.GetCategoryContentIds(db, category, userId)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
		template, err := template.ParseFiles("web/templates/posts.html")
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
		template.Execute(w, string(jsonIds))
	} else {
		withCreatedAndDeleted := false
		if userId != 0 {
			withCreatedAndDeleted = true
		}
		categories, err := GetCategories(db, withCreatedAndDeleted)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
		template, err := template.ParseFiles("web/templates/categories.html")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
		template.Execute(w, categories)
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
