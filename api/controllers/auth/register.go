package c_auth

import (
	"chess/api/components/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	//"chess/api/components/broker"
	"chess/api/components/input"
	//"chess/api/components/sms"
	"chess/api/components/user_init"
	grpcServer "chess/api/grpc"
	"chess/common/helper"
	"chess/common/log"
	pb "chess/api/proto"
	"chess/common/config"
	"chess/models"
	"golang.org/x/net/context"
)

const (
	USER_HAS_REG     = 3
	USER_QUICK_LOGIN = 2
)

type RegisterMobileParams struct {
	// @todo email...
	MobileNumber string `json:"mobile_number" binding:"required" description:"手机号"`
	Password     string `json:"password" description:"密码"`
	Gender       int    `json:"gender" description:"性别0未知 1男 2女(可以不传，默认为0)"` // 性别，可以不传，默认为0
	Code         string `json:"code" binding:"required" description:"验证码"`
	From         string `json:"from" description:"来源"`
	UniqueId     string `json:"unique_id"description:"唯一标识0" `
	Channel      string `json:"channel" description:"渠道"`
	// BindingParam broker.CheckCodeParams `json:"binding_param"`
}

// @Title 手机号注册
// @Description 手机注册
// @Summary 手机号注册
// @Accept json
// @Param   body     body    c_auth.RegisterMobileParams  true        "post 数据"
// @Success 200 {object} c_auth.LoginResult
// @router /auth/register/mobile [post]
func RegisterMobile(c *gin.Context) {
	var result LoginResult
	var post RegisterMobileParams
	var user = new(models.UsersModel)

	_conf, ok1 := c.Get("config")
	cConf, ok2 := _conf.(*config.ApiConfig)
	if !ok1 || !ok2 {
		result.Msg = "Get config fail."
		c.JSON(http.StatusOK, result)
		return
	}

	if input.BindJSON(c, &post, cConf) == nil {
		clientIp := helper.ClientIP(c)
		if post.Channel == "" {
			post.Channel = "default"
		}
		if post.From == "" {
			post.From = "default"
		}
		post.From = strings.ToLower(post.From)
		// 判断手机号码是否被注册了
		// Query the user is exists
		err := models.Users.GetByMobileNumber(post.MobileNumber, user)
		if err == nil {
			// 注册过了
			if user.Pwd == "" {
				result.Ret = USER_QUICK_LOGIN
				result.Msg = "plz login wish sms code"
				c.JSON(http.StatusOK, result)
				return
			}
			result.Ret = USER_HAS_REG
			result.Msg = "this mobile is exit"
			c.JSON(http.StatusOK, result)
			return
		}
		//result.Ret, result.Msg, err = sms.CheckCode(post.MobileNumber, post.Code, sms.SMS_REGISTER, cConf)
		//if err != nil {
		//    // 验证不通过
		//    c.JSON(http.StatusOK, result)
		//    return
		//}

		pwdRet := auth.Passwords.CheckPasswordStrong(post.Password)
		if pwdRet != 1 {
			result.Ret = pwdRet
			result.Msg = "password not strong engough"
			c.JSON(http.StatusOK, result)
			return
		}
		hp, err := auth.Passwords.Hash(post.Password)
		if err != nil {
			result.Ret = 0
			result.Msg = "server error"
			c.JSON(http.StatusOK, result)
			return
		}
		var userId int
		// Create a new user
		// TODO: add the user channel
		user.MobileNumber = post.MobileNumber
		user.Pwd = hp

		user.RegIp = clientIp
		user.LastLoginIp = user.RegIp
		user.Channel = post.Channel
		user.AppFrom = post.From
		user.ContactMobile = post.MobileNumber
		user.Avatar = cConf.User.DefaultAvatar
		user.Gender = post.Gender
		user.Type = models.TYPE_MOBILE_REG
		user.Nickname = helper.GenMobileNickname(post.MobileNumber)
		userId, err = models.Users.Insert(user)
		if err != nil {
			log.Error(err)
			result.Msg = "Could not create new user."
			c.JSON(http.StatusOK, result)
			return
		}
		user.Id = userId

		// Init user wallet
		models.UsersWallet.Init(userId)

		extra := make(map[string]interface{})
		extra["from"] = user.AppFrom
		extra["unique_id"] = post.UniqueId
		extra["idfv"] = c.Query("idfv")
		extra["idfa"] = c.Query("idfa")
		//user.IsFresh, err = user_init.UserInit(*user, extra, cConf)
		if err != nil {
			log.Error(err)
		}
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

		// Create login token
		// Create login token
		AuthClient := grpcServer.GetAuthGrpc()
		authResult, err := AuthClient.RefreshToken(context.Background(), &pb.RefreshTokenArgs{UserId: int32(user.Id), AppFrom: user.AppFrom, UniqueId: post.UniqueId})
		//authResult, err := auth.LoginUser(userId, post.From, post.UniqueId)
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
