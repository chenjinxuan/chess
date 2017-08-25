package handler

import(
    "chess/models"
    . "chess/srv/srv-task/proto"
    "chess/common/log"
    "chess/common/define"
    "chess/srv/srv-task/redis"
    "github.com/garyburd/redigo/redis"
    "encoding/json"
    "time"
)


type TaskHandlerManager struct {
     TaskRequired []models.TaskRequiredModel
     TaskType []models.TaskTypeModel
     TaskRewardType []models.TaskRewardTypeModel
     Type int
     GameOver chan GameInfoArgs
     PlayerEvent chan PlayerActionArgs
     TaskHandlerGameOverRedisKey string
     TaskHandlerPlayerEventRedisKey string
}

var TaskMgr *TaskHandlerManager

func GetTaskHandlerMgr() *TaskHandlerManager {
    if TaskMgr == nil {
	TaskMgr = new(TaskHandlerManager)
	err := TaskMgr.Init()
	if err != nil {
	    log.Errorf("init taskhandler fail (%s)",err)
	    TaskMgr = nil
	}
    }
    return TaskMgr
}

func (m *TaskHandlerManager) Init() (err error) {
    log.Info("init TaskHandler ,manager ...")
    m.TaskHandlerGameOverRedisKey = define.TaskLoopHandleGameOverRedisKey
    m.TaskHandlerPlayerEventRedisKey = define.TaskLoopHandlePlayerEventRedisKey
    m.GameOver =make(chan GameInfoArgs, 1)
    m.PlayerEvent = make(chan PlayerActionArgs, 1)
    err = m.initTaskRequired()
    if err != nil {
	log.Errorf("init taskrequired fail .(%s)",err)
	return
    }
    err = m.initTaskRewardType()
    if err != nil {
	log.Errorf("init taskrewardtype fail .(%s)",err)
	return
    }
    err = m.initTaskType()
    if err != nil {
	log.Errorf("init taskType fail .(%s)",err)
	return
    }
   return nil
}

func (m *TaskHandlerManager) initTaskRequired() (err error)  {
     m.TaskRequired ,err = models.TaskRequired.List()
    return
}
func (m *TaskHandlerManager) initTaskType() (err error)  {
    m.TaskType ,err = models.TaskType.List()
    return
}
func (m *TaskHandlerManager) initTaskRewardType() (err error)  {
    m.TaskRewardType ,err = models.TaskRewardType.List()
    return
}

func (m *TaskHandlerManager) SubLoop() {
    defer func() {
	if err := recover(); err != nil {
	    panic(err)
	}
    }()
   go func() {
       for {
	   res,err := task_redis.Redis.Task.Brpop(m.TaskHandlerGameOverRedisKey,30)
	   if err != nil {
	       if err == redis.ErrNil {
		   log.Debug("QueueTimeout...")

	       }else {
		   log.Errorf("Get game over channel info fail (%s)",err)
	       }
	       continue
	   }
	   log.Debugf("get over channel info(%s)",res)
	   key := res[1]
	   gameInfo := GameInfoArgs{}
	   err = json.Unmarshal([]byte(key),&gameInfo)
	   if err != nil {
	       log.Errorf("game info  could not be marshaled")
	       continue
	   }
	   m.GameOver<-gameInfo
       }
   }()
   go func() {
       for {
	   res,err := task_redis.Redis.Task.Brpop(m.TaskHandlerPlayerEventRedisKey,30)
	   if err != nil {
	       if err == redis.ErrNil {
		   log.Debug("QueueTimeout...")
	       }else {
		   log.Errorf("Get Player evernt channel info fail (%s)",err)
	       }
	       continue
	   }
	   log.Debugf("get Player evernt channel info(%s)",res)
	   key := res[1]
	   ActionInfo := PlayerActionArgs{}
	   err = json.Unmarshal([]byte(key),&ActionInfo)
	   if err != nil {
	       log.Errorf("game info  could not be marshaled")
	       continue
	   }
	   m.PlayerEvent<-ActionInfo
       }

   }()
}

func (m *TaskHandlerManager) Loop(){
     go func() {
	 for  {
	     select {
	     case _gameInfo := <-m.GameOver:
		 func(gameInfo GameInfoArgs) {
		     //判断现有的任务要求 类型,中奖类型
                      overTime :=time.Unix(int64(gameInfo.Time),0)
		     //循环查出所有该局用户的任务信息并更新信息
		     for _,v:=range gameInfo.Players {
                         taskList,err:=models.UserTask.Get(int(v.Id))
			 if err != nil {
			     log.Errorf("get user(%v) tasklist fail (%s)",v.Id,err)
			     continue
			 }
			 for _,taskInfo:=range taskList.List  {
			     if taskInfo.TaskRequiredRoomType != 0 { //是否要求房间场次类型
				 if taskInfo.TaskRequiredRoomType != int(gameInfo.RoomType) {
				     continue
				 }

			     }
			     if taskInfo.TaskRequiredMatchType != 0  { //是否要求赛事类型
				 if taskInfo.TaskRequiredMatchType != int(gameInfo.MatchType) {
				     continue
				 }
			     }

			     if taskInfo.IsWin != 0 { //是否需要胜利
				 if v.Id != gameInfo.Winner {
				     continue
				 }
			     }
			     if taskInfo.TaskRequiredHandLevel != 0 { //是否需要手牌等级
				 if int(v.HandLevel) != taskInfo.TaskRequiredHandLevel {
				     continue
				 }
			     }
			     if taskInfo.TaskTypeExpireType == define.PermanentTask { //是否过期
				 if  taskInfo.ExpireTime.Unix() <= int64(gameInfo.Time) {
				     //删除任务
				     err = models.UserTask.RemoveByTaskId(int(v.Id),taskInfo.TaskId)
				     if err != nil {
					 log.Errorf("remove user(%v) expire task(%v) fail (%s)",v.Id,taskInfo.TaskId,err)
				     }
				     continue
				 }

			     }
			     if taskInfo.TaskTypeExpireType == define.TodayTask {
				 if time.Unix(taskInfo.LastUpdate,0).Format(define.FormatDate) != overTime.Format(define.FormatDate){
				     continue
				 }
			     }
			     if taskInfo.TaskTypeExpireType == define.WeekTask {
				 _,week1 := time.Unix(taskInfo.LastUpdate,0).ISOWeek()
				 _,week2 := overTime.ISOWeek()
				 if week1 != week2{
				     continue
				 }
			     }
			     //排除不符合任务要求后,,,更新任务
			     taskInfo.LastUpdate = int64(gameInfo.Time)
			     taskInfo.AlreadyCompleted = taskInfo.AlreadyCompleted + 1

			     err = models.UserTask.UpdateOneTask(taskList.UserId,taskInfo.TaskId,taskInfo)
			     if err != nil {
				 log.Errorf("update user(%v) taskid(%v) fail  (%s)",taskList.UserId,taskInfo.TaskId,err)
			     }
			 }
		     }

		 }(_gameInfo)

	     case _playAction := <-m.PlayerEvent:
		 func(playAction PlayerActionArgs){
		     actionTime :=time.Unix(int64(playAction.Time),0)
		     //查出该用户的任务信息
		     taskList,err:=models.UserTask.Get(int(playAction.Id))
		     if err != nil {
			 log.Errorf("get user(%v) tasklist fail (%s)",playAction.Id,err)
		     }
		     for _,taskInfo:=range taskList.List {
			 if taskInfo.TaskRequiredPlayerAction != 0 {
			     if taskInfo.TaskRequiredPlayerAction != int(playAction.Type)  {
				 continue
			     }
			 }
			 if taskInfo.TaskRequiredRoomType != 0 { //是否要求房间场次类型
			     if taskInfo.TaskRequiredRoomType != int(playAction.RoomType) {
				 continue
			     }

			 }
			 if taskInfo.TaskRequiredMatchType != 0  { //是否要求赛事类型
			     if taskInfo.TaskRequiredMatchType != int(playAction.MatchType) {
				 continue
			     }
			 }
			 if taskInfo.TaskTypeExpireType == define.PermanentTask { //是否过期
			     if  taskInfo.ExpireTime.Unix() <= int64(playAction.Time) {
				 //删除任务
				 err = models.UserTask.RemoveByTaskId(int(playAction.Id),taskInfo.TaskId)
				 if err != nil {
				     log.Errorf("remove user(%v) expire task(%v) fail (%s)",playAction.Id,taskInfo.TaskId,err)
				 }
				continue
			     }

			 }
			 if taskInfo.TaskTypeExpireType == define.TodayTask {
			     if time.Unix(taskInfo.LastUpdate,0).Format(define.FormatDate) != actionTime.Format(define.FormatDate){
				 continue
			     }
			 }
			 if taskInfo.TaskTypeExpireType == define.WeekTask {
			     _,week1 := time.Unix(taskInfo.LastUpdate,0).ISOWeek()
			     _,week2 := actionTime.ISOWeek()
			     if week1 != week2{
				 continue
			     }
			 }
			 //排除不符合任务要求后,,,更新任务
			 taskInfo.LastUpdate = int64(playAction.Time)
			 taskInfo.AlreadyCompleted = taskInfo.AlreadyCompleted + 1

			 err = models.UserTask.UpdateOneTask(taskList.UserId,taskInfo.TaskId,taskInfo)
			 if err != nil {
			     log.Errorf("update user(%v) taskid(%v) fail  (%s)",taskList.UserId,taskInfo.TaskId,err)
			 }
		     }

		 }(_playAction)


	     }
	     log.Error("11111111111111")
	 }
     }()
}

