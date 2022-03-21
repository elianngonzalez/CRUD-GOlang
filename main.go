package main

import (
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var Port = ":3000"
var templates = template.Must(template.ParseGlob("templates/*"))

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/crear", Crear)
	log.Println("Listening on port" + Port)
	http.ListenAndServe(Port, nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello World")
	templates.ExecuteTemplate(w, "index", nil)
}

func Crear(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello World")
	templates.ExecuteTemplate(w, "crear", nil)
}
