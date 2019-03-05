package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pelletier/go-toml"

	_ "github.com/go-sql-driver/mysql"
)

type dbConfig struct {
	DB   string `toml:"dbname"`
	User string `toml:"dbuser"`
	Pass string `toml:"dbpass"`
}

type config struct {
	APIVersion string `toml:"apiVersion"`
	Database   dbConfig
}

func main() {
	conf := loadConfig("config.toml")

	dbConf := conf.Database
	dsName := fmt.Sprintf("%s:%s@/%s", dbConf.User, dbConf.Pass, dbConf.DB)
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	prefix := fmt.Sprintf("/%s", strings.TrimLeft(conf.APIVersion, "/"))
	r := mux.NewRouter()
	s := r.PathPrefix(prefix).Subrouter()

	// authentication (public routes)
	auth := s.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
	}).Methods("GET")

	serverAddr := "localhost:8080"
	log.Printf("Server started at port %v", serverAddr)
	log.Fatal(http.ListenAndServe(
		serverAddr,
		handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(false),
		)(r),
	))
}

func loadConfig(filename string) (conf config) {
	conf = config{}

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	decoder := toml.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		panic(err)
	}

	return
}
