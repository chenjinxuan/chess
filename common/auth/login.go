package auth

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
)

const (
	AuthFailedStatus = -998
)

var (
	UnexpectedSigningMethod = errors.New("Unexpected signing method.")
	AuthFailed              = errors.New("Auth failed.")
)

func CreateLoginToken(loginData string, expire int64, secret string) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := make(jwt.MapClaims)
	claims["userid"] = loginData
	claims["exp"] = expire
	token.Claims = claims
	return token.SignedString([]byte(secret))
}

func AuthLoginToken(tokenString string, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		method := jwt.GetSigningMethod("HS256")
		if token.Method.Alg() != method.Alg() {
			return nil, UnexpectedSigningMethod
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims["userid"].(string), nil
}
