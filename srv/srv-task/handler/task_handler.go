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

type TaskHandlerManager struct {
	TaskRequired                   []models.TaskRequiredModel
	TaskType                       []models.TaskTypeModel
	TaskRewardType                 []models.TaskRewardTypeModel
	RoomType                       map[int]models.RoomsModel
	Type                           int
	GameOver                       chan GameTableInfoArgs
	PlayerEvent                    chan PlayerActionArgs
	TaskHandlerGameOverRedisKey    string
	TaskHandlerPlayerEventRedisKey string
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
	m.GameOver = make(chan GameTableInfoArgs, 2)
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

func (m *TaskHandlerManager) initTaskRequired() (err error) {
	m.TaskRequired, err = models.TaskRequired.List()
	return
}
func (m *TaskHandlerManager) initTaskType() (err error) {
	m.TaskType, err = models.TaskType.List()
	return
}
func (m *TaskHandlerManager) initTaskRewardType() (err error) {
	m.TaskRewardType, err = models.TaskRewardType.List()
	return
}

func (m *TaskHandlerManager) initRoomType() (err error) {
	m.RoomType = make(map[int]models.RoomsModel)
	data, err := models.Rooms.GetAll()
	if err != nil {
		return
	}
	for _, v := range data {
		m.RoomType[v.Id] = v
	}
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
			res, err := task_redis.Redis.Task.Brpop(m.TaskHandlerGameOverRedisKey, 30)
			if err != nil {
				if err == redis.ErrNil {
					log.Debug("QueueTimeout...")

				} else {
					log.Errorf("Get game over channel info fail (%s)", err)
				}
				continue
			}
			log.Debugf("get over channel info(%s)", res)
			key := res[1]
			gameInfo := GameTableInfoArgs{}
			err = json.Unmarshal([]byte(key), &gameInfo)
			if err != nil {
				log.Errorf("game info  could not be marshaled")
				continue
			}
			m.GameOver <- gameInfo
		}
	}()
	go func() {
		for {
			res, err := task_redis.Redis.Task.Brpop(m.TaskHandlerPlayerEventRedisKey, 30)
			if err != nil {
				if err == redis.ErrNil {
					log.Debug("QueueTimeout...")
				} else {
					log.Errorf("Get Player evernt channel info fail (%s)", err)
				}
				continue
			}
			log.Debugf("get Player evernt channel info(%s)", res)
			key := res[1]
			ActionInfo := PlayerActionArgs{}
			err = json.Unmarshal([]byte(key), &ActionInfo)
			if err != nil {
				log.Errorf("game info  could not be marshaled")
				continue
			}
			m.PlayerEvent <- ActionInfo
		}

	}()
}

func (m *TaskHandlerManager) Loop() {//可能会出现,,两个动作同时发生的情况
	go func() {
		for {
			select {
			case _gameInfo := <-m.GameOver:
				func(gameInfo GameTableInfoArgs) {
					//判断现有的任务要求 类型,中奖类型
					overTime := time.Unix(int64(gameInfo.End), 0)
				        //判断胜利者
				    //比牌型
				    var handLevel int32
				    var  handVaule int32
				    var winner []int32
				    for _,v:=range gameInfo.Players {
					if v==nil {
					    continue
					}
					if handLevel==0 {
					    handLevel=v.HandLevel
					    handVaule=v.HandFinalValue
					    winner=append(winner,v.Id)
					}else {
					    if (handLevel==v.HandLevel && handVaule<= v.HandFinalValue) ||(handLevel<v.HandLevel){
						handLevel=v.HandLevel
						handVaule=v.HandFinalValue
						if  handLevel==v.HandLevel && handVaule == v.HandFinalValue{
						    winner=append(winner,v.Id)
						}else {
						    winner = nil
						    winner=append(winner,v.Id)
						}
					    }
					}
				    }
					//循环查出所有该局用户的任务信息并更新信息
					for _, v := range gameInfo.Players{
						taskList, err := models.UserTask.Get(int(v.Id))
						if err != nil {
							log.Errorf("get user(%v) tasklist fail (%s)", v.Id, err)
							continue
						}
						var newTaskList models.UserTaskMongoModel
						newTaskList.UserId = int(v.Id)
						for _, taskInfo := range taskList.List {

							if taskInfo.TaskRequiredRoomType != 0 { //是否要求房间场次类型
								if taskInfo.TaskRequiredRoomType != m.RoomType[int(gameInfo.RoomId)].RoomsTypeId {
									newTaskList.List = append(newTaskList.List, taskInfo)
									continue
								}

							}
							//if taskInfo.TaskRequiredMatchType != 0 { //是否要求赛事类型
							//	if taskInfo.TaskRequiredMatchType != int(gameInfo.MatchType) {v
							//		newTaskList.List = append(newTaskList.List, taskInfo)
							//		continue
							//	}
							//}

							if taskInfo.IsWin != define.RequiredCommon { //是否需要胜利
							    w:=0
							    for _,winnerId:=range winner {
								if v.Id==winnerId {
								    w=1
								}
							    }
								if w==1 && taskInfo.IsWin == define.RequiredWin {
									newTaskList.List = append(newTaskList.List, taskInfo)
									continue
								}
							} else if taskInfo.IsWin == define.RequiredTotalBalance {
								newTaskList.List = append(newTaskList.List, taskInfo)
								continue
							}

							if taskInfo.TaskRequiredHandLevel != 0 { //是否需要手牌等级
								if int(v.HandLevel) != taskInfo.TaskRequiredHandLevel {
									newTaskList.List = append(newTaskList.List, taskInfo)
									continue
								}
							}
							if taskInfo.TaskTypeExpireType == define.PermanentTask { //是否过期
								if taskInfo.ExpireTime.Unix() <= int64(gameInfo.End) {
									//删除任务
									err = models.UserTask.RemoveByTaskId(int(v.Id), taskInfo.TaskId)
									if err != nil {
										log.Errorf("remove user(%v) expire task(%v) fail (%s)", v.Id, taskInfo.TaskId, err)
									}
									continue
								}

							}
							if taskInfo.TaskTypeExpireType == define.TodayTask {
								if time.Unix(taskInfo.LastUpdate, 0).Format(define.FormatDate) != overTime.Format(define.FormatDate) {
									newTaskList.List = append(newTaskList.List, taskInfo)
									continue
								}
							}
							if taskInfo.TaskTypeExpireType == define.WeekTask {
								_, week1 := time.Unix(taskInfo.LastUpdate, 0).ISOWeek()
								_, week2 := overTime.ISOWeek()
								if week1 != week2 {
									newTaskList.List = append(newTaskList.List, taskInfo)
									continue
								}
							}
							//排除不符合任务要求后,,,更新任务
							taskInfo.LastUpdate = int64(gameInfo.End)
							taskInfo.AlreadyCompleted = taskInfo.AlreadyCompleted + 1

							newTaskList.List = append(newTaskList.List, taskInfo)
							//err = models.UserTask.UpdateOneTask(taskList.UserId, taskInfo.TaskId, taskInfo)
							//if err != nil {
							//	log.Errorf("update user(%v) taskid(%v) fail  (%s)", taskList.UserId, taskInfo.TaskId, err)
							//}
						}
						err = models.UserTask.Upsert(int(v.Id), newTaskList)
						if err != nil {
							log.Errorf("update user(%v)  fail  (%s)", taskList.UserId, err)
						}
					}

				}(_gameInfo)

			case _playAction := <-m.PlayerEvent:
				func(playAction PlayerActionArgs) {
					actionTime := time.Unix(int64(playAction.Time), 0)
					//查出该用户的任务信息
					taskList, err := models.UserTask.Get(int(playAction.Id))
					if err != nil {
						log.Errorf("get user(%v) tasklist fail (%s)", playAction.Id, err)
					}
					for _, taskInfo := range taskList.List {
						//不是动作类型的任务直接跳过
						if taskInfo.TaskRequiredPlayerAction == 0 {
							continue
						}

						if taskInfo.TaskRequiredPlayerAction != int(playAction.Type) {
							continue
						}

						if taskInfo.TaskRequiredRoomType != 0 { //是否要求房间场次类型
							if taskInfo.TaskRequiredRoomType != m.RoomType[int(playAction.RoomId)].RoomsTypeId {
								continue
							}

						}
						if taskInfo.TaskRequiredMatchType != 0 { //是否要求赛事类型
							if taskInfo.TaskRequiredMatchType != int(playAction.MatchType) {
								continue
							}
						}
						if taskInfo.TaskTypeExpireType == define.PermanentTask { //是否过期
							if taskInfo.ExpireTime.Unix() <= int64(playAction.Time) {
								//删除任务
								err = models.UserTask.RemoveByTaskId(int(playAction.Id), taskInfo.TaskId)
								if err != nil {
									log.Errorf("remove user(%v) expire task(%v) fail (%s)", playAction.Id, taskInfo.TaskId, err)
								}
								continue
							}

						}
						if taskInfo.TaskTypeExpireType == define.TodayTask {
							if time.Unix(taskInfo.LastUpdate, 0).Format(define.FormatDate) != actionTime.Format(define.FormatDate) {
								continue
							}
						}
						if taskInfo.TaskTypeExpireType == define.WeekTask {
							_, week1 := time.Unix(taskInfo.LastUpdate, 0).ISOWeek()
							_, week2 := actionTime.ISOWeek()
							if week1 != week2 {
								continue
							}
						}
						//排除不符合任务要求后,,,更新任务
						taskInfo.LastUpdate = int64(playAction.Time)
						taskInfo.AlreadyCompleted = taskInfo.AlreadyCompleted + 1

						err = models.UserTask.UpdateOneTask(taskList.UserId, taskInfo.TaskId, taskInfo)
						if err != nil {
							log.Errorf("update user(%v) taskid(%v) fail  (%s)", taskList.UserId, taskInfo.TaskId, err)
						}
					}

				}(_playAction)

			}
		}
	}()
}
