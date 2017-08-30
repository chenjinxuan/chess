package c_auth

import (
	"chess/api/components/auth"
	"chess/api/components/input"
	"chess/api/components/tp"
	"chess/api/redis"
	"chess/common/config"
	"chess/common/define"
	"chess/common/log"
	"chess/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TokenInfoParams struct {
	UserId   int    `json:"user_id" binding:"required" description:"用户id"`
	From     string `json:"from" description:"来源"`
	UniqueId string `json:"unique_id" description:"唯一设备标识"`
	Token    string `json:"token" binding:"required" description:"token"`
	Ver      int    `json:"ver" binding:"required" description:"版本号"`
}

type TokenInfoResult struct {
	define.BaseResult
	Expire int64 `json:"expire,omitempty"`
}

// @Title 获取token信息
// @Description 获取token信息
// @Summary 获取token信息
// @Accept json
// @Param   body     body    c_auth.TokenInfoParams  true        "post 数据"
// @Success 200 {object} c_auth.TokenInfoResult
// @router /auth/token/info [post]
func TokenInfo(c *gin.Context) {
	var result TokenInfoResult
	var post TokenInfoParams

	_conf, ok1 := c.Get("config")
	cConf, ok2 := _conf.(*config.ApiConfig)
	if !ok1 || !ok2 {
		result.Msg = "Get config fail."
		c.JSON(http.StatusOK, result)
		return
	}

	if input.BindJSON(c, &post, cConf) == nil {
		//defer log.Log.Debug("token-info-result-debug:", result)
		//log.Log.Info(post)
		post.From = strings.ToLower(post.From)
		// Check token
		defer log.Debug(result)
		//判断黑名单
		msg, err := api_redis.Redis.Login.GetInt(post.Token)
		if err == nil {
			result.Ret = define.AuthALreadyLogin
			result.Msg = define.AuthMsgMap[msg]
			c.JSON(http.StatusOK, result)
			return
		}
		loginData, err := auth.AuthLoginToken(post.Token, cConf.TokenSecret)
		if err != nil {
			log.Error(err)
			result.Ret = define.AuthFailedStatus
			result.Msg = auth.AuthFailed.Error()
			c.JSON(http.StatusOK, result)
			return
		}

		if strconv.Itoa(post.UserId) != loginData {
			log.Debug("strconv fail")
			result.Msg = auth.AuthFailed.Error()
			c.JSON(http.StatusOK, result)
			return
		}

		// Get the session
		session, err := models.Session.Get(post.UserId)
		if err != nil {
			log.Debugf("userid:%d,from:%s,uniqueId:%s", post.UserId, post.From, post.UniqueId)
			result.Msg = auth.AuthFailed.Error()
			c.JSON(http.StatusOK, result)
			return
		}
		//判断过期时间12小时
		now := time.Now().Unix()
		if session.Token.Expire-now < 60*60*12 {
			result.Ret = define.AuthReToken
			result.Msg = define.AuthMsgMap[define.AuthReToken]
			c.JSON(http.StatusOK, result)
			return
		}
		// TODO: 检查黑名单

		//var checkConfig config.TokenInfoCheckDetail
		//if post.From == "ios" {
		//	checkConfig = config.C.TokenInfoCheck.Ios
		//}
		//if post.From == "android" {
		//	checkConfig = config.C.TokenInfoCheck.Android
		//}
		//
		//if !checkConfig.Check || post.Ver < checkConfig.Min || post.Ver > checkConfig.Max {
		//	result.Ret = 1
		//	result.Expire = session.Token.Expire
		//	c.JSON(http.StatusOK, result)
		//	return
		//}

		// wechat union id check
		var tpUser models.UsersTpModel
		err = models.UsersTp.GetByUid(post.UserId, &tpUser)
		if err == nil {
			// 找到tp用户, 判断为微信且union id为空
			if tpUser.Type == tp.Wechat && tpUser.WxUnionId == "" {
				result.Ret = -1
				result.Expire = session.Token.Expire
				c.JSON(http.StatusOK, result)
				return
			}
		}

		result.Ret = 1
		result.Expire = session.Token.Expire
		c.JSON(http.StatusOK, result)
		return
	} else {
		log.Debug("params invalid")
	}

	result.Msg = "Params invaild."
	c.JSON(http.StatusOK, result)
}
