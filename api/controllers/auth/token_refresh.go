package c_auth

import (
	"chess/api/components/input"
	grpcServer "chess/api/grpc"
	pb "chess/api/proto"
	"chess/common/auth"
	"chess/common/config"
	"chess/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net/http"
)

type TokenRefreshParams struct {
	UserId       int    `json:"user_id" binding:"required" description:"用户id"`
	From         string `json:"from" description:"请求来源"`
	UniqueId     string `json:"unique_id" form:"unique_id" description:"唯一设备"`
	OldToken     string `json:"token" binding:"required" description:"旧的token"`
	RefreshToken string `json:"refresh_token" binding:"required" description:"旧的刷新码"`
}

type TokenRefreshResult struct {
	LoginResult
}

// @Title 刷新token
// @Description 刷新token
// @Summary 刷新token
// @Accept json
// @Param   body     body    c_auth.TokenRefreshParams  true        "post 数据"
// @Success 200 {object} c_auth.TokenRefreshResult
// @router /auth/token/refresh [post]
func TokenRefrash(c *gin.Context) {
	var result TokenRefreshResult
	var post TokenRefreshParams

	_conf, ok1 := c.Get("config")
	cConf, ok2 := _conf.(*config.ApiConfig)
	if !ok1 || !ok2 {
		result.Msg = "Get config fail."
		c.JSON(http.StatusOK, result)
		return
	}

	if input.BindJSON(c, &post, cConf) == nil {
		// Get the session
		session, err := models.Session.Get(post.UserId)
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
		AuthClient, ret := grpcServer.GetAuthGrpc()
		if ret == 0 {
			result.Msg = "rpc fail"
			c.JSON(http.StatusOK, result)
			return
		}
		authResult, err := AuthClient.RefreshToken(context.Background(), &pb.RefreshTokenArgs{UserId: int32(post.UserId), AppFrom: post.From, UniqueId: post.UniqueId})

		result.Ret = 1
		result.UserId = post.UserId
		result.Token = authResult.Token
		result.Expire = authResult.Expire
		result.RefreshToken = authResult.RefreshToken
		c.JSON(http.StatusOK, result)
		return
	}

	result.Msg = "Params invaild."
	c.JSON(http.StatusOK, result)
}
