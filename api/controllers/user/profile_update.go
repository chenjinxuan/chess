package c_user

import (
	"chess/api/components/convert"
	"chess/api/components/input"
	"chess/api/components/profile"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	//"chess/api/components/sms"
	"chess/common/config"
	"chess/common/define"
	"chess/common/helper"
	"chess/common/log"
	"chess/models"
)

const (
	Old_Pwd_Wrong = -11

	// /profile/gender
	NoGender    = 0
	Male        = 1
	Female      = 2
	OtherGender = 3

	// Nickname
	Nickname_Wrong       = -41 // 含有敏感词或者非法字符
	Nickname_Check_Error = -42 // 敏感词检查出现问题
	System_Wrong         = -43 // 更新出问题

	// mobile num
	Mobile_Exist     = -51
	Mobile_Not_Empty = -52
	Mobile_Binded    = -53 //该手机好已绑定第三方账号
)

type ProfileNicknameParams struct {
	Nickname string `form:"nickname" json:"nickname" binding:"required" description:"昵称"`
}

type ProfileMobileParams struct {
	MobileNumber string `json:"mobile_number" form:"mobile_num" binding:"required" description:"手机号"`
	Code         string `json:"code" form:"code" description:"短信验证码"`
}

type ProfileAvatarParams struct {
	AvatarFileName string `json:"avatar_filename" form:"avatar_filename" binding:"required" description:"头像链接"`
}

type ProfileGenderParams struct {
	Gender int `json:"gender" form:"gender" description:"性别 0 未知,1男,2女"`
}

type ProfileMobileResult struct {
	define.BaseResult
	DiamondBalance int `json:"diamond_balance" description:"钻石余额"`
	Balance        int `json:"balance" description:"金币余额"`
}

type ProfileAvatarResult struct {
	define.BaseResult
	Avatar string `json:"avatar" description:"头像链接"`
}

// @Title 更新昵称
// @Description 更新昵称
// @Summary 更新昵称
// @Accept json
// @Param   body     body    c_user.ProfileNicknameParams  true        "post 数据"
// @Success 200 {object} define.BaseResult
// @router /user/{user_id}/profile/nickname [post]
func ProfileNicknameUpdate(c *gin.Context) {
	var result define.BaseResult
	var form ProfileNicknameParams

	_conf, ok1 := c.Get("config")
	cConf, ok2 := _conf.(*config.ApiConfig)
	if !ok1 || !ok2 {
		result.Msg = "Get config fail."
		c.JSON(http.StatusOK, result)
		return
	}

	uidStr := c.Param("user_id")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		result.Ret = 0
		c.JSON(http.StatusOK, result)
		return
	}

	if input.BindJSON(c, &form, cConf) == nil {
		var user models.UsersModel
		err = models.Users.Get(uid, &user)

		nickname, isOK, err := profile.CheckNickname(form.Nickname)
		if err != nil {
			result.Ret = Nickname_Check_Error
			result.Msg = "system error"
			c.JSON(http.StatusOK, result)
			return
		}
		if !isOK {
			result.Ret = Nickname_Wrong
			result.Msg = "nickname include bad char"
			c.JSON(http.StatusOK, result)
			return
		}
		err = user.UpdateNickname(uid, helper.ConverUnsupportStr(nickname))

		if err != nil {
			result.Ret = 0
			result.Msg = "update failed"
			c.JSON(http.StatusOK, result)
			return
		}

		result.Ret = 1
		result.Msg = "ok"
		c.JSON(http.StatusOK, result)

	} else {
		result.Ret = 0
		result.Msg = "wrong params"
		c.JSON(http.StatusOK, result)
		return
	}
}

// @Title 绑定手机
// @Description 绑定手机
// @Summary 绑定手机
// @Accept json
// @Param   body     body    c_user.ProfileMobileParams  true        "post 数据"
// @Success 200 {object} c_user.ProfileMobileResult
// @router /user/{user_id}/profile/mobile [post]
func ProfileMobile(c *gin.Context) {
	var result ProfileMobileResult
	var form ProfileMobileParams
	var user models.UsersModel

	_conf, ok1 := c.Get("config")
	cConf, ok2 := _conf.(*config.ApiConfig)
	if !ok1 || !ok2 {
		result.Msg = "Get config fail."
		c.JSON(http.StatusOK, result)
		return
	}

	//clientIp := helper.ClientIP(c)

	uidStr := c.Param("user_id")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		result.Ret = 0
		c.JSON(http.StatusOK, result)
		return
	}

	if input.BindJSON(c, &form, cConf) == nil {
		err = models.Users.Get(uid, &user)
		if err != nil {
			result.Ret = 0
			c.JSON(http.StatusOK, result)
			return
		}
		//if user.MobileNumber != "" {
		//	// 校验修改绑定手机验证码
		//	result.Ret, result.Msg, err = sms.CheckChangeMobileCode(user.MobileNumber)
		//	if err != nil {
		//		c.JSON(http.StatusOK, result)
		//		return
		//	}
		//}

		// 判断手机号码是否被注册了
		// Query the user is exists
		err := models.Users.GetByMobileNumber(form.MobileNumber, &user)
		if err == nil {
			var userTp models.UsersTpModel
			//判断该手机是否已绑定第三方账号
			err := models.UsersTp.GetByMobile(form.MobileNumber, &userTp)
			if err == nil {
				result.Ret = Mobile_Binded
				result.Msg = "this mobile has binded"
				c.JSON(http.StatusOK, result)
				return
			}
			var balance, diamond_balance int
			balance, diamond_balance, err = models.UsersWallet.GetBalanceByMobile(form.MobileNumber)
			log.Debug(balance)
			if err != nil {
				result.Ret = 0
				c.JSON(http.StatusOK, result)
				return
			}
			result.Balance = balance
			result.DiamondBalance = diamond_balance
			result.Ret = Mobile_Exist
			result.Msg = "this mobile is exist"
			c.JSON(http.StatusOK, result)
			return
		}

		// 校验验证码
		//result.Ret, result.Msg, err = sms.CheckCode(form.MobileNumber, form.Code, sms.SMS_CONTACT_MOBIlE, cConf)
		//if err != nil {
		//	// 验证不通过
		//	c.JSON(http.StatusOK, result)
		//	return
		//}

		// update mobile
		err = user.UpdateMobile(form.MobileNumber)
		if err != nil {
			result.Ret = 0
			result.Msg = "server error"
			c.JSON(http.StatusOK, result)
			return
		}

		result.Ret = 1
		result.Msg = "OK"
		c.JSON(http.StatusOK, result)
		return
	}
	result.Ret = 0
	result.Msg = "Params invaild."
	c.JSON(http.StatusOK, result)
	return
}

// @Title 更新头像
// @Description 更新头像
// @Summary 更新头像
// @Accept json
// @Param   body     body    c_user.ProfileAvatarParams  true        "post 数据"
// @Success 200 {object} c_user.ProfileAvatarResult
// @router /user/{user_id}/profile/avatar [post]
func ProfileAvatar(c *gin.Context) {
	var result ProfileAvatarResult
	var form ProfileAvatarParams
	var user models.UsersModel

	_conf, ok1 := c.Get("config")
	cConf, ok2 := _conf.(*config.ApiConfig)
	if !ok1 || !ok2 {
		result.Msg = "Get config fail."
		c.JSON(http.StatusOK, result)
		return
	}

	uidStr := c.Param("user_id")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		result.Ret = 0
		c.JSON(http.StatusOK, result)
		return
	}
	err = models.Users.Get(uid, &user)
	if err != nil {
		result.Ret = 0
		c.JSON(http.StatusOK, result)
		return
	}
	if input.BindJSON(c, &form, cConf) == nil {
		avatar := form.AvatarFileName
		// 不允许传url
		if helper.IsUrl(avatar) {
			result.Ret = 0
			result.Msg = "server error"
			c.JSON(http.StatusOK, result)
			return
		}
		err = user.UpdateAvatar(avatar)
		if err != nil {
			result.Ret = 0
			result.Msg = "server error"
			c.JSON(http.StatusOK, result)
			return
		}
		result.Ret = 1
		result.Msg = "OK"
		result.Avatar = convert.ToFullAvatarUrl(avatar, cConf.Storage.QiniuAvatarUrl, cConf.User.DefaultAvatar)
		c.JSON(http.StatusOK, result)
		return
	}
	result.Ret = 0
	result.Msg = "Params invaild."
	c.JSON(http.StatusOK, result)
	return
}

// @Title 更新性别
// @Description 更新性别
// @Summary 更新性别
// @Accept json
// @Param   body     body    c_user.ProfileGenderParams  true        "post 数据"
// @Success 200 {object} define.BaseResult
// @router /user/{user_id}/profile/gender [post]
func ProfileGender(c *gin.Context) {
	var result define.BaseResult
	var form ProfileGenderParams
	var user models.UsersModel

	_conf, ok1 := c.Get("config")
	cConf, ok2 := _conf.(*config.ApiConfig)
	if !ok1 || !ok2 {
		result.Msg = "Get config fail."
		c.JSON(http.StatusOK, result)
		return
	}

	uidStr := c.Param("user_id")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		result.Ret = 0
		c.JSON(http.StatusOK, result)
		return
	}
	err = models.Users.Get(uid, &user)
	if err != nil {
		result.Ret = 0
		c.JSON(http.StatusOK, result)
		return
	}
	if input.BindJSON(c, &form, cConf) == nil {

		if form.Gender != NoGender && form.Gender != Male && form.Gender != Female && form.Gender != OtherGender {
			result.Ret = 0
			result.Msg = "gender not allow "
			c.JSON(http.StatusOK, result)
			return
		}

		err = user.UpdateGender(form.Gender)
		if err != nil {
			result.Ret = 0
			result.Msg = "server error"
			c.JSON(http.StatusOK, result)
			return
		}
		result.Ret = 1
		result.Msg = "OK"
		c.JSON(http.StatusOK, result)
		return
	}
	result.Ret = 0
	result.Msg = "Params invaild."
	c.JSON(http.StatusOK, result)
	return
}
