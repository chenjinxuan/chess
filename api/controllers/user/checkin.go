package c_user

import (
	"chess/common/config"
	"chess/common/define"
	"chess/common/helper"
	"chess/common/log"
	"chess/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type CheckinResult struct {
	define.BaseResult
	CheckinDays     int    `json:"checkin_days" description:"签到天数"`
	LastCheckinTime string `json:"last_checkin_time" description:"上次签到时间"`
}
type CheskinParams struct {
	IsMore int `form:"is_more"`
}

// @Title 用户签到
// @Description 用户签到
// @Summary 用户签到
// @Accept json
// @Param   user_id     path    int  true        "用户id"
// @Param   token     query    string  true        "token"
// @Param   is_more     query    string  true        "is_more 领取额外奖励的时候为1"
// @Success 200 {object} c_user.CheckinResult
// @router /user/{user_id}/checkin [get]
func Checkin(c *gin.Context) {
	var result CheckinResult
	var params CheskinParams
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		result.Msg = "bind fail ."
		c.JSON(http.StatusOK, result)
		return
	}
	if err := c.Bind(&params); err != nil {
		result.Msg = "bind fail ."
		c.JSON(http.StatusOK, result)
		return
	}
	//获取签到信息

	checkDay, checkTime, err := models.Users.CheckinInfo(userId)
	if err != nil {
		log.Errorf("models.Users.CheckinInfo", err)
		result.Msg = "already checkin today"
		c.JSON(http.StatusOK, result)
		return
	}
	//
	nowTime := time.Now()
	lastCheckTime := checkTime.Format(define.FormatDate)
	if lastCheckTime == nowTime.Format(define.FormatDate) {
		if params.IsMore == 1 && checkDay > models.CHECKIN_DAYS_REWARD_MORE_SEVEN { //如果是领取额外奖励
			//判断今天是否已经领取额外奖励
			isMore, err := models.UsersCheckinReceive.GetIsMore(userId)
			if err != nil {

			}
			if isMore > 0 { //今天已经领取
				result.Ret = -1
				result.Msg = "already receive more"
				c.JSON(http.StatusOK, result)
				return
			}
			err = receive(userId, models.CHECKIN_DAYS_REWARD_MORE)
			if err != nil {
				result.Msg = "receive fail ."
				c.JSON(http.StatusOK, result)
				return
			}
			result.Ret = 1
			c.JSON(http.StatusOK, result)
			return

		} else {
			result.Ret = -1
			result.Msg = "already checkin today"
			c.JSON(http.StatusOK, result)
			return
		}

	}
	//签到
	yesterday := helper.GetYesterdayDate()
	if lastCheckTime != yesterday {
		checkDay = 0
	}
	checkDay = checkDay + 1
	err = models.Users.Checkin(userId, checkDay, nowTime.Format(define.FormatDatetime))
	if err != nil {
		log.Errorf("models.Users.Checkin", err)
		result.Msg = "already checkin today"
		c.JSON(http.StatusOK, result)
		return
	}
	//签到奖励
	moreDay := checkDay
	if checkDay > models.CHECKIN_DAYS_REWARD_MORE_SEVEN { //超过7天,,,,奖励为第八天的奖励
		moreDay = models.CHECKIN_DAYS_REWARD_MORE_EIGHT
	}
	err = receive(userId, moreDay)
	if err != nil {
		log.Errorf("receive reward fail", err)
		result.Msg = "receive reward fail"
		c.JSON(http.StatusOK, result)
		return
	}
	result.CheckinDays = checkDay
	result.LastCheckinTime = nowTime.Format(define.FormatDatetime)
	result.Ret = 1
	c.JSON(http.StatusOK, result)
	return
}

func receive(userId, days int) error {
	checkinData, err := models.CheckinDaysReward.Get(days)
	if err != nil {
		return err
	}
	//判断奖励类型,1 金币,给用户加金币u
	if checkinData.Type == models.CHECKIN_DAYS_REWARD_TYPE_GOLD {
		err = models.UsersWallet.AddBlance(userId, checkinData.Number)
		if err != nil {
			return err
		}
	}

	//领取记录表
	var isMore int
	if days == models.CHECKIN_DAYS_REWARD_MORE {
		isMore = 1
	}
	var checkinReceive = new(models.UsersCheckinReceiveModel)
	checkinReceive.UserId = userId
	checkinReceive.IsMore = isMore
	checkinReceive.CheckinDaysRewardId = checkinData.Id
	err = checkinReceive.Insert()
	return err
}

type CheckinListResult struct {
	define.BaseResult
	AlreadyCheckin  int                             `json:"already_checkin" description:"已经签到天数"`
	LastCheckinTime string                          `json:"last_checkin_time" description:"上次签到时间"`
	IsMore          int                             `json:"is_more"  description:"是否已经领取7天之后的额外奖励 1为已领取"`
	List            []models.CheckinDaysRewardModel `json:"list"`
}

// @Title 用户签到列表
// @Description 用户签到列表
// @Summary 用户签到列表
// @Accept json
// @Param   user_id     path    int  true        "用户id"
// @Param   token     query    string  true        "token"
// @Success 200 {object} c_user.CheckinListResult
// @router /user/{user_id}/checkin/list [get]
func CheckinList(c *gin.Context) {
	var result CheckinListResult
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		result.Msg = "bind fail ."
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
	//获取签到信息
	checkDay, checkTime, err := models.Users.CheckinInfo(userId)
	if err != nil {
		log.Errorf("models.Users.CheckinInfo", err)
		result.Msg = "already checkin today"
		c.JSON(http.StatusOK, result)
		return
	}
	lastCheckTime := checkTime.Format(define.FormatDate)
	t2, _ := time.Parse("2006-01-02", lastCheckTime)
	yesterday := helper.GetYesterdayDate()
	t1, _ := time.Parse("2006-01-02", yesterday)
	if t2.Before(t1) {
		checkDay = 0
	}
	if checkDay > models.CHECKIN_DAYS_REWARD_MORE_SEVEN {
		//判断是否已经领取额外奖励
		isMore, err := models.UsersCheckinReceive.GetIsMore(userId)
		if err != nil {
			log.Errorf("models.UsersCheckinReceive.GetIsMore", err)
			result.Msg = "get isMore fail."
			c.JSON(http.StatusOK, result)
			return
		}
		if isMore > 0 { //今天已经领取
			result.IsMore = 1
		}
	}
	result.AlreadyCheckin = checkDay
	result.LastCheckinTime = lastCheckTime
	list, err := models.CheckinDaysReward.GetAll()
	if err != nil {
		log.Errorf("models.CheckinDaysReward.GetAll err(%s)", err)
		result.Msg = "already checkin today"
		c.JSON(http.StatusOK, result)
		return
	}
	for _, v := range list {
		v.Image = cConf.Backend.ImageDomainPrefix + v.Image
		result.List = append(result.List, v)
	}
	result.Ret = 1
	c.JSON(http.StatusOK, result)
	return
}
