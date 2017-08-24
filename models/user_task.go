package models

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/Sirupsen/logrus"
    "chess/common/log"
    "time"
    "fmt"
)

type UserTaskModel struct {
    TaskId int `bson:"task_id" json:"task_id" description:"任务id"`
    ParentId int `bson:"parent_id" json:"parent_id"  description:"任务父级id"`
    TaskTypeId int `bson:"task_type_id" json:"task_type_id" description:"任务类型id"`
    TaskTypeName string `bson:"task_type_name" json:"task_type_name" description:"任务类型名字"`
    TaskTypeExpireType int `bson:"task_type_expire_type" json:"task_type_expire_type" description:"任务过期类型"`
    Name string `bson:"name" json:"name" description:"任务名字"`
    Image string `bson:"image" json:"image" description:"任务图片"`
    ImageDescribe string `bson:"image_describe" json:"image_describe" description:"任务图片描述"`
    TaskRewardTypeId int `bson:"task_reward_type_id" json:"task_reward_type_id" description:"任务奖品类型id 1金币"`
    RewardNum int `bson:"reward_num" json:"reward_num" description:"奖品数量"`
    TaskRequiredId int `bson:"task_required_id" json:"task_required_id" description:"任务要求id"`
    TaskRequiredRoomType int `bson:"task_required_room_type" json:"task_required_room_type" description:"任务要求房间类型"`
    TaskRequiredMatchType int `bson:"task_required_match_type" json:"task_required_match_type" description:"任务要求赛事类型"`
    RequiredDescribe string `bson:"required_describe" json:"required_describe" description:"任务要求描述"`
    RequiredNum int `bson:"required_num" json:"required_num" description:"任务要求数量"`
    AlreadyCompleted int `bson:"already_completed" json:"already_completed" description:"任务已经完成数"`
    LastUpdate int64 `bson:"last_update" json:"last_update" description:"上次更新时间时间戳"`
    IsWin int `bson:"is_win" json:"is_win" description:"是非要赢才算"`
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

func (m *UserTaskModel) GetByParentId(taskId int) (t UserTaskModel ,err error) {
    sqlStr :=`SELECT a.id,a.parent_id,b.id,b.name,b.expire_type,a.name,a.image,a.image_describe,a.task_reward_type_id,a.reward_num,c.id,c.room_type,c.match_type,a.required_describe,a.required_num,a.is_win
    FROM task AS a
    LEFT JOIN task_type AS b ON a.task_type_id = b.id
    LEFT JOIN task_required AS c ON a.task_required_id = c.id
    WHERE a.parent_id = ? AND a.status = 1`
    err=Mysql.Chess.QueryRow(sqlStr,taskId).Scan(&t.TaskId,&t.ParentId,&t.TaskTypeId,&t.TaskTypeName,&t.TaskTypeExpireType,&t.Name,&t.Image,&t.ImageDescribe,&t.TaskRewardTypeId,&t.RewardNum,&t.TaskRequiredId,&t.TaskRequiredRoomType,&t.TaskRequiredMatchType,&t.RequiredDescribe,&t.RequiredNum,&t.IsWin)
    return
}

type UserTaskMongoModel struct {
    UserId int `bson:"user_id"`
    List []UserTaskModel `bson:"list" json:"list" `
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

func (m *UserTaskModel) GetOne(userId,taskId int) (*UserTaskModel,error) {
    var task = new(UserTaskModel)
    err := Mongo.Chess.M(MongoDBStr, MongoColUserTask, func(c *mgo.Collection) error {
	query := bson.M{
	    "user_id": userId,
	    "list.task_id":taskId,
	}
	return c.Find(query).One(&task)
    })
    fmt.Println(task)
    return task ,err
}

func (m *UserTaskModel) UpdateOneTask(userId ,taskId int,task UserTaskModel) error{
    return Mongo.Chess.M(MongoDBStr, MongoColUserTask, func(c *mgo.Collection) error {
        query:=bson.M{
	    "user_id": userId,
	    "list.task_id":taskId,
	}
	update :=bson.M{
	   "list.$.task_id": task.TaskId,
	    "list.$.parent_id": task.ParentId,
	    "list.$.task_type_id": task.TaskTypeId,
	    "list.$.task_type_name": task.TaskTypeName,
	    "list.$.task_type_expire_type": task.TaskTypeExpireType,
	    "list.$.name": task.Name,
	    "list.$.image": task.Image,
	    "list.$.image_describe": task.ImageDescribe,
	    "list.$.task_reward_type_id": task.TaskRewardTypeId,
	    "list.$.reward_num": task.RewardNum,
	    "list.$.task_required_id": task.TaskRequiredId,
	    "list.$.task_required_room_type": task.TaskRequiredRoomType,
	    "list.$.task_required_match_type": task.TaskRequiredMatchType,
	    "list.$.required_describe": task.RequiredDescribe,
	    "list.$.required_num": task.RequiredNum,
	    "list.$.already_completed": task.AlreadyCompleted,
	    "list.$.last_update": task.LastUpdate,
	    "list.$.is_win": task.IsWin,

	}
	err := c.Update(query, bson.M{"$set":update})

	return err
    })
}

func (m *UserTaskModel) RemoveByTaskId(userId int ,task UserTaskModel ) error {
    return Mongo.Chess.M(MongoDBStr, MongoColUserTask, func(c *mgo.Collection) error {
	query := bson.M{
	    "user_id": userId,
	}
	update :=bson.M{
	    "list":bson.M{
		"task_id": task.TaskId,
	    },

	}
	err := c.Update(query, bson.M{"$pull":update})

	return err
    })
}