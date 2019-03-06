package manager

import (
	"database/sql"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gchumillas/ucms/token"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Username string
}

type userClaims struct {
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

func (user *User) NewToken(privateKey string) string {
	claims := userClaims{UserID: user.ID}

	return token.New(privateKey, claims)
}

// CreateUser creates a user.
func (user *User) CreateUser(db *sql.DB, password string) {
	stmt, err := db.Prepare(`
		insert into user(username, password)
		values(?, ?)`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if _, err := stmt.Exec(user.Username, pwd); err != nil {
		panic(err)
	}
}

func (user *User) ReadUser(db *sql.DB, ID string) (found bool) {
	stmt, err := db.Prepare(`
		select id, username
		from user where id = ?`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	switch err := stmt.QueryRow(ID).Scan(&user.ID, &user.Username); {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		panic(err)
	}

	return true
}

func (user *User) ReadUserByCredentials(db *sql.DB, uname string, upass string) (found bool) {
	stmt, err := db.Prepare(`
		select id, username, password
		from user where username = ?`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	var hashedPassword []byte
	switch err := stmt.QueryRow(uname).Scan(&user.ID, &user.Username, &hashedPassword); {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		panic(err)
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(upass)); err != nil {
		return false
	}

	return true
}

func (user *User) ReadUserByToken(db *sql.DB, privateKey, signedToken string) (found bool) {
	claims := &userClaims{UserID: user.ID}
	_, err := token.Parse(privateKey, signedToken, claims)
	if err != nil {
		panic(err)
	}

	return user.ReadUser(db, claims.UserID)
}
