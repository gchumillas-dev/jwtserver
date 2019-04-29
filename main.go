package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gchumillas/jwtserver/handler"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	apiVersion := os.Getenv("apiVersion")
	serverAddr := os.Getenv("serverAddr")
	privateKey := os.Getenv("privateKey")
	expiration, _ := time.ParseDuration(os.Getenv("expiration"))
	dbName := os.Getenv("dbName")
	dbUser := os.Getenv("dbUser")
	dbPass := os.Getenv("dbPass")

	dsName := fmt.Sprintf("%s:%s@/%s", dbUser, dbPass, dbName)
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	env := &handler.Env{
		DB:         db,
		PrivateKey: privateKey,
		Expiration: expiration}
	prefix := fmt.Sprintf("/%s", strings.TrimLeft(apiVersion, "/"))
	r := mux.NewRouter()

	// public routes
	public := r.PathPrefix(prefix).Subrouter()
	public.HandleFunc("/sign/in", env.SignIn).Methods("POST")

	// private routes
	private := r.PathPrefix(prefix).Subrouter()
	private.HandleFunc("/home", env.Home).Methods("GET")
	private.Use(env.AuthMiddleware)

	// CORS support
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})

	log.Printf("Server started at port %s", serverAddr)
	log.Fatal(http.ListenAndServe(
		serverAddr,
		handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(false),
		)(handlers.CORS(headersOk, originsOk, methodsOk)(r)),
	))
}
