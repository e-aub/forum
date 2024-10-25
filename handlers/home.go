package handlers

import "net/http"
import "fmt"

var Mux = http.NewServeMux()


func Home(w http.ResponseWriter, r  *http.Request){
	fmt.Fprintln(w, "welcome to forum project")
}
