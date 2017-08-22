package c_auth

import (
	"chess/api/components/input"
	"chess/api/components/user_init"
	grpcServer "chess/api/grpc"
	pb "chess/api/proto"
	"chess/common/config"
	"chess/common/helper"
	"chess/common/log"
	"chess/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"strings"
)

type TouristParams struct {
	From     string `json:"from" form:"from" description:"请求来源"`
	UniqueId string `json:"unique_id" form:"unique_id" description:"唯一标识"`
	Channel  string `json:"channel" description:"渠道"`
}

// @Title 游客登录
// @Description 游客登录
// @Summary 游客登录
// @Accept json
// @Param   body     body    c_auth.TouristParams  true        "post 数据"
// @Success 200 {object} c_auth.LoginResult
// @router /auth/login/tourist [post]
func TouristLogin(c *gin.Context) {
	var parms TouristParams
	var result LoginResult
	_conf, ok1 := c.Get("config")
	cConf, ok2 := _conf.(*config.ApiConfig)
	if !ok1 || !ok2 {
		result.Msg = "Get config fail."
		c.JSON(http.StatusOK, result)
		return
	}

	var user = new(models.UsersModel)
	if input.BindJSON(c, &parms, cConf) == nil {
		clientIp := helper.ClientIP(c)
		if parms.Channel == "" {
			parms.Channel = "default"
		}
		if parms.From == "" {
			parms.From = "default"
		}
		parms.From = strings.ToLower(parms.From)
		user.RegIp = clientIp
		user.LastLoginIp = clientIp
		user.Channel = parms.Channel
		user.Type = models.TYPE_TOURIST
		user.AppFrom = parms.From
		userId, err := models.Users.Insert(user)
		if err != nil {
			log.Error(err)
			result.Ret = 0
			result.Msg = "Could not create new user."
			c.JSON(http.StatusOK, result)
			return
		}
		err = models.Users.UpdateNickname(userId, strconv.Itoa(userId))
		if err != nil {

		}
		models.UsersWallet.Init(userId)
		go func() {
			//更新设备信息
			err = user_init.DeviceInit(user.Id, user.AppFrom, parms.UniqueId, c.Query("idfv"), c.Query("idfa"))
			if err != nil {
				log.Error(err)
			}

		}()

		// Create login token
		AuthClient := grpcServer.GetAuthGrpc()
		authResult, err := AuthClient.RefreshToken(context.Background(), &pb.RefreshTokenArgs{UserId: int32(user.Id), AppFrom: user.AppFrom, UniqueId: parms.UniqueId})
		if err != nil {
			result.Msg = "login failed"
			c.JSON(http.StatusOK, result)
			return
		}

		result.Ret = 1
		result.UserId = userId
		result.IsFresh = user.IsFresh
		result.Token = authResult.Token
		result.Expire = authResult.Expire
		result.RefreshToken = authResult.RefreshToken
		c.JSON(http.StatusOK, result)
		return

	}
	result.Msg = "Params invaild."
	c.JSON(http.StatusOK, result)
}
