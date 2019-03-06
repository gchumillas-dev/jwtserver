package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gchumillas/ucms/manager"
)

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

	json.NewEncoder(w).Encode(u.NewToken(env.PrivateKey))
}
