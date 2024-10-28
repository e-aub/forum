package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"forum/handlers"
)

func Controlle_Api(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api" {
		http.Error(w, "not found", 404)
	}
	if r.Method != "GET" {
		http.Error(w, "method Not allowed", 405)
	}
	id := r.FormValue("id")
	if id != "" {
		idint, _ := strconv.Atoi(id)
		post := handlers.Read_Post(idint)
		json, err := json.Marshal(post)
		if err != nil {
			log.Fatal(err)
		}
		_, _ = w.Write(json)
		return
	}
	lastindex := handlers.Get_Last()
	json, err := json.Marshal(lastindex)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = w.Write(json)
}
