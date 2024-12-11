package tmpl

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

type Err struct {
	Message      string
	Unauthorized bool
}

func ExecuteTemplate(w http.ResponseWriter, templateName string, statusCode int, data any) {
	basePath := "./web/templates/"

	templateFiles := []string{
		filepath.Join(basePath, "base.html"),
		filepath.Join(basePath, templateName+".html"),
	}

	tmpl, err := template.ParseFiles(templateFiles...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	fmt.Println(data)
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
