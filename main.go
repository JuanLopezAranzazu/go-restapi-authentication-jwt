package main

import (
	"log"
	"net/http"

	"github.com/JuanLopezAranzazu/go-restapi-authentication-jwt/db"
	"github.com/JuanLopezAranzazu/go-restapi-authentication-jwt/routes"
	"github.com/gorilla/mux"
)

func main() {
	// conexion con la base de datos
	db.DBConnection()

	// manejar rutas
	r := mux.NewRouter()
	// index
	r.HandleFunc("/", routes.HomeHandler)

	// iniciar servidor
	log.Println("Servidor iniciado en http://localhost:3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
	}
}
