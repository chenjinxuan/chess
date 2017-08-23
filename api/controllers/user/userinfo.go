package c_user

import (
	"chess/common/define"
	"chess/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
        "chess/api/components/convert"
        "chess/common/config"
)

type UserInfo struct {
	Id              int    `json:"id"`
	NickName        string `json:"nick_name" description:"昵称"`
	MobileNumber    string `json:"mobile_number" description:"手机号"`
	Gender          int    `json:"gender" description:"性♂别 // 0 未知  1男2女 "`
	Avatar          string `json:"avatar" description:"头像"`
	Type            int    `json:"type" description:"用户类型"`
	Status          int    `json:"status" description:"状态"`
	IsFresh         int    `json:"is_fresh" description:"是否新用户"`
	Balance         int    `json:"balance" description:"金币余额"`
	DiamondBalance  int    `json:"diamond_balance" description:"钻石余额"`
	CheckinDays     int    `json:"checkin_days" description:"签到天数"`
	LastCheckinTime string `json:"last_checkin_time" description:"上次签到时间"`
}

type UserInfoResult struct {
	define.BaseResult
	Data UserInfo `json:"data"`
}

// @Title 获取用户基本信息
// @Description 获取用户基本信息
// @Summary 获取用户基本信息
// @Accept json
// @Param   token     query    string   true        "token"
// @Param   user_id     path    int   true        "user_id"
// @Success 200 {object} c_user.UserInfoResult
// @router /user/{user_id}/info [get]
func GetUserInfo(c *gin.Context) {
	var result UserInfoResult
	UserId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		result.Msg = "bind params fail ."
		c.JSON(http.StatusOK, result)
		return
	}
	_conf, ok1 := c.Get("config")
	cConf, ok2 := _conf.(*config.ApiConfig)
	if !ok1 || !ok2 {
	    result.Msg = "Get config fail."
	    c.JSON(http.StatusOK, result)
	    return
	}
	var user = new(models.UsersModel)
	//获取信息
	err = models.Users.Get(UserId, user)
	if err != nil {
		result.Msg = "get info fail 1"
		c.JSON(http.StatusOK, result)
		return
	}
	result.Data.Id = user.Id
	result.Data.NickName = user.Nickname
	result.Data.MobileNumber = user.MobileNumber
	result.Data.Gender = user.Gender
	result.Data.Avatar = convert.ToFullAvatarUrl(user.Avatar, cConf.Storage.QiniuAvatarUrl, cConf.User.DefaultAvatar)
	result.Data.Type = user.Type
	result.Data.Status = user.Status
	result.Data.IsFresh = user.IsFresh
	result.Data.CheckinDays = user.CheckinDays
	result.Data.LastCheckinTime = user.LastCheckinTime.Format(define.FormatDatetime)
	//余额查询
	result.Data.Balance, result.Data.DiamondBalance, err = models.UsersWallet.GetBalance(UserId)
	if err != nil {
		result.Msg = "get info fail 2"
		c.JSON(http.StatusOK, result)
		return
	}
	result.Ret = 1
	c.JSON(http.StatusOK, result)
	return
}

type GetUserInfoDetailResult struct {
    define.BaseResult
    WinRate float64 `json:"win_rate" description:"胜率"`
    TotalGame int `json:"total_game" description:"总局数"`
    Cards []UserBestCards `json:"cards" description:"最大牌"`
    BestWinner int `json:"best_winner" description:"最大赢取筹码"`
    ShowdownRate float64 `json:"showdown_rate" description:"摊牌率"`
    InboundRate float64 `json:"inbound_rate" description:"入局率"`
}
type UserBestCards struct {
    Suit  int `json:"suit" description:"花色"`//程序统一标准：0是黑桃、1是红桃、2是梅花、3是方片
    Value int `json:"value" description:"大小"`//0代表‘牌2’、1代表‘牌3’...etc
}
// @Title 获取用户详细信息
// @Description 获取用户详细信息
// @Summary 获取用户详细信息
// @Accept json
// @Param   token     query    string   true        "token"
// @Param   user_id     path    int   true        "user_id"
// @Success 200 {object} c_user.GetUserInfoDetailResult
// @router /user/{user_id}/detail [get]
func GetUserInfoDetail(c *gin.Context) {
       	var result GetUserInfoDetailResult
	//UserId, err := strconv.Atoi(c.Param("user_id"))
	//if err != nil {
	//    result.Msg = "bind params fail ."
	//    c.JSON(http.StatusOK, result)
	//    return
	//}
        //获取相应的数据
    result.WinRate=0.5
    result.TotalGame = 30
    result.BestWinner = 1333
    result.ShowdownRate = 0.6
    result.InboundRate = 0.8
    var cards UserBestCards
    cards.Value = 1
    cards.Suit=1
    for  i:=0 ;i<=4 ;i++ {
	cards.Value=i
	result.Cards=append(result.Cards,cards)
    }
    result.Ret = 1
    c.JSON(http.StatusOK, result)
    return
}