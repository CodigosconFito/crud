package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasena := ""
	Nombre := "sistema"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasena+"@tcp(127.0.0.1)/"+Nombre)

	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	log.Println("Servidor corriendo...")
	http.ListenAndServe(":8080", nil)
}

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func Home(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionDB()
	registros, err := conexionEstablecida.Query("select * from empleados")

	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for registros.Next() {
		var id int
		var nombre, correo string

		err = registros.Scan(&id, &nombre, &correo)

		if err != nil {
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)
	}

	fmt.Println(arregloEmpleado)
	conexionEstablecida.Close()

	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado)
}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionDB()
		insertarRegistro, err := conexionEstablecida.Prepare("insert into empleados (nombre, correo) values (?,?)")

		if err != nil {
			panic(err.Error())
		}
		insertarRegistro.Exec(nombre, correo)
		log.Println("Datos insertados")
		http.Redirect(w, r, "/", 301)
	}
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")

	conexionEstablecida := conexionDB()
	insertarRegistro, err := conexionEstablecida.Prepare("delete from empleados where id=?")

	if err != nil {
		panic(err.Error())
	}
	insertarRegistro.Exec(idEmpleado)
	http.Redirect(w, r, "/", 301)

	fmt.Println("usuario eliminado")
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")

	conexionEstablecida := conexionDB()
	registro, err := conexionEstablecida.Query("select * from empleados where id = ?", idEmpleado)

	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{}

	for registro.Next() {
		var id int
		var nombre, correo string

		err = registro.Scan(&id, &nombre, &correo)

		if err != nil {
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
	}

	plantillas.ExecuteTemplate(w, "editar", empleado)
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		println("id " + id)
		println("nombre" + nombre)
		println("correo" + correo)

		conexionEstablecida := conexionDB()
		modificarRegistro, err := conexionEstablecida.Prepare("UPDATE empleados SET nombre=?, correo=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}
		modificarRegistro.Exec(nombre, correo, id)
		log.Println("Datos actualizados")
		http.Redirect(w, r, "/", 301)
	}
}
