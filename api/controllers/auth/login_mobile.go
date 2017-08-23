package c_auth

import (
	"chess/api/components/input"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	//"chess/api/components/sms"
	"chess/api/components/user_init"
	grpcServer "chess/api/grpc"
	pb "chess/api/proto"
	"chess/common/config"
	"chess/common/helper"
	"chess/common/log"
	"chess/models"
	"golang.org/x/net/context"
)

var (
	VerifyCodeExpired = errors.New("Verify code is expired.")
)

type LoginMobileParams struct {
	MobileNumber string `json:"mobile_number" form:"mobile_num" binding:"required" description:"手机号"`
	Password     string `json:"password" binding:"required" description:"密码"`
	From         string `json:"from" description:"来源"`
	UniqueId     string `json:"unique_id"  description:"唯一标识"`
	Channel      string `json:"channel" description:"渠道"`
	//Q            string                 `json:"q"`
	//BindingParam broker.CheckCodeParams `json:"binding_param"`
}

// @Title 手机号注册
// @Description 手机注册
// @Summary 手机号注册
// @Accept json
// @Param   body     body    c_auth.LoginMobileParams  true        "post 数据"
// @Success 200 {object} c_auth.LoginResult
// @router /auth/login/quick [post]
func LoginMobile(c *gin.Context) {
	var result LoginResult
	var post LoginMobileParams

	_conf, ok1 := c.Get("config")
	cConf, ok2 := _conf.(*config.ApiConfig)
	if !ok1 || !ok2 {
		result.Msg = "Get config fail."
		c.JSON(http.StatusOK, result)
		return
	}

	//var err error
	var user = new(models.UsersModel)
	if input.BindJSON(c, &post, cConf) == nil {
		clientIp := helper.ClientIP(c)
		if post.Channel == "" {
			post.Channel = "default"
		}
		if post.From == "" {
			post.From = "default"
		}
		post.From = strings.ToLower(post.From)
		//result.Ret, result.Msg, err = sms.CheckCode(post.MobileNumber, post.Password, sms.SMS_LOGIN, cConf)
		//if err != nil {
		//	// 验证不通过
		//	c.JSON(http.StatusOK, result)
		//	return
		//}

		// Query the user is exists
		var userId int
		err := models.Users.GetByMobileNumber(post.MobileNumber, user)

		if err != nil {
			// Create a new user
			user.MobileNumber = post.MobileNumber
			user.RegIp = clientIp
			user.LastLoginIp = user.RegIp
			user.ContactMobile = post.MobileNumber
			user.Avatar = cConf.User.DefaultAvatar
			user.Channel = post.Channel
			user.Type = models.TYPE_MOBILE_QUICK_LOGIN
			user.AppFrom = post.From
			user.Nickname = helper.GenMobileNickname(post.MobileNumber)
			userId, err = models.Users.Insert(user)
			if err != nil {
				log.Error(err)
				result.Ret = 0
				result.Msg = "Could not create new user."
				c.JSON(http.StatusOK, result)
				return
			}
			user.Id = userId

			// Init user wallet
			models.UsersWallet.Init(userId)

			go func() {
				//更新设备信息
				err = user_init.DeviceInit(user.Id, user.AppFrom, post.UniqueId, c.Query("idfv"), c.Query("idfa"))
				if err != nil {
					log.Error(err)
				}
				//初始化用户任务系统
				err =  user_init.TaskInit(user.Id)
				if err != nil {
				    log.Error(err)
				}

			}()

		} else {

			userId = user.Id

			// update last login ip
			go models.Users.UpdateLastLoginIp(userId, clientIp)
		}

		extra := make(map[string]interface{})
		extra["from"] = user.AppFrom
		extra["unique_id"] = post.UniqueId
		extra["idfv"] = c.Query("idfv")
		extra["idfa"] = c.Query("idfa")
		//user.IsFresh, err = user_init.UserInit(*user, extra, cConf)
		if err != nil {
			log.Debugf("%+v", extra)
			log.Error(err)
		}

		// Create login token
		AuthClient := grpcServer.GetAuthGrpc()
		authResult, err := AuthClient.RefreshToken(context.Background(), &pb.RefreshTokenArgs{UserId: int32(user.Id), AppFrom: user.AppFrom, UniqueId: post.UniqueId})
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
