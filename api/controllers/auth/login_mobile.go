package c_auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"chess/api/components/auth"
	"chess/api/components/input"
	"chess/api/components/sms"
	"chess/api/components/user_init"
	"chess/common/config"
	"chess/api/helper"
	"chess/api/log"
	"chess/models"
    "fmt"
)

var (
	VerifyCodeExpired = errors.New("Verify code is expired.")
)

type LoginMobileParams struct {
	MobileNumber string                 `json:"mobile_number" form:"mobile_num" binding:"required"`
	Password     string                 `json:"password" binding:"required" `
	From         string                 `json:"from" form:"from"`
	UniqueId     string                 `json:"unique_id" form:"unique_id"`
	Channel      string                 `json:"channel"`
	Q            string                 `json:"q"`
	//BindingParam broker.CheckCodeParams `json:"binding_param"`
}

type LoginMobileResult struct {
	LoginResult
}

func LoginMobile(c *gin.Context) {
	var result LoginMobileResult
	var post LoginMobileParams

	_conf, ok1 := c.Get("config")
	cConf, ok2 := _conf.(*config.ApiConfig)
	if !ok1 || !ok2 {
		result.Msg = "Get config fail."
		c.JSON(http.StatusOK, result)
		return
	}

	var err error
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
		result.Ret, result.Msg, err = sms.CheckCode(post.MobileNumber, post.Password, sms.SMS_LOGIN, cConf)
		if err != nil {
			// 验证不通过
			c.JSON(http.StatusOK, result)
			return
		}


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
				log.Log.Error(err)
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
					log.Log.Error(err)
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
			log.Log.Debugf("%+v", extra)
			log.Log.Error(err)
		}

		// Create login token
		authResult, err := auth.LoginUser(userId, post.From, post.UniqueId)
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
