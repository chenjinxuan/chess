package c_auth

import (
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"chess/api/components/auth"
	"chess/api/components/input"
	"chess/api/components/tp"
	"chess/api/components/tp/qq"
	"chess/api/components/tp/wechat"
	"chess/api/components/user_init"
	"chess/common/config"
	"chess/api/define"
	"chess/api/helper"
	"chess/api/log"
	"chess/models"
)

const (
	TpQQ       = 1
	TpWeibo    = 2
	TpWechat   = 3
	TpWechatH5 = 4
	TpQQH5     = 5
	TpQQInner  = 6 // QQ应用内登录
)

type TpLoginResult struct {
	define.BaseResult
	LoginResult
}

type TpLoginParams struct {
	Type         int                    `json:"type" form:"type" binding:"required"`
	Key          string                 `form:"key" json:"key" binding:"required"`
	From         string                 `form:"from" json:"from" binding:"required"`
	UniqueId     string                 `json:"unique_id" form:"unique_id"`
	Channel      string                 `json:"channel"`
	//BindingParam broker.CheckCodeParams `json:"binding_param"`
}

func TpLogin(c *gin.Context) {
	var result TpLoginResult
	var form TpLoginParams
	var userId int

	_conf, ok1 := c.Get("config")
	cConf, ok2 := _conf.(*config.ApiConfig)
	if !ok1 || !ok2 {
		result.Msg = "Get config fail."
		c.JSON(http.StatusOK, result)
		return
	}

	var isfresh int
	var isNew bool
	clientIp := helper.ClientIP(c)

	if err := input.BindJSON(c, &form, cConf); err != nil {
		log.Log.Error("BindJSON error: ", err)
		result.Msg = "wrong params"
		c.JSON(http.StatusOK, result)
		return
	}
	if form.Channel == "" {
		form.Channel = "default"
	}
	if form.From == "" {
		form.From = "default"
	}
	form.From = strings.ToLower(form.From)
	defer func() {
		// Debug
		log.Log.WithFields(logrus.Fields{
			"key":    form.Key,
			"type":   form.Type,
			"form":   form.From,
			"result": result,
		}).Debug("c_auth.login.tp")
	}()


	var userInfo models.UsersModel
	switch form.Type {
	case TpQQ, TpQQH5, TpQQInner:
		client := qqsdk.NewClient(cConf.Tp.QQ.AppId, cConf.Tp.QQ.AppSecret, cConf.Tp.QQ.RedirectUrl)
		if form.Type == TpQQH5 || form.Type == TpQQInner {
			client = qqsdk.NewClient(cConf.Tp.H5QQ.AppId, cConf.Tp.H5QQ.AppSecret, cConf.Tp.H5QQ.RedirectUrl)
			if form.Type == TpQQInner { // QQ应用内登录，先用Code获取accesstoken
				token, _ := client.GetAccessToken(form.Key, "")
				form.Key = token.AccessToken
			}
		}

		isnew, user, msg, err := tp.LoginByQQ(form.Key, clientIp, form.Channel, form.From, client)
		if err != nil {
			log.Log.Error("qq login failed:", err)
			result.Ret = 0
			result.Msg = msg
			c.JSON(http.StatusOK, result)
			return
		}

		isNew = isnew
		userId = user.Id
		userInfo = user
        //
	//case TpWeibo:
	//	isnew, user, msg, err := tp.LoginByWeibo(form.Key, clientIp, form.Channel, form.From, cConf)
	//	if err != nil {
	//		log.Log.Error("weibo login failed", err)
	//		result.Ret = 0
	//		result.Msg = msg
	//		c.JSON(http.StatusOK, result)
	//		return
	//		break
	//	}
	//	isNew = isnew
	//	userId = user.Id
	//	userInfo = user

	case TpWechat, TpWechatH5:
		client := wechat.NewClient(cConf.Tp.Wechat.AppId, cConf.Tp.Wechat.AppSecret)
		if form.Type == TpWechatH5 {
			client = wechat.NewClient(cConf.Tp.H5Wechat.AppId, cConf.Tp.H5Wechat.AppSecret)
		}

		isnew, user, msg, err := tp.LoginByWechat(form.Key, clientIp, form.Channel, form.From, client)
		if err != nil {
			log.Log.Error("wechat login failed", err)
			result.Ret = 0
			result.Msg = msg
			c.JSON(http.StatusOK, result)
			return
			break
		}
		isNew = isnew
		userId = user.Id
		userInfo = user

	default:
		result.Msg = "unsupport type"
		c.JSON(http.StatusOK, result)
		return
	}

	// 新用户处理
	extra := make(map[string]interface{})
	extra["from"] = form.From
	extra["unique_id"] = form.UniqueId
	extra["idfv"] = c.Query("idfv")
	extra["idfa"] = c.Query("idfa")
    _=userInfo
	//isfresh, err := user_init.UserInit(userInfo, extra, cConf)
	//if err != nil {
	//	log.Log.Error(err)
	//}

	// 绑定邀请码
	if isNew {
		//go broker.Binding(userId, form.BindingParam, isfresh, bindType, c.Query("from"), cConf)

		//更新设备信息
		err := user_init.DeviceInit(userId, form.From, form.UniqueId, c.Query("idfv"), c.Query("idfa"))
		if err != nil {
			log.Log.Error(err)
		}

	} else {
		go models.Users.UpdateLastLoginIp(userId, clientIp)
	}

	authResult, err := auth.LoginUser(userId, form.From, form.UniqueId)
	if err != nil {
		result.Ret = 0
		result.Msg = "login failed"
		c.JSON(http.StatusOK, result)
		return
	}
	result.Ret = 1
	result.Token = authResult.Token
	result.Expire = authResult.Expire
	result.RefreshToken = authResult.RefreshToken
	result.UserId = authResult.UserId
	result.IsFresh = isfresh
	c.JSON(http.StatusOK, result)
	return

}
