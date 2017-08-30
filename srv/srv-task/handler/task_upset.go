package handler


import (
    "chess/common/define"
    "chess/common/log"
    "chess/models"
    . "chess/srv/srv-task/proto"
    "chess/srv/srv-task/redis"
    "encoding/json"
    "github.com/garyburd/redigo/redis"
    "time"
)

type TaskUpsetManager struct {
    TaskAll
}

var TaskMgr *TaskHandlerManager

func GetTaskHandlerMgr() *TaskHandlerManager {
    if TaskMgr == nil {
	TaskMgr = new(TaskHandlerManager)
	err := TaskMgr.Init()
	if err != nil {
	    log.Errorf("init taskhandler fail (%s)", err)
	    TaskMgr = nil
	}
    }
    return TaskMgr
}

func (m *TaskHandlerManager) Init() (err error) {
    log.Info("init TaskHandler ,manager ...")
    m.TaskHandlerGameOverRedisKey = define.TaskLoopHandleGameOverRedisKey
    m.TaskHandlerPlayerEventRedisKey = define.TaskLoopHandlePlayerEventRedisKey
    m.GameOver = make(chan GameInfoArgs, 2)
    m.PlayerEvent = make(chan PlayerActionArgs, 2)
    err = m.initTaskRequired()
    if err != nil {
	log.Errorf("init taskrequired fail .(%s)", err)
	return
    }
    err = m.initTaskRewardType()
    if err != nil {
	log.Errorf("init taskrewardtype fail .(%s)", err)
	return
    }
    err = m.initTaskType()
    if err != nil {
	log.Errorf("init taskType fail .(%s)", err)
	return
    }
    err = m.initRoomType()
    if err != nil {
	log.Errorf("init roomType fail .(%s)", err)
	return
    }
    return nil
}

