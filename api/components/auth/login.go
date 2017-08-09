package auth

import (
	"errors"
	"github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"chess/common/define"
	"chess/common/log"
)

const (
	AuthFailedStatus = -998
)

var (
	UnexpectedSigningMethod = errors.New("Unexpected signing method.")
	AuthFailed              = errors.New("Auth failed.")
)


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

func Login(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var result define.BaseResult
		userId := c.Param("user_id")
		tokenString := c.Query("token")

		loginData, err := AuthLoginToken(tokenString, secret)

		log.Debugf("Auth",logrus.Fields{
			"User ID":   userId,
			"Token":     tokenString,
			"LoginData": loginData,
			"Error":     err,
		})

		if err != nil || userId != loginData {
			result.Ret = AuthFailedStatus
			result.Msg = AuthFailed.Error()
			c.JSON(http.StatusOK, result)
			c.Abort()
			return
		}
	}
}
