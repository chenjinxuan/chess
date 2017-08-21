package c_user

import (
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

// @Title 用户签到
// @Description 用户签到
// @Summary 用户签到
// @Accept json
// @Param   user_id     path    int  true        "用户id"
// @Param   token     query    string  true        "token"
// @Success 200 {object} c_user.CheckinResult
// @router /user/{user_id}/checkin [get]
func Checkin(c *gin.Context) {
	var result CheckinResult
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
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
		log.Info("already checkin today")
		result.Ret = -1
		result.Msg = "already checkin today"
		c.JSON(http.StatusOK, result)
		return
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

	result.CheckinDays = checkDay
	result.LastCheckinTime = nowTime.Format(define.FormatDatetime)
	result.Ret = 1
	c.JSON(http.StatusOK, result)
	return
}
