package c_task

import (
	"chess/common/define"
	"chess/common/log"
	"chess/models"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type ReceiveTaskRewardParams struct {
	TaskId int `form:"task_id"  description:"任务id"`
}
type ReceiveTaskRewardResult struct {
	define.BaseResult
	Data models.UserTaskModel `json:"data" description:"替换已完成的任务"`
}

// @Title 任务达标领取
// @Description 任务达标领取
// @Summary 任务达标领取
// @Accept json
// @Param   token     query    string   true        "token"
// @Param   task_id     query    string   true        "task_id"
// @Param   user_id     path    int   true        "user_id"
// @Success 200 {object} c_task.ReceiveTaskRewardResult
// @router /task/{user_id}/receive [get]
func ReceiveTaskReward(c *gin.Context) {
	var result ReceiveTaskRewardResult
	var params ReceiveTaskRewardParams
	UserId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		result.Msg = "bind params fail ."
		c.JSON(http.StatusOK, result)
		return
	}
	if err := c.Bind(&params); err != nil {
		result.Msg = "bind fail ."
		c.JSON(http.StatusOK, result)
		return
	}
	//查出用户改任务的完成情况
	data, err := models.UserTask.Get(UserId)
	if err != nil {
		log.Errorf("models.UserTask.Get err %s", err)
		result.Msg = "get fail ."
		c.JSON(http.StatusOK, result)
		return
	}

	t := time.Now()
	var oldTask models.UserTaskModel
	for _, v := range data.List {
		if v.TaskId == params.TaskId {
			if v.RequiredNum > v.AlreadyCompleted {
				break
			}
			//判断类型//过期类型

			if v.TaskTypeExpireType == define.TodayTask {
				if time.Unix(v.LastUpdate, 0).Format(define.FormatDate) != t.Format(define.FormatDate) {
					break
				}
			}
			if v.TaskTypeExpireType == define.WeekTask {
				_, week1 := time.Unix(v.LastUpdate, 0).ISOWeek()
				_, week2 := t.ISOWeek()
				if week1 != week2 {
					break
				}
			}
			oldTask = v
		}
	}
	if oldTask.AlreadyCompleted == 0 {
		result.Msg = "required fail ."
		c.JSON(http.StatusOK, result)
		return
	}
	//查出该任务的下一级任务
	task, err := models.UserTask.GetByParentId(params.TaskId)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("models.UserTask.GetByParentId err %s", err)
		result.Msg = "get task fail ."
		c.JSON(http.StatusOK, result)
		return
	}
	task.LastUpdate = t.Unix()
	task.AlreadyCompleted = oldTask.AlreadyCompleted
	if err == sql.ErrNoRows {
		//删除旧的任务
		err = models.UserTask.RemoveByTaskId(UserId, oldTask.TaskId)
		if err != nil {
			log.Errorf("models.UserTask.RemoveByTaskId err %s", err)
			result.Msg = "delete fail ."
			c.JSON(http.StatusOK, result)
			return
		}
	} else {
		//更新任务
		err = models.UserTask.UpdateOneTask(UserId, oldTask.TaskId, task)
		if err != nil {
			log.Errorf("models.UserTask.UpdateOneTask err %s", err)
			result.Msg = "update fail ."
			c.JSON(http.StatusOK, result)
			return
		}
	}

	//插入记录表
	var taskPrize = new(models.TaskPriceReceiveModel)
	taskPrize.UserId = UserId
	taskPrize.TaskId = oldTask.TaskId
	taskPrize.RewardNum = oldTask.RewardNum
	taskPrize.RewardType = oldTask.TaskRewardTypeId
	err = taskPrize.Insert()
	if err != nil {
		log.Errorf("taskPrize.Insert err %s", err)
		result.Msg = "update fail ."
		c.JSON(http.StatusOK, result)
		return
	}
	//发放奖品   金币
	if oldTask.TaskRewardTypeId == 1 {
		err = models.UsersWallet.AddBlance(UserId, oldTask.RewardNum)
		if err != nil {
			log.Errorf("models.UsersWallet.AddBlance err %s", err)
			result.Msg = "add balance fail ."
			c.JSON(http.StatusOK, result)
			return
		}
	}
	result.Ret = 1
	result.Data = task
	c.JSON(http.StatusOK, result)
	return
}

type ListResult struct {
	define.BaseResult
	Data models.UserTaskMongoModel `json:"data"`
}

// @Title 任务列表
// @Description 任务列表
// @Summary 任务列表
// @Accept json
// @Param   token     query    string   true        "token"
// @Param   user_id     path    int   true        "user_id"
// @Success 200 {object} c_task.ListResult
// @router /task/{user_id}/list [get]
func List(c *gin.Context) {
	var result ListResult
	UserId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		result.Msg = "bind params fail ."
		c.JSON(http.StatusOK, result)
		return
	}
	data, err := models.UserTask.Get(UserId)
	if err != nil {
		log.Errorf("models.UserTask.Get err %s", err)
		result.Msg = "get fail ."
		c.JSON(http.StatusOK, result)
		return
	}
	t := time.Now()
	var reData []models.UserTaskModel
	for _, v := range data.List {
		//判断类型//过期类型
		if v.TaskTypeExpireType == define.PermanentTask { //是否过期
			if v.ExpireTime.Unix() <= int64(t.Unix()) {
				//删除任务
				err = models.UserTask.RemoveByTaskId(UserId, v.TaskId)
				if err != nil {
					log.Errorf("remove user(%v) expire task(%v) fail (%s)", UserId, v.TaskId, err)
				}
				continue
			}

		}
		if v.TaskTypeExpireType == define.TodayTask {
			if time.Unix(v.LastUpdate, 0).Format(define.FormatDate) != t.Format(define.FormatDate) {
				//重置今日任务
				task, err := models.UserTask.GetById(v.TaskId)
				if err != nil {
					log.Errorf("models.UserTask.GetById err %s", err)
					result.Msg = "get fail ."
					c.JSON(http.StatusOK, result)
					return
				}
				err = resetTask(UserId, task, t)
				if err != nil {
					log.Errorf("resetTask err %s", err)
					result.Msg = "get fail ."
					c.JSON(http.StatusOK, result)
					return
				}
				reData = append(reData, task)
				continue
			}
			reData = append(reData, v)
			continue
		}
		if v.TaskTypeExpireType == define.WeekTask {
			_, week1 := time.Unix(v.LastUpdate, 0).ISOWeek()
			_, week2 := t.ISOWeek()
			if week1 != week2 {
				//重置今日任务
				task, err := models.UserTask.GetById(v.TaskId)
				if err != nil {

				}
				err = resetTask(UserId, task, t)
				if err != nil {
					log.Errorf("resetTask err %s", err)
					result.Msg = "get fail ."
					c.JSON(http.StatusOK, result)
					return
				}
				reData = append(reData, task)
				continue

			}
		}
		reData = append(reData, v)
	}
	result.Data.List = reData
	result.Data.UserId = UserId
	result.Ret = 1
	c.JSON(http.StatusOK, result)
	return
}

func resetTask(userId int, task models.UserTaskModel, t time.Time) error {

	task.LastUpdate = t.Unix()
	err := models.UserTask.UpdateOneTask(userId, task.TaskId, task)
	return err
}
