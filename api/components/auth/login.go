package auth

import (
	"chess/api/redis"
	"chess/common/define"
	"chess/common/log"
	"errors"
	"github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
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
		//判断黑名单
		msg, err := redis.Redis.Login.GetInt(tokenString)
		if err == nil {
			result.Ret = define.AuthALreadyLogin
			result.Msg = define.AuthMsgMap[msg]
			c.JSON(http.StatusOK, result)
			c.Abort()
			return
		}
		loginData, err := AuthLoginToken(tokenString, secret)

		log.Debugf("Auth", logrus.Fields{
			"User ID":   userId,
			"Token":     tokenString,
			"LoginData": loginData,
			"Error":     err,
		})

		if err != nil || userId != loginData {
			result.Ret = define.AuthFailedStatus
			result.Msg = AuthFailed.Error()
			c.JSON(http.StatusOK, result)
			c.Abort()
			return
		}
	}
}
