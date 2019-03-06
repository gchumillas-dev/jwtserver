package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gchumillas/jwtserver/manager"
)

// SignIn handler.
func (env *Env) SignIn(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string
		Password string
	}
	parseBody(w, r, &body)

	u := manager.NewUser()
	if !u.ReadUserByCredentials(env.DB, body.Username, body.Password) {
		httpError(w, docNotFoundError)
		return
	}

	json.NewEncoder(w).Encode(u.NewToken(env.PrivateKey, env.Expiration))
}

// Home handler.
func (env *Env) Home(w http.ResponseWriter, r *http.Request) {
	u := getUser(r)

	msg := fmt.Sprintf("Hello %s!\nYou are ahthorized.", u.Username)
	io.WriteString(w, msg)
}
