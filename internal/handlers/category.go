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
	switch r.Method {
	case "GET":
		category := r.URL.Query().Get("category")
		if category != "" {
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
			return
		}
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
	default:
		utils.RespondWithError(w, utils.Err{Message: "Method not allowed", Unauthorized: false}, http.StatusMethodNotAllowed)
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
