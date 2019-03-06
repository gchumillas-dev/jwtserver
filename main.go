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
	ServerAddr string `toml:"serverAddr"`
	PrivateKey string `toml:"privateKey"`
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

	env := &handler.Env{DB: db}
	prefix := fmt.Sprintf("/%s", strings.TrimLeft(conf.APIVersion, "/"))
	r := mux.NewRouter()
	public := r.PathPrefix(prefix).Subrouter()
	private := r.PathPrefix(prefix).Subrouter()
	private.Use(func(next http.Handler) http.Handler {
		return env.AuthMiddleware(next, conf.PrivateKey)
	})

	// authentication
	public.HandleFunc("/sign/in", func(w http.ResponseWriter, r *http.Request) {
		env.SignIn(w, r, conf.PrivateKey)
	}).Methods("POST")
	public.HandleFunc("/sign/up", env.SignUp)
	private.HandleFunc("/sign/out", env.SignOut).Methods("POST")

	log.Printf("Server started at port %s", conf.ServerAddr)
	log.Fatal(http.ListenAndServe(
		conf.ServerAddr,
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
