package handler

import (
	"log"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gchumillas/ucms/manager"
)

func (env *Env) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signedToken := ""
		items := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(items) > 1 {
			signedToken = items[1]
		}

		if len(signedToken) == 0 {
			httpError(w, unauthorizedError)
			return
		}

		// u.ReadUserByToken(token)
		claim := manager.UserClaim{}
		jwt.ParseWithClaims(signedToken, &claim, func(t *jwt.Token) (interface{}, error) {
			privateKey := os.Getenv("privateKey")
			return []byte(privateKey), nil
		})
		log.Println(claim.UserID)

		// u := &manager.User{AccessToken: token}
		// if !u.ReadUserByToken(env.DB) {
		// 	httpError(w, unauthorizedError)
		// 	return
		// }
		//
		// ctx := context.WithValue(r.Context(), contextAuthUser, u)
		// next.ServeHTTP(w, r.WithContext(ctx))

		next.ServeHTTP(w, r)
	})
}
