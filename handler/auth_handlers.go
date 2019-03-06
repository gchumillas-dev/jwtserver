package handler

import (
	"encoding/json"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gchumillas/ucms/manager"
)

func (env *Env) SignIn(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string
		Password string
	}
	parseBody(w, r, &body)

	u := manager.NewUser()
	if err := u.ReadUserByCredentials(env.DB, body.Username, []byte(body.Password)); err != nil {
		httpError(w, docNotFoundError)
		return
	}

	// NewToken(privateKey)
	// u.GenerateToken()
	privateKey := os.Getenv("privateKey")
	claim := &manager.UserClaim{UserID: u.ID}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString([]byte(privateKey))
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(signedToken)
}

func (env *Env) SignOut(w http.ResponseWriter, r *http.Request) {

}

func (env *Env) SignUp(w http.ResponseWriter, r *http.Request) {

}
