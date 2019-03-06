package handler

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gchumillas/ucms/manager"
)

func (env *Env) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := ""
		items := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(items) > 1 {
			token = items[1]
		}

		if len(token) == 0 {
			httpError(w, unauthorizedError)
			return
		}

		privateKey := os.Getenv("privateKey")
		u := manager.NewUser()
		u.ReadUserByToken(env.DB, privateKey, token)
		log.Println(token)

		next.ServeHTTP(w, r)
	})
}
