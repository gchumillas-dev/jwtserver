package handler

import (
	"net/http"
)

func (env *Env) Home(w http.ResponseWriter, r *http.Request) {
	u := getUser(r)
}
