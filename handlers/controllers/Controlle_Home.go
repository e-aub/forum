package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"forum/handlers"
)

func Controlle_Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Println("Error loading template:", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, nil); err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
		}
		return
	}
	if r.Method == http.MethodPost {
		newpost := r.Body
		post := handlers.Creat_New_Post()
		err := json.NewDecoder(newpost).Decode(post)
		if err != nil {
			log.Fatal(err)
			return
		}
		handlers.Insert_Post(post)
		return
	}
}
