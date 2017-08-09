//@TODO 提取流程

package tp

import (
	"errors"
	"chess/api/components/tp/qq"
	"chess/api/components/tp/wechat"
	"chess/api/helper"
	"chess/api/log"
	"chess/models"
    "fmt"
)

const (
	QQ     = "qq"
	Weibo  = "weibo"
	Wechat = "wechat"
)

// QQ Weibo Wechat 第三方登录处理接口
// QQ Weibo 通过token处理
// Wechat 通过code换取token处理
// @param token/code  ip  @return userid 第三方openid绑定的userid   msg 错误信息/OK  err

func LoginByQQ(token string, ip, channel, from string, client *qqsdk.Client) (isNew bool, user models.UsersModel, msg string, err error) {
	TpType := QQ
	appid, err := qqsdk.GetAppId(token)
	// appid 获取出错 token有误
	if err != nil {
		msg = "token wrong"
		return
	}

	// token 非官方appid获取
	if appid != client.GetAppId() {
		msg = "token wrong"
		err = errors.New("token is not verify")
		return
	}

	var openid, unionid string
	openid, unionid, err = qqsdk.GetOpenIdAndUnionId(token)
	if err != nil {
		openid, err = qqsdk.GetOpenId(token)
		if err != nil {
			msg = "cant get openid and unionid"
			return
		}
	}

	// 判断openid是否存在
	UserId, err := models.UsersTp.IsReg(openid, TpType)
	user.Id = UserId

	// openid,存在,返回user_id
	if err == nil {
		if unionid != "" {
			err = models.UsersTp.UpdateWxUnionIdByUid(UserId, unionid)
			if err != nil {
				log.Log.Error(err)
			}
		}
		msg = "ok"
		return
	}

	// 没有open id, 判断union id
	if unionid != "" {
		var tpId int
		tpId, _, err = models.UsersTp.CheckWxUnionId(unionid, QQ)
		if err == nil {
			// 找得到union id
			var tpUser models.UsersTpModel
			cerr := models.UsersTp.GetId(tpId, &tpUser)
			err = cerr
			if err != nil {
				msg = "cant get user"
				return
			}
			// 更新openid
			err = models.UsersTp.UpdateOpenidIdById(tpId, openid)
			if err != nil {
				log.Log.Error("store openid error", err)
			}

			user.Id = tpUser.UserID
			return
		}
	}

	// openid,union id都不存在,创建用户
	// 获取QQ信息
	var qqUser = new(qqsdk.UserInfo)
	qqUser, err = qqsdk.GetUserInfo(token, openid, client.GetAppId())
	if err != nil {
		msg = "Could get qq user info"
		return
	}
	// 插入users表
	user.RegIp = ip
	nickname := qqUser.Nickname
	user.Nickname = helper.ConvertNickname(nickname)
	user.Channel = channel
	user.LastLoginIp = user.RegIp
	user.Avatar = qqUser.Figureurl_qq_2
	user.Type = models.TYPE_QQ
        user.AppFrom = from
	userid, err := models.Users.Insert(&user)
	if err != nil {
		msg = "Could not create new user."
		return isNew, user, msg, err
	}

	user.Id = userid
	// 插入wallet
	// Init user wallet
	err = models.UsersWallet.Init(userid)
	if err != nil {
		log.Log.Error(err)
		return isNew, user, msg, err
	}

	// 插入users_tp表
	userQQ := new(models.UsersTpModel)
	userQQ.Type = TpType
	userQQ.OpenID = openid
	userQQ.WxUnionId = unionid
	userQQ.UserID = userid
	_, err = models.UsersTp.Insert(userQQ)
	if err != nil {
		log.Log.Error(err)
		msg = "Could not create new qq user."
		return isNew, user, msg, err
	}

	isNew = true
	return
}

//
//func LoginByWeibo(token string, ip, channel, from string, cConf *config.Config) (isNew bool, user models.UsersModel, msg string, err error) {
//	TpType := Weibo
//	weiboApi := weibo.NewApi(token, cConf.Tp.Weibo.AppId, cConf.Tp.Weibo.AppSecret)
//
//	tokenInfo, err := weibo.GetTokenInfo(token, weiboApi)
//	// appid 获取出错 token有误
//	if err != nil {
//		msg = "token wrong"
//		return
//	}
//
//	// token 非官方appid获取
//	if tokenInfo.AppKey != cConf.Tp.Weibo.AppId {
//		msg = "token wrong"
//		err = errors.New("token is not verify")
//		return
//	}
//
//	openid := strconv.Itoa(tokenInfo.UID)
//
//	UserId, err := models.UsersTp.IsReg(openid, TpType)
//	if err != nil {
//		// openid不存在，创建用户并更新userid
//
//		// 获取微博信息
//		//weiboUser := weibo.NewWeigoUser()
//		weiboUser, wberr := weibo.GetUserShowByUid(token, openid, cConf.Tp.Weibo.AppId, weiboApi)
//		err = wberr
//		// @TODO 错误处理
//		if err != nil {
//			msg = "Could get weibo user info"
//			return
//		}
//		// 插入users表
//		//var user = new(models.UsersModel)
//		user.RegIp = ip
//		nickname := weiboUser.Name
//		nickname = helper.ConvertNickname(nickname)
//		//nickname = helper.ConverUnsupportStr(nickname)
//		//if len(nickname) > 20 {
//		//	nickname = nickname[:20]
//		//}
//		user.Nickname = nickname
//		user.Channel = channel
//		user.LastLoginIp = user.RegIp
//		user.Avatar = weiboUser.Avatar_large
//		user.Type = models.TYPE_WEIBO
//	        user.AppFrom = from
//		userid, err := models.Users.Insert(&user)
//		if err != nil {
//			msg = "Could not create new user."
//			return isNew, user, msg, err
//		}
//
//		user.Id = userid
//		// 插入wallet
//		// Init user wallet
//		err = models.UsersWallet.Init(userid)
//		if err != nil {
//			log.Log.Error(err)
//			return isNew, user, msg, err
//		}
//		// 插入users_tp表
//		tpUser := new(models.UsersTpModel)
//		tpUser.Type = TpType
//		tpUser.OpenID = openid
//		tpUser.UserID = userid
//		_, err = models.UsersTp.Insert(tpUser)
//		// @TODO 错误处理
//		if err != nil {
//			log.Log.Error(err)
//			msg = "Could not create new weibo user."
//			return isNew, user, msg, err
//		}
//		isNew = true
//	} else {
//		user.Id = UserId
//	}
//	return
//}

func LoginByWechat(code string, ip, channel, from string, client *wechat.Client) (isNew bool, user models.UsersModel, msg string, err error) {
	TpType := Wechat

	//tokenInfo := client.Token
	tokenInfo, err := client.GetTokenByCode(code)
	// appid 获取出错 token有误
	if err != nil {
		//log.Log.Error(err)
		msg = "token wrong"
		return
	}
	defer log.Log.Debug("token info", tokenInfo)

	// 判断open id
	openid := tokenInfo.OpenId
	UserId, err := models.UsersTp.IsReg(openid, TpType)
	user.Id = UserId

	// openid,存在,返回user_id
	if err == nil {
		err = models.UsersTp.UpdateWxUnionIdByUid(UserId, tokenInfo.UnionId)
		if err != nil {
			log.Log.Error(err)
		}
		msg = "ok"
		return
	}

	// 没有open id, 判断union id
	tpId, _, err := models.UsersTp.CheckWxUnionId(tokenInfo.UnionId, Wechat)
	if err == nil {
		// 找得到union id
		var tpUser models.UsersTpModel
		cerr := models.UsersTp.GetId(tpId, &tpUser)
		err = cerr
		if err != nil {
			msg = "cant get user"
			return
		}
		// 更新openid
		err = models.UsersTp.UpdateOpenidIdById(tpId, tokenInfo.OpenId)
		if err != nil {
			log.Log.Error("store openid error", err)
		}

		user.Id = tpUser.UserID
		return
	}

	// openid,union id都不存在,创建用户
	var wechatUser wechat.UserInfo
	// openid不存在，创建用户并更新userid
	// 获取Wechat
	//wechatUser := wechat.GetNewWechatUserInfo()
	//werr := client.GetUserInfo(wechatUser, "zh-cn")
	var werr error
	wechatUser, werr = client.GetUserInfoByToken(tokenInfo)
	err = werr
	// @TODO 错误处理
	if err != nil {
		msg = "Could get wechat user info"
		return
	}
	// 插入users表
	//var user = new(models.UsersModel)
	user.RegIp = ip
	nickname := wechatUser.Nickname
	nickname = helper.ConvertNickname(nickname)
	//nickname = helper.ConverUnsupportStr(nickname)
	//if len(nickname) > 20 {
	//	nickname = nickname[:20]
	//}
	user.Nickname = nickname
	user.LastLoginIp = user.RegIp
	user.Gender = wechatUser.Sex
	user.Channel = channel
	user.Avatar = wechatUser.HeadImageURL
	user.Type = models.TYPE_WECHAT
        user.AppFrom = from
	userid, err := models.Users.Insert(&user)
	if err != nil {
	    fmt.Println(err)
		msg = "Could not create new user."
		return
	}

	user.Id = userid
	// 插入wallet
	// Init user wallet
	err = models.UsersWallet.Init(userid)
	if err != nil {
		log.Log.Error(err)
		msg = "system error"
		return
	}

	// 插入users_tp表
	tpUser := new(models.UsersTpModel)
	tpUser.Type = TpType
	tpUser.OpenID = openid
	tpUser.UserID = userid
	tpUser.WxUnionId = wechatUser.UnionId
	_, err = models.UsersTp.Insert(tpUser)
	// @TODO 错误处理
	if err != nil {
		log.Log.Error(err)
		msg = "Could not create new wechat user."
		return
	}
	isNew = true

	msg = "ok"
	return
}
