package models

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/Sirupsen/logrus"
    "chess/common/log"
    "time"
)

type UserTaskModel struct {
    TaskId int `json:"task_id"`
    ParentId int `json:"parent_id"`
    TaskTypeId int `json:"task_type_id"`
    TaskTypeName string `json:"task_type_name"`
    TaskTypeExpireType int `json:"task_type_expire_type"`
    Name string `json:"name"`
    Image string `json:"image"`
    ImageDescribe string `json:"image_describe"`
    TaskRewardTypeId int `json:"task_reward_type_id"`
    RewardNum int `json:"reward_num"`
    TaskRequiredId int `json:"task_required_id"`
    TaskRequiredRoomType int `json:"task_required_room_type"`
    TaskRequiredMatchType int `json:"task_required_match_type"`
    RequiredDescribe string `json:"required_describe"`
    RequiredNum int `json:"required_num"`
    AlreadyCompleted int `json:"already_completed"`
    LastUpdate int64 `json:"last_update"`
    IsWin int `json:"is_win"`
}
var UserTask =  new(UserTaskModel)
func (m *UserTaskModel) GetInit() (list []UserTaskModel,err error) {
    sqlStr := `SELECT a.id,a.parent_id,b.id,b.name,b.expire_type,a.name,a.image,a.image_describe,a.task_reward_type_id,a.reward_num,c.id,c.room_type,c.match_type,a.required_describe,a.required_num,a.is_win
    FROM task AS a
    LEFT JOIN task_type AS b ON a.task_type_id = b.id
    LEFT JOIN task_required AS c ON a.task_required_id = c.id
    WHERE a.parent_id = 0 AND a.status = 1`
    rows ,err := Mysql.Chess.Query(sqlStr)
    if err != nil {
	return
    }
    defer rows.Close()
    for rows.Next()  {
	var t UserTaskModel
	err=rows.Scan(&t.TaskId,&t.ParentId,&t.TaskTypeId,&t.TaskTypeName,&t.TaskTypeExpireType,&t.Name,&t.Image,&t.ImageDescribe,&t.TaskRewardTypeId,&t.RewardNum,&t.TaskRequiredId,&t.TaskRequiredRoomType,&t.TaskRequiredMatchType,&t.RequiredDescribe,&t.RequiredNum,&t.IsWin)
	if err != nil {
	    continue
	}
	t.LastUpdate=time.Now().Unix()
	list=append(list,t)
    }
    return
}

type UserTaskMongoModel struct {
    UserId int `json:"user_id"`
    List []UserTaskModel
}

func (m *UserTaskModel) Get(userId int) (task UserTaskMongoModel, err error) {
    err = Mongo.Chess.M(MongoDBStr, MongoColUserTask, func(c *mgo.Collection) error {
	query := bson.M{
	    "user_id": userId,
	}
	return c.Find(query).One(&task)
    })

    return
}
func (m *UserTaskModel) Upsert(userId int, Task UserTaskMongoModel) error {

    return Mongo.Chess.M(MongoDBStr, MongoColUserTask, func(c *mgo.Collection) error {
	query := bson.M{
	    "user_id": userId,
	}
	changeInfo, err := c.Upsert(query, Task)
	// Debug
	log.Debugf("UserTaskModel.Upsert", logrus.Fields{
	    "User ID":     userId,
	    "Query":       query,
	    "Change Info": changeInfo,
	    "Error":       err,
	})

	return err
    })
}