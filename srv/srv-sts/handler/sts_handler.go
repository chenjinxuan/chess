package handler

import (
	"chess/common/define"
	"chess/common/log"
	"chess/models"
	. "chess/srv/srv-sts/proto"
	"chess/srv/srv-sts/redis"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type StsHandlerManager struct {
	GameInfo chan GameTableInfoArgs
}

var StsMgr *StsHandlerManager

func GetStsHandlerMgr() *StsHandlerManager {
	if StsMgr == nil {
	    StsMgr = new(StsHandlerManager)
		err := StsMgr.Init()
		if err != nil {
			log.Errorf("init stshandler fail (%s)", err)
		    StsMgr = nil
		}
	}
	return StsMgr
}

func (m *StsHandlerManager) Init() (err error) {
	log.Info("init StsHandler ,manager ...")
        m.GameInfo=make(chan GameTableInfoArgs)
	return nil
}


func (m *StsHandlerManager) SubLoop() {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	go func() {
		for {
			res, err := sts_redis.Redis.Sts.Brpop(define.STS_GAME_INFO_REDIS_KEY, 30)
			if err != nil {
				if err == redis.ErrNil {
					log.Debug("QueueTimeout...")

				} else {
					log.Errorf("Get game info channel info fail (%s)", err)
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
			m.GameInfo <- gameInfo
		}
	}()
}

func (m *StsHandlerManager) Loop() {
	go func() {
		for {
			select {
			case _gameInfo := <-m.GameInfo:
				func(gameInfo GameTableInfoArgs) {//玩家统计
				//查出该局所有玩家
				    //比牌型
				   var handLevel int32
				   var  handVaule int32
				    var winner int32
				    for _,v:=range gameInfo.Player {
					if v==nil {
					    continue
					}
					if handLevel==0 {
					    handLevel=v.HandLevel
					    handVaule=v.HandFinalValue
					    winner=v.Id
					}else {
					    if (handLevel==v.HandLevel && handVaule<= v.HandFinalValue) ||(handLevel<v.HandLevel){
						handLevel=v.HandLevel
						handVaule=v.HandFinalValue
						winner=v.Id
					    }
					}
				    }
				    for _,v:=range gameInfo.Player {
					if v==nil {
					    continue
					}
					//取出该玩家信息
					userInfo,err:=models.UserGameSts.Get(int(v.Id))
					if err != nil {
					    log.Errorf("models.UserGameSts.Get",err)
 						continue
					}
					//对比数据
					//胜利数
					if v.Id == winner {
					    userInfo.Win++
					}
					if userInfo.BestWinner<int(v.Win) {
					    userInfo.BestWinner=int(v.Win)
					}
					if (userInfo.HandLevel==int(v.HandLevel) && userInfo.HandFinalValue<int(v.HandFinalValue) ) || (userInfo.HandLevel<int(v.HandLevel)){
					         userInfo.HandLevel=int(v.HandLevel)

						userInfo.HandFinalValue=int(v.HandFinalValue)
						//此时更新牌
						userInfo.Cards=nil
						for _,val:=range gameInfo.Cards {
						    var card models.Card
						    card.Value=int(val.Value)
						    card.Suit=int(val.Suit)
						    userInfo.Cards=append(userInfo.Cards,card)

						}
						for _,val:=range v.Cards {
						    var card models.Card
						    card.Value=int(val.Value)
						    card.Suit=int(val.Suit)
						    userInfo.Cards=append(userInfo.Cards,card)
						}

					}

					//入局数
					if v.Bet>gameInfo.Bb {
					    userInfo.Inbound++
					}
					//摊牌数  计算动作次数,,大于4次就是摊牌
					if len(v.Actions)>=4 {
					    userInfo.Showdown++
					}

					//总局数
					userInfo.TotalGame++

					//更新数据
					err=models.UserGameSts.Upsert(int(v.Id),userInfo)
					if err != nil {
					    log.Errorf("models.UserGameSts.Upsert",err)
					}


				    }


				}(_gameInfo)


			}
		}
	}()
}
