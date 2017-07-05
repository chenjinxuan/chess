package middleware

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func SetContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		appFrom := strings.ToLower(c.Query("from"))
		if appFrom == "" {
			appFrom = "ios"
		}
		c.Set("app_from", appFrom)

		appVer, _ := strconv.Atoi(c.Query("ver"))
		c.Set("app_ver", appVer)

		appChannel := c.Query("channel")
		c.Set("app_channel", appChannel)
	}
}

func GetString(c *gin.Context, key string) string {
	_val, ok1 := c.Get(key)
	val, ok2 := _val.(string)
	if ok1 && ok2 {
		return val
	}
	return ""
}

func GetInt(c *gin.Context, key string) int {
	_val, ok1 := c.Get(key)
	val, ok2 := _val.(int)
	if ok1 && ok2 {
		return val
	}
	return 0
}
