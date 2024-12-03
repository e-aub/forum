package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"forum/internal/database"
	models "forum/internal/database/models"
	"forum/internal/utils"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
		postIds, err := database.GetCategoryContentIds(db, category)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
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
		//err = template.Execute(w, string(jsonIds))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
	} else {
		categories, err := GetCategories(db)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
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
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
		feed := struct {
			Style      string
			Categories []models.Category
		}{
			Style:      "categories.css",
			Categories: categories,
		}
		err = template.ExecuteTemplate(w, "base", feed)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.RespondWithError(w, utils.Err{Message: "internal server error", Unauthorized: false}, http.StatusInternalServerError)
			return
		}
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
