package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gchumillas/ucms/manager"
	"github.com/gchumillas/ucms/token"
)

func (env *Env) SignIn(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string
		Password string
	}
	parseBody(w, r, &body)

	u := manager.NewUser()
	if err := u.ReadUserByCredentials(env.DB, body.Username, body.Password); err != nil {
		httpError(w, docNotFoundError)
		return
	}

	privateKey := os.Getenv("privateKey")
	claims := manager.UserClaims{UserID: u.ID}
	signedToken := token.New(privateKey, claims)

	json.NewEncoder(w).Encode(signedToken)
}

func (env *Env) SignOut(w http.ResponseWriter, r *http.Request) {

}

func (env *Env) SignUp(w http.ResponseWriter, r *http.Request) {

}
