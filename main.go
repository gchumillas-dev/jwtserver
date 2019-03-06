package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gchumillas/ucms/handler"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	apiVersion := os.Getenv("apiVersion")
	serverAddr := os.Getenv("serverAddr")
	dbName := os.Getenv("dbName")
	dbUser := os.Getenv("dbUser")
	dbPass := os.Getenv("dbPass")

	dsName := fmt.Sprintf("%s:%s@/%s", dbUser, dbPass, dbName)
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	env := &handler.Env{DB: db}
	prefix := fmt.Sprintf("/%s", strings.TrimLeft(apiVersion, "/"))
	r := mux.NewRouter()
	public := r.PathPrefix(prefix).Subrouter()
	private := r.PathPrefix(prefix).Subrouter()
	private.Use(env.AuthMiddleware)

	// authentication
	public.HandleFunc("/sign/in", env.SignIn).Methods("POST")
	public.HandleFunc("/sign/up", env.SignUp).Methods("POST")
	private.HandleFunc("/sign/out", env.SignOut).Methods("POST")

	log.Printf("Server started at port %s", serverAddr)
	log.Fatal(http.ListenAndServe(
		serverAddr,
		handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(false),
		)(r),
	))
}
