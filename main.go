package main

import(
	"fmt"
	"net/http"
	"log"
	"text/template"
	_"github.com/go-sql-driver/mysql"
	"database/sql"
)

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

type Usuario struct{
	Id int
	Nombre string
	Correo string
}

func main(){
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/update", Update)
	log.Println("Servidor Trabajando...")
	http.ListenAndServe(":3000", nil)
}

func conexionDB()(conexion *sql.DB){
	Driver:= "mysql"
	Usuario:= "root"
	Contra:= "1234"
	Nombre:= "crudgo"

	conexion, err:= sql.Open(Driver, Usuario+":"+Contra+"@tcp(127.0.0.1)/"+Nombre)
	if err!=nil {
		panic(err.Error())
	}
	return conexion
}

func Update(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" {
		id:= r.FormValue("txtId")
		nombre:= r.FormValue("txtNombre")
		correo:= r.FormValue("txtCorreo")
		
		conexion:= conexionDB()
		modificar, err:= conexion.Prepare("UPDATE usuario SET nombre = ?, correo =? WHERE id = ?")
		if err != nil{
			panic(err.Error())
		}
		modificar.Exec(nombre, correo, id)
		http.Redirect(w, r, "/", 301)
	}
}

func Editar(w http.ResponseWriter, r *http.Request){
	idUsuario:= r.URL.Query().Get("id")
	conexion:= conexionDB()
	editar, err:= conexion.Query("SELECT * FROM usuario WHERE id = ?", idUsuario)
	if err!=nil {
		panic(err.Error())
	}
	usuario:= Usuario{}
	for editar.Next(){
		var id int
		var nombre, correo string
		err = editar.Scan(&id, &nombre, &correo)
		if err!=nil {
			panic(err.Error())
		}
		usuario.Id = id
		usuario.Nombre = nombre
		usuario.Correo = correo
	}
	plantillas.ExecuteTemplate(w, "editar", usuario)
}

func Borrar(w http.ResponseWriter, r *http.Request){
	idUsuario:= r.URL.Query().Get("id")
	conexion:= conexionDB()
	borrar, err:= conexion.Prepare("DELETE FROM usuario WHERE id = ?")
	if err!=nil {
		panic(err.Error())
	}
	borrar.Exec(idUsuario)
	http.Redirect(w, r, "/", 301)
}

func Insertar(w http.ResponseWriter, r *http.Request){
	//fmt.Fprintf(w, "Hola DICIS")
	if r.Method == "POST" {
		nombre:= r.FormValue("txtNombre")
		correo:= r.FormValue("txtCorreo")
		
		conexion:= conexionDB()
		insertar, err:= conexion.Prepare("INSERT INTO usuario(nombre, correo) VALUES(?, ?)")
		if err != nil{
			panic(err.Error())
		}
		insertar.Exec(nombre, correo)
		http.Redirect(w, r, "/", 301)
	}
}

func Inicio(w http.ResponseWriter, r *http.Request){
	//fmt.Fprintf(w, "Hola DICIS")
	conexion:= conexionDB()
	usuario:= Usuario{}
	arregloUsuario:= []Usuario{}
	seleccionar, err:= conexion.Query("SELECT * FROM usuario")
	if err != nil{
		panic(err.Error())
	}
	for seleccionar.Next(){
		var id int
		var nombre, correo string
		err = seleccionar.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		usuario.Id = id
		usuario.Nombre = nombre
		usuario.Correo = correo
		arregloUsuario = append(arregloUsuario, usuario)
	}
	fmt.Println(arregloUsuario)
	//seleccionar.Exec()
	plantillas.ExecuteTemplate(w, "inicio", arregloUsuario)
}

func Crear(w http.ResponseWriter, r *http.Request){
	//fmt.Fprintf(w, "Hola DICIS")
	plantillas.ExecuteTemplate(w, "crear", nil)
}
