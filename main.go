package main

import (
	"log"
	"net/http"

	"github.com/JuanLopezAranzazu/go-restapi-authentication-jwt/db"
	"github.com/JuanLopezAranzazu/go-restapi-authentication-jwt/middlewares"
	"github.com/JuanLopezAranzazu/go-restapi-authentication-jwt/models"
	"github.com/JuanLopezAranzazu/go-restapi-authentication-jwt/routes"
	"github.com/gorilla/mux"
)

func main() {
	// conexion con la base de datos
	db.DBConnection()
	// migraciones de las tablas
	if err := db.DB.AutoMigrate(models.User{}); err != nil {
		log.Fatal("Error en migración de tablas: ", err)
	}
	// manejar rutas
	r := mux.NewRouter()
	// index
	r.HandleFunc("/", routes.HomeHandler)

	s := r.PathPrefix("/api/v1").Subrouter()

	// autenticacion de usuarios
	s.HandleFunc("/auth/login", routes.LoginHandler).Methods("POST")
	s.HandleFunc("/auth/register", routes.RegisterHandler).Methods("POST")
	s.HandleFunc("/auth/refresh", routes.RefreshHandler).Methods("POST")

	s.Handle("/auth/me", middlewares.JWTMiddleware(http.HandlerFunc(routes.MeHandler))).Methods("GET")

	// iniciar servidor
	log.Println("Servidor iniciado en http://localhost:3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
	}
}
