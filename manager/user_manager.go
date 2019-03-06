package manager

import (
	"database/sql"

	jwt "github.com/dgrijalva/jwt-go"
	tk "github.com/gchumillas/ucms/token"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Username string
	Password []byte
}

type UserClaims struct {
	UserID string
	jwt.StandardClaims
}

// NewUser creates a user.
func NewUser(userID ...string) *User {
	id := ""
	if len(userID) > 0 {
		id = userID[0]
	}

	return &User{ID: id}
}

// CreateUser creates a user.
func (user *User) CreateUser(db *sql.DB) {
	stmt, err := db.Prepare(`
		insert into user(username, password)
		values(?, ?)`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	pwd, err := bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)
	if _, err := stmt.Exec(user.Username, pwd); err != nil {
		panic(err)
	}
}

func (user *User) ReadUserByCredentials(db *sql.DB, uname string, upass string) error {
	stmt, err := db.Prepare(`
		select id, username, password
		from user where username = ?`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(uname).
		Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(upass)); err != nil {
		return err
	}

	return nil
}

func (user *User) ReadUserByToken(db *sql.DB, privateKey, token string) {
	claims := &UserClaims{UserID: user.ID}
	tk.Parse(privateKey, token, claims)
}
