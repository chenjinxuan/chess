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
        grpcServer "chess/api/grpc"
    pb "chess/api/proto"
    "golang.org/x/net/context"
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
	MobileNumber string `json:"mobile_number" form:"mobile_number" binding:"required" description:"手机号"`
	Password     string `json:"password" form:"password" binding:"required" description:"密码"`
	From         string `json:"from" form:"from" description:"请求来源"`
	UniqueId     string `json:"unique_id" form:"unique_id" description:"唯一标识"`
	Captcha      string `json:"captcha" description:"验证码"`
}

type LoginResult struct {
	define.BaseResult
	NeedCaptcha  bool   `json:"need_captcha" description:"是否需要验证码"`
	UserId       int    `json:"user_id" description:"用户Id"`
	IsFresh      int    `json:"fresh" description:"是否是新用户"`
	Token        string `json:"token" description:"请求授权的token"`
	Expire       int64  `json:"expire" description:"token过期时间"`
	RefreshToken string `json:"refresh_token" description:"当 Token 过期时，使用该 Token 获取一个新的"`
}
// @Title 手机号 + 密码登录
// @Description 手机号 + 密码登录
// @Summary 手机号 + 密码登录
// @Accept json
// @Param   body     body    c_auth.LoginParams  true        "post 数据"
// @Success 200 {object} c_auth.LoginResult
// @router /auth/login [post]
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
                AuthClient:=grpcServer.GetAuthGrpc()
		authResult, err := AuthClient.RefreshToken(context.Background(), &pb.RefreshTokenArgs{UserId: int32(user.Id),AppFrom:form.From,UniqueId:form.UniqueId})
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