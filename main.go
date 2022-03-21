package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexDB() (conexion *sql.DB) {
	driver := "mysql"
	user := "root"
	pass := ""
	dbName := "go_empleados"

	conexion, err := sql.Open(driver, user+":"+pass+"@tcp(127.0.0.1)/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var Port = ":3000"
var templates = template.Must(template.ParseGlob("templates/*"))

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/editar/{id}", Editar)

	log.Println("Listening on port" + Port)
	http.ListenAndServe(Port, nil)
}

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func Index(w http.ResponseWriter, r *http.Request) {
	conEstablecida := conexDB()
	consulta, err := conEstablecida.Query("SELECT * FROM empleados")
	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}
	empleadosLista := []Empleado{}

	for consulta.Next() {
		var id int
		var nombre, correo string

		err = consulta.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())

		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		empleadosLista = append(empleadosLista, empleado)
	}

	log.Println(empleadosLista)

	templates.ExecuteTemplate(w, "index", empleadosLista)
}

func Crear(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello World")
	templates.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")
		conEstablecida := conexDB()
		conRegistros, err := conEstablecida.Prepare("INSERT INTO empleados(nombre, correo) VALUES(?,?)")

		if err != nil {
			panic(err.Error())
		}
		conRegistros.Exec(nombre, correo)
		http.Redirect(w, r, "/", 301)
	}
}

func Editar(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	nombre := r.FormValue("nombre")
	correo := r.FormValue("correo")
	conEstablecida := conexDB()

	if r.Method == "get" {
		consulta, err := conEstablecida.Query("SELECT * FROM empleados WHERE id=?", id)
		if err != nil {
			panic(err.Error())
		}

		empleado := Empleado{}

		for consulta.Next() {
			var id int
			var nombre, correo string

			err = consulta.Scan(&id, &nombre, &correo)
			if err != nil {
				panic(err.Error())
			}

			empleado.Id = id
			empleado.Nombre = nombre
			empleado.Correo = correo
		}

		if r.Method == "post" {
			conRegistros, err := conEstablecida.Prepare("UPDATE empleados SET nombre=?, correo=? WHERE id=?")
			if err != nil {
				panic(err.Error())
			}
			conRegistros.Exec(nombre, correo, id)
			http.Redirect(w, r, "/", 301)
		}

		templates.ExecuteTemplate(w, "editar", empleado)
	}

}
