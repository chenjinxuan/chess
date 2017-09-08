package handler

import (
	"chess/common/define"
	"chess/common/log"
	"chess/models"
	"chess/srv/srv-task/redis"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"sync"
	"time"
)

type TaskUpsetManager struct {
	TaskAll map[int]models.UserTaskModel
	UserId  chan int
	Mutex   sync.Mutex
}

var TaskUpsetMgr *TaskUpsetManager

func GetTaskUpsetMgr() *TaskUpsetManager {
	if TaskUpsetMgr == nil {
		TaskUpsetMgr = new(TaskUpsetManager)
		err := TaskUpsetMgr.Init()
		if err != nil {
			log.Errorf("init taskhandler fail (%s)", err)
			TaskUpsetMgr = nil
		}
	}
	return TaskUpsetMgr
}

func (m *TaskUpsetManager) Init() (err error) {
	log.Info("init TaskUpset ,manager ...")
	m.UserId = make(chan int, 2)
	err = m.initAllTask()
	return err
}
func (m *TaskUpsetManager) initAllTask() (err error) {
	data, err := models.UserTask.GetAll()
	if err != nil {
		return
	}
	m.TaskAll = make(map[int]models.UserTaskModel)
	for _, v := range data {
		m.TaskAll[v.TaskId] = v
	}
	return
}

//定时更新任务数据
func (m *TaskUpsetManager) LoopGetALLTask() {
	go func() {
		for {
			data, err := models.UserTask.GetAll()
			if err != nil {
				log.Errorf("models.UserTask.GetAll fail (%s)", err)
			}
			m.Mutex.Lock()
			m.TaskAll = nil
			m.TaskAll = make(map[int]models.UserTaskModel)
			for _, v := range data {
				m.TaskAll[v.TaskId] = v
			}
			m.Mutex.Unlock()
			time.Sleep(time.Duration(60) * time.Second)
		}

	}()
}
func (m *TaskUpsetManager) SubLoop() {
	go func() {
		for {
			userIdStr, err := task_redis.Redis.Task.Spop(define.TaskUpsetRedisKey)
			if err != nil {
				if err == redis.ErrNil {
					time.Sleep(time.Duration(10) * time.Second)
					log.Debug("upset QueueTimeout. ")
					continue
				}
			}
			userId, err := strconv.Atoi(userIdStr)
			if userId == 0 {
				continue
			}
			m.UserId <- userId
		}
	}()
}

func (m *TaskUpsetManager) Loop() {
	go func() {
		for {
			select {
			case userId := <-m.UserId:
				func(userId int) {
					m.Mutex.Lock()
					defer m.Mutex.Unlock()
					var taskMap = m.TaskAll
					//一直取出
					//取出所有该用户已经领取的任务id
					allReceive, err := models.TaskPriceReceive.GetAllByUserId(userId)
					if err != nil {
						log.Errorf("models.TaskPriceReceive.GetAllByUserIduser() fail(%s)", userId, err)
						return
					}
					for _, v := range allReceive {
						delete(taskMap, v)
					}
					data, err := models.UserTask.Get(userId)
					if err != nil {
						log.Errorf("models.UserTask.Get user() fail(%s)", userId, err)
						return
					}
					var upsetTaskList models.UserTaskMongoModel
					var isUpset int
					for _, v := range data.List {
						if tsk, ok := taskMap[v.TaskId]; ok {
							delete(taskMap, v.TaskId)
							//对比两个任务的版本号,,如果不同则更新
							if v.Version != tsk.Version {
								//更新
								tsk.AlreadyCompleted = v.AlreadyCompleted
								tsk.LastUpdate = time.Now().Unix()
								upsetTaskList.List = append(upsetTaskList.List, tsk)
								isUpset++
								continue
							} else {
								upsetTaskList.List = append(upsetTaskList.List, v)
								continue
							}
						}
					}
					for _, v := range taskMap {
						if v.ParentId == 0 { //只有一级新任务才会插入
							upsetTaskList.List = append(upsetTaskList.List, v)
							isUpset++
						}
					}
					//更新
					if isUpset >= 1 {
						log.Errorf("you gengxin")
						upsetTaskList.UserId = userId
						err = models.UserTask.Upsert(userId, upsetTaskList)
						if err != nil {
							log.Errorf("models.UserTask.Upsert user() fail(%s)", userId, err)
							return
						}
					}
				}(userId)
			}
		}
	}()
}
