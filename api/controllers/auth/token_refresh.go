package c_auth

import (
	// "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"net/http"
	"strconv"
	"time"
	"treasure/components/auth"
	"treasure/components/input"
	"treasure/config"
	"treasure/define"
	// "treasure/log"
	"strings"
	"treasure/models"
)

type TokenRefreshParams struct {
	UserId       int    `json:"user_id" binding:"required"`
	From         string `json:"from"`
	UniqueId     string `json:"unique_id" form:"unique_id"`
	OldToken     string `json:"token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokenRefreshResult struct {
	LoginResult
}

func TokenRefrash(c *gin.Context) {
	var result TokenRefreshResult
	var post TokenRefreshParams

	_conf, ok1 := c.Get("config")
	cConf, ok2 := _conf.(*config.Config)
	if !ok1 || !ok2 {
		result.Msg = "Get config fail."
		c.JSON(http.StatusOK, result)
		return
	}

	if input.BindJSON(c, &post, cConf) == nil {
		// Get the session
		session, err := models.Session.Get(post.UserId, strings.ToLower(post.From), post.UniqueId)
		if err != nil {
			result.Msg = auth.AuthFailed.Error()
			c.JSON(http.StatusOK, result)
			return
		}

		// Verify old token
		if session.Token.Data != post.OldToken {
			result.Msg = auth.AuthFailed.Error()
			c.JSON(http.StatusOK, result)
			return
		}

		// Verify refresh token
		if session.RefreshToken != post.RefreshToken {
			result.Msg = auth.AuthFailed.Error()
			c.JSON(http.StatusOK, result)
			return
		}

		// Generate a new login token
		expire := time.Now().Add(time.Second * time.Duration(config.C.Login.TokenExpire)).Unix()
		tokenString, err := auth.CreateLoginToken(strconv.Itoa(post.UserId), expire, define.JwtSecret)
		if err != nil {
			result.Msg = "Could not generate token."
			c.JSON(http.StatusOK, result)
			return
		}

		// Generate a new refresh token
		u := uuid.NewV4()
		refreshToken := u.String()

		// Update session and update database
		sessionUpdated := new(models.SessionModel)
		sessionUpdated.UserId = session.UserId
		sessionUpdated.From = session.From
		sessionUpdated.UniqueId = session.UniqueId
		sessionUpdated.Token = &models.SessionToken{tokenString, expire}
		sessionUpdated.RefreshToken = refreshToken
		sessionUpdated.Updated = time.Now()
		sessionUpdated.Created = session.Created
		err = models.Session.Upsert(post.UserId, post.From, post.UniqueId, sessionUpdated)
		if err != nil {
			result.Msg = "Could not generate session."
			c.JSON(http.StatusOK, result)
			return
		}

		result.Ret = 1
		result.UserId = post.UserId
		result.Token = tokenString
		result.Expire = expire
		result.RefreshToken = refreshToken
		c.JSON(http.StatusOK, result)
		return
	}

	result.Msg = "Params invaild."
	c.JSON(http.StatusOK, result)
}
