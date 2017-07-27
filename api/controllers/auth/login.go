package c_auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"chess/api/components/auth"
	"chess/api/components/captcha"
	"chess/api/components/input"
	"chess/common/config"
	"chess/api/define"
	"chess/api/helper"
	"chess/api/log"
	"chess/models"
    "fmt"
)

const (
	No_Captcha       = -3
	Wrong_Captcha    = -4
	Login_Fail_Limit = -5

	Login_Wich_Code = 2
	Password_Wrong  = -12
)

type LoginParams struct {
	LoginBy      string `json:"loginby" form:"loginby"`
	MobileNumber string `json:"mobile_number" form:"mobile_number" binding:"required"`
	Password     string `json:"password" form:"password" binding:"required"`
	From         string `json:"from" form:"from"`
	UniqueId     string `json:"unique_id" form:"unique_id"`
	Captcha      string `json:"captcha"`
}

type LoginResult struct {
	define.BaseResult
	NeedCaptcha  bool   `json:"need_captcha"`
	UserId       int    `json:"user_id"`
	IsFresh      int    `json:"fresh"`
	Token        string `json:"token"`
	Expire       int64  `json:"expire"`
	RefreshToken string `json:"refresh_token"`
}

// 手机号 + 密码登录
func Login(c *gin.Context) {
	var result LoginResult
	var form LoginParams

	_conf, ok1 := c.Get("config")
	cConf, ok2 := _conf.(*config.ApiConfig)
	if !ok1 || !ok2 {
		result.Msg = "Get config fail."
		c.JSON(http.StatusOK, result)
		return
	}
	//var err error
	var user = new(models.UsersModel)
	if input.BindJSON(c, &form, cConf) == nil {
		clientIp := helper.ClientIP(c)
		form.From = strings.ToLower(form.From)
		// check fail time
		failCount, err := auth.GetFailLoginCount(form.MobileNumber)
		if err != nil {
			log.Log.Error(err)
			result.Ret = 0
			result.Msg = "server error"
			c.JSON(http.StatusOK, result)
			return
		}
		if failCount >= cConf.Login.FailLimit {
			// 达到错误上线,不允许登录
			result.Ret = Login_Fail_Limit
			result.Msg = "login limit"
			c.JSON(http.StatusOK, result)
			return
		}
		if failCount >=  cConf.Login.ShowCaptchaCount {
			result.NeedCaptcha = true
			if form.Captcha == "" {
				// 未填写验证码
				result.Ret = No_Captcha
				result.Msg = "require captcha"
				c.JSON(http.StatusOK, result)
				return
			}
			// check captcha
			isOK := captcha.VerifyByMobileLogin(form.MobileNumber, form.Captcha)
			if !isOK {
				// 验证码错误
				result.Ret = Wrong_Captcha
				result.Msg = "wrong captcha"
				c.JSON(http.StatusOK, result)
				return
			}
		}

		err = models.Users.GetByMobileNumber(form.MobileNumber, user)
		if err != nil {
			//result.Ret = 0
			// result.Msg = "server error"
			result.Ret = Password_Wrong
			result.Msg = "password wrong"
			c.JSON(http.StatusOK, result)
			return
		}

		if user.Pwd == "" {
			result.Ret = Login_Wich_Code
			result.Msg = "password not set,plz login with sms code"
			c.JSON(http.StatusOK, result)
			return
		}
	    fmt.Println(1)
fmt.Println(form.Password)
		err = auth.Passwords.Check(user.Pwd, form.Password)
		if err != nil {
			// fail count plus one
			err = auth.FailCountPlusOne(form.MobileNumber)
			if err != nil {
				result.Ret = 0
				result.Msg = "system error"
				c.JSON(http.StatusOK, result)
				return
			}
			result.Ret = Password_Wrong
			result.Msg = "password wrong"
			c.JSON(http.StatusOK, result)
			return
		}

		authResult, err := auth.LoginUser(user.Id, form.From, form.UniqueId)
		if err != nil {
			result.Ret = 0
			result.Msg = "login failed"
			c.JSON(http.StatusOK, result)
			return
		}

		// update last login ip
		go models.Users.UpdateLastLoginIp(user.Id, clientIp)

		// @todo 清楚登录失败记录
		result.Ret = 1
		result.UserId = user.Id
		result.IsFresh = user.IsFresh
		result.Token = authResult.Token
		result.Expire = authResult.Expire
		result.RefreshToken = authResult.RefreshToken
		c.JSON(http.StatusOK, result)
		return
	}
	result.Ret = 0
	result.Msg = "Params invaild"
	c.JSON(http.StatusOK, result)
}

func Ttest(c *gin.Context)  {
    _conf, ok1 := c.Get("apiconfig")
    cConf, ok2 := _conf.(*config.ApiConfig)
    if !ok1 || !ok2 {

    }

    fmt.Println(cConf)
}