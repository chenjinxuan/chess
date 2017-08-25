package user_init

import (
	"chess/models"
)

//用户任务系统初始化
func TaskInit(user_id int) error {
	//查出所有一级任务
	list, err := models.UserTask.GetInit()
	if err != nil {
		return err
	}
	var userTaskInfo models.UserTaskMongoModel
	userTaskInfo.UserId = user_id
	userTaskInfo.List = list
	err = models.UserTask.Upsert(user_id, userTaskInfo)
	return err
}
