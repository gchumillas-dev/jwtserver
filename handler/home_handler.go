package handler

import (
	"fmt"
	"io"
	"net/http"
)

func (env *Env) Home(w http.ResponseWriter, r *http.Request) {
	u := getUser(r)

	msg := fmt.Sprintf("Hello %s!\nYou are ahthorized.", u.Username)
	io.WriteString(w, msg)
}
