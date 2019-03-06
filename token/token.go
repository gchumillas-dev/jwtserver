package token

import jwt "github.com/dgrijalva/jwt-go"

func New(privateKey string, claims jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(privateKey))
	if err != nil {
		panic(err)
	}

	return signedToken
}

func Parse(privateKey, token string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(privateKey), nil
	})
}
