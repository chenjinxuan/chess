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
    "fmt"
)

type TaskChargeManager struct {
    GoodsMap map[int]models.GoodsModel
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
    data, err := models.Goods.List()
    if err != nil {
	return
    }
    m.GoodsMap = make(map[int]models.GoodsModel)
    for _, v := range data {
	m.GoodsMap[v.Id] = v
    }
    return
}

//定时更新任务数据
func (m *TaskChargeManager) LoopGetALLGoods() {
    go func() {
	for {
	    data, err := models.Goods.List()
	    if err != nil {
		log.Errorf("models.UserTask.GetAll fail (%s)", err)
	    }
	    m.Mutex.Lock()
	    m.GoodsMap = nil
	    m.GoodsMap = make(map[int]models.GoodsModel)
	    for _, v := range data {
		m.GoodsMap[v.Id] = v
	    }
	    m.Mutex.Unlock()
	    time.Sleep(time.Duration(60) * time.Second)
	}

    }()
}

type ChargeInfo struct {
    UserId int `json:"user_id"`
    ChargeGoodsId string `json:"charge_goods_id"`
    GoodsId int `json:"goods_id"`
    Price   int `json:"price"`
    BuyTime int64 `json:"buy_time"`
    From    string `json:"from"`
    ChargeFrom string `json:"charge_from"`
} 
func (m *TaskChargeManager) SubLoop() {
    go func() {
	for {
	    data, err := task_redis.Redis.Task.Spop(define.TaskChargeGoodsRedisKey)
	    if err != nil {
		if err == redis.ErrNil {
		    time.Sleep(time.Duration(10) * time.Second)
		    log.Debug("upset QueueTimeout. ")
		    continue
		}
	    }
	    var charge ChargeInfo
	    err=json.Unmarshal([]byte(data),&charge)
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
		    //购买商品,,要给用户的背包加商品
		    //先查出用户是否有这个商品
		    data,err:=models.UserBag.Get(charge.UserId)
		    if err != nil {
			if fmt.Sprint(err) == "not found" {
			   var data models.UserBagMongoModel
			    data.UserId =charge.UserId
			}else {
			    log.Errorf("models.UserBag.Get",err)
			    return
			}
		    }
		    goodsExist := 0 
		    for _,v:=range data.List{
			if v.GoodsId == charge.GoodsId {
			    v.Number++
			    goodsExist=1
			}
		    }
		    if goodsExist == 0 {
			info :=models.UserBagModel{GoodsId:charge.GoodsId,Number:1}
			data.List=append(data.List,info)
		    }
		    //更新
		    err = models.UserBag.Upsert(charge.GoodsId,data)
		    if err != nil {
			log.Errorf("models.UserBag.Upsert",err)
		    }
		}(charge)
	    }
	}
    }()
}
