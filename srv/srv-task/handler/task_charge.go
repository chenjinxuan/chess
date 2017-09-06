package handler

import (
    "chess/common/define"
    "chess/common/log"
    "chess/models"
    "chess/srv/srv-task/redis"
    "github.com/garyburd/redigo/redis"
    "sync"
    "time"
    "encoding/json"
)

type TaskChargeManager struct {
    ChargeGoodsMap map[string]models.ChargeGoodsModel
    Charge  chan ChargeInfo
    Mutex   sync.Mutex
}

var TaskChargeMgr *TaskChargeManager

func GetTaskChargeMgr() *TaskChargeManager {
    if TaskChargeMgr == nil {
	TaskChargeMgr = new(TaskChargeManager)
	err := TaskChargeMgr.Init()
	if err != nil {
	    log.Errorf("init taskhandler fail (%s)", err)
	    TaskChargeMgr = nil
	}
    }
    return TaskChargeMgr
}

func (m *TaskChargeManager) Init() (err error) {
    log.Info("init TaskUpset ,manager ...")
    m.Charge = make(chan ChargeInfo, 2)
    err=m.initAllGoods()
    return err
}
func (m *TaskChargeManager) initAllGoods() (err error) {
    data, err := models.ChargeGoods.List()
    if err != nil {
	return
    }
    m.ChargeGoodsMap = make(map[string]models.ChargeGoodsModel)
    for _, v := range data {
	m.ChargeGoodsMap[v.ChargeGoodsId] = v
    }
    return
}

//定时更新任务数据
func (m *TaskChargeManager) LoopGetALLGoods() {
    go func() {
	for {
	    data, err := models.ChargeGoods.List()
	    if err != nil {
		log.Errorf("models.UserTask.GetAll fail (%s)", err)
	    }
	    m.Mutex.Lock()
	    m.ChargeGoodsMap = nil
	    m.ChargeGoodsMap = make(map[string]models.ChargeGoodsModel)
	    for _, v := range data {
		m.ChargeGoodsMap[v.ChargeGoodsId] = v
	    }
	    m.Mutex.Unlock()
	    time.Sleep(time.Duration(60) * time.Second)
	}

    }()
}

type ChargeInfo struct {
    UserId int `json:"user_id"`
    ChargeGoodsId string `json:"charge_goods_id"`
    Price   int `json:"price"`
    BuyTime int64 `json:"buy_time"`
    From    string `json:"from"`
    ChargeFrom string `json:"charge_from"`
} 
func (m *TaskChargeManager) SubLoop() {
    go func() {
	for {
	    res, err := task_redis.Redis.Task.Brpop(define.TaskChargeGoodsRedisKey,60)
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
	    var charge ChargeInfo
	    err=json.Unmarshal([]byte(key),&charge)
	    if err != nil {

	    }

	    m.Charge <- charge
	}
    }()
}

func (m *TaskChargeManager) Loop() {
    go func() {
	for {
	    select {
	    case charge := <-m.Charge:
		func(charge ChargeInfo) {
		 //处理支付项目发过来的消息

		}(charge)
	    }
	}
    }()
}
