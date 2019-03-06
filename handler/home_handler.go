package handler

import (
	"log"
	"net/http"
)

func (env *Env) Home(w http.ResponseWriter, r *http.Request) {
	u := getUser(r)
	log.Println(u.ID)
}
