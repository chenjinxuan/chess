package handler

import (
    "chess/common/define"
    "chess/common/log"
    "chess/models"
    "chess/srv/srv-task/redis"
    "github.com/garyburd/redigo/redis"
    . "chess/srv/srv-task/proto"
   // "sync"
    "encoding/json"
    "fmt"
)

type TaskUserBagManager struct {
    UserInfo chan UpdateBagArgs
    //Mutex   sync.Mutex
}
type UserInfo struct {
    UserId int `json:"user_id"`
    GoodsId int `json:"goods_id"`
}
var TaskUserBagMgr *TaskUserBagManager

func GetTaskUserBagMgr() *TaskUserBagManager {
    if TaskUserBagMgr == nil {
	TaskUserBagMgr = new(TaskUserBagManager)
	err := TaskUserBagMgr.Init()
	if err != nil {
	    log.Errorf("init taskhandler fail (%s)", err)
	    TaskUserBagMgr = nil
	}
    }
    return TaskUserBagMgr
}

func (m *TaskUserBagManager) Init() (err error) {
    log.Info("init TaskUpset ,manager ...")
    m.UserInfo = make(chan UpdateBagArgs, 1)
    return
}



func (m *TaskUserBagManager) SubLoop() {
    defer func() {
	if err := recover(); err != nil {
	    panic(err)
	}
    }()
    go func() {
	for {
	    res, err := task_redis.Redis.Task.Brpop(define.TaskUserBagRedisKey,60)
	    if err != nil {
		if err == redis.ErrNil {
		    log.Debug("QueueTimeout...")

		} else {
		    log.Errorf("Get userbag channel info fail (%s)", err)
		}
		continue
	    }
	    log.Debugf("get userbag channel info(%s)", res)
	    key := res[1]
	    var userInfo UpdateBagArgs
	    err=json.Unmarshal([]byte(key),&userInfo)
	    if err != nil {
               log.Errorf("cannot jsondecord ",err)
		continue
	    }
	    log.Debug(userInfo.UserId)
	    m.UserInfo<- userInfo
	}
    }()
}

func (m *TaskUserBagManager) Loop() {
    defer func() {
	if err := recover(); err != nil {
	    panic(err)
	}
    }()
    go func() {
	for {
	    select {
	    case userInfo := <-m.UserInfo:
		func(userInfo UpdateBagArgs) {
		    if userInfo.UserId == 0  {
			return
		    }
		    //查出bag
		    bag, err := models.UserBag.Get(int(userInfo.UserId))
		    if err != nil {
			if fmt.Sprint(err)=="not found" {
			    bag =models.UserBagMongoModel{}
			}
		    }
		    exist :=0
		    var newBag models.UserBagMongoModel
		    newBag.UserId=int(userInfo.UserId)
		    for _, v := range bag.List {
			if v.GoodsId == int(userInfo.GoodsId) { //找到该商品
			    v.Number++
			    exist=1
			}
			newBag.List=append(newBag.List,v)
		    }
		    if exist==0 {
			newBag.List=append(newBag.List,models.UserBagModel{Number:1,GoodsId:int(userInfo.GoodsId)})
		    }
		    //更新
		    err=models.UserBag.Upsert(int(userInfo.UserId),newBag)
		    if err != nil {
			log.Errorf("models.UserBag.Upsert",err)
		    }

		}(userInfo)
	    }
	}
    }()
}
