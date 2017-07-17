package c_auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"treasure/components/auth"
	"treasure/components/input"
	"treasure/components/sms"
	"treasure/config"
	"treasure/define"
	"treasure/models"
)

type PasswordResetParams struct {
	MobileNumber string `json:"mobile_number" form:"mobile_number" binding:"required"`
	Password     string `json:"password" form:"password"`
	Code         string `json:"code" form:"code" binding:"required" `
	From         string `json:"from" form:"from"`
}

type PasswordResetResult struct {
	define.BaseResult
}

func PasswordReset(c *gin.Context) {
	var result PasswordResetResult
	var form PasswordResetParams

	_conf, ok1 := c.Get("config")
	cConf, ok2 := _conf.(*config.Config)
	if !ok1 || !ok2 {
		result.Msg = "Get config fail."
		c.JSON(http.StatusOK, result)
		return
	}

	var user = new(models.UsersModel)
	var err error

	if input.BindJSON(c, &form, cConf) == nil {

		result.Ret, result.Msg, err = sms.CheckCode(form.MobileNumber, form.Code, sms.SMS_PWD_RESET, cConf.Sms)
		if err != nil {
			// 验证不通过
			c.JSON(http.StatusOK, result)
			return
		}

		pwdRet := auth.Passwords.CheckPasswordStrong(form.Password)
		if pwdRet != 1 {
			result.Ret = pwdRet
			result.Msg = "password not strong engough"
			c.JSON(http.StatusOK, result)
			return
		}
		// hash password
		hp, err := auth.Passwords.Hash(form.Password)
		if err != nil {
			result.Ret = 0
			result.Msg = "server error"
			c.JSON(http.StatusOK, result)
			return
		}
		// get user
		err = models.Users.GetByMobileNumber(form.MobileNumber, user)
		if err != nil {
			result.Ret = 0
			result.Msg = "server error"
			c.JSON(http.StatusOK, result)
			return
		}
		err = user.UpdatePassword(hp)
		if err != nil {
			result.Ret = 0
			result.Msg = "server error"
			c.JSON(http.StatusOK, result)
			return
		}
		result.Ret = 1
		result.Msg = "ok"
		c.JSON(http.StatusOK, result)
		return

	}
	result.Ret = 0
	result.Msg = "Params invaild."
	c.JSON(http.StatusOK, result)
	return
}
