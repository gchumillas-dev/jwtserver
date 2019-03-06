package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/gchumillas/ucms/manager"
)

// AuthMiddleware verifies that the user was successful authorized.
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

		u := manager.NewUser()
		if !u.ReadUserByToken(env.DB, env.PrivateKey, token) {
			httpError(w, unauthorizedError)
			return
		}

		ctx := context.WithValue(r.Context(), contextUserKey, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
