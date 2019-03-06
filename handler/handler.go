package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/gchumillas/ucms/manager"
)

// Env contains common variables, such as the database access, etc.
type Env struct {
	DB         *sql.DB
	PrivateKey string
}

// Context keys.
type contextKey string

const contextUserKey = contextKey("context-user-key")

// Common HTTP status errors.
type httpStatus struct {
	code int
	msg  string
}

var (
	docNotFoundError  = httpStatus{404, "Document not found."}
	unauthorizedError = httpStatus{401, "Not authorized."}
)

func getUser(r *http.Request) *manager.User {
	return r.Context().Value(contextUserKey).(*manager.User)
}

func httpError(w http.ResponseWriter, status httpStatus) {
	http.Error(w, status.msg, status.code)
	log.Printf("http error (%d): %s", status.code, status.msg)
	return
}

func parseBody(w http.ResponseWriter, r *http.Request, body interface{}) {
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(body); err != nil {
		panic(err)
	}

	// Removes whitespaces around texts.
	elem := reflect.ValueOf(body).Elem()
	switch reflect.TypeOf(elem).Kind() {
	case reflect.Struct:
		count := elem.NumField()
		for i := 0; i < count; i++ {
			field := elem.Field(i)
			if field.Type().Kind() == reflect.String {
				field.SetString(strings.Trim(field.String(), " "))
			}
		}
	}
}
