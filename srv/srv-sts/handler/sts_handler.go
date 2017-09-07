package handler

import (
	"chess/common/define"
	"chess/common/log"
	"chess/models"
	. "chess/srv/srv-sts/proto"
	"chess/srv/srv-sts/redis"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
    	"fmt"
    "sync"
    "time"
    "sort"
)

type StsHandlerManager struct {
	GameInfo chan GameTableInfoArgs
    	GradeExperienceList []models.GradeExperienceModel
        Mutex   sync.Mutex
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
        m.GameInfo=make(chan GameTableInfoArgs,1)
    	err=m.initGrade()
	return err
}

func (m *StsHandlerManager) initGrade() (err error ) {
    m.GradeExperienceList, err = models.GradeExperience.GetAll()
    var gradeSort []models.GradeExperienceModel
    gradeSort, err = models.GradeExperience.GetAll()
    if err != nil {
	log.Errorf("models.UserTask.GetAll fail (%s)", err)
    }
    sort.Sort(GradeList(gradeSort))
    m.GradeExperienceList = gradeSort
    return
}
    //定时更新任务数据
func (m *StsHandlerManager) LoopGetALLGrade() {
    go func() {
	for {
	    log.Debug("grade枷锁")
	    m.Mutex.Lock()
	    var err error
	    var gradeSort []models.GradeExperienceModel
	    gradeSort, err = models.GradeExperience.GetAll()
	    if err != nil {
		log.Errorf("models.UserTask.GetAll fail (%s)", err)
	    }
	    sort.Sort(GradeList(gradeSort))
	    m.GradeExperienceList=gradeSort
	    m.Mutex.Unlock()
	    log.Debug("grade解锁")
	    time.Sleep(time.Duration(60) * time.Second)
	}

    }()
}
type GradeList []models.GradeExperienceModel

func (a GradeList) Len() int {
    return len(a)
}
func (a GradeList) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}
func (a GradeList) Less(i, j int) bool {

    return a[j].Experience < a[i].Experience
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
			if err != nil || res == nil {
				if err == redis.ErrNil {
					log.Debug("QueueTimeout...")

				} else {
					log.Errorf("Get game info channel info fail (%s)", err)
				}
				continue
			}

			log.Debugf("get over channel info(%s)", res)
			val := res[1]
			gameInfo := GameTableInfoArgs{}
			err = json.Unmarshal([]byte(val), &gameInfo)
			if err != nil {
				log.Errorf("game info  could not be marshaled")
				continue
			}
			m.GameInfo <- gameInfo
		    log.Debug("sts channel input end..")
		}
	}()
}

func (m *StsHandlerManager) Loop() {
	go func() {
		for {
		    log.Debug("sts loop start ")
			select {
			case _gameInfo := <-m.GameInfo:
			    log.Debug("sts channel out start..")
				func(gameInfo GameTableInfoArgs) {//玩家统计
				//查出该局所有玩家
				    //比牌型
				    //比牌型
				    var handLevel int32
				    var  handVaule int32
				    var winner []int32
				    for _,v:=range gameInfo.Players {
					if v.Id==0 {
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
				    m.Mutex.Lock()
				    defer m.Mutex.Unlock()
				    timeExpenditure:=gameInfo.End-gameInfo.Start
				    log.Debug("到循环玩家这里")
				    for _,v:=range gameInfo.Players {
					if v.Id == 0 {
					    continue
					}
					log.Debugf("玩家id(%v)",v.Id)
					//取出该玩家信息
					userInfo,err:=models.UserGameSts.Get(int(v.Id))
					if err != nil {
					    if fmt.Sprint(err)=="not found" {
						userInfo=models.UserGameStsModel{UserId:int(v.Id)}
					    }else {
						log.Errorf("models.UserGameSts.Get",err)
						continue
					    }

					}

					//对比数据
					//胜利数
					var experience int32
					for _,winnerId:=range winner {
					    //计算本局所得经验秒数乘以胜利或失败的
					    if winnerId == v.Id{
						userInfo.Win++
						experience=timeExpenditure*2
					    }else {
						experience=timeExpenditure*1
					    }
					}
					//加经验
					userInfo.Experience=userInfo.Experience+int(experience)
					//判断等级
					for k,gradeVal:=range m.GradeExperienceList {
					    if gradeVal.Experience<=userInfo.Experience {
						userInfo.Grade = m.GradeExperienceList[k].Grade
						userInfo.GradeDescribe = m.GradeExperienceList[k].GradeDescribe
						next_Experience:=userInfo.NextExperience
						if k!=0 {
						    next_Experience=m.GradeExperienceList[k-1].Experience
						}else {
						    next_Experience=userInfo.Experience
						}
						userInfo.NextExperience = next_Experience
						break
					    }
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
					num:=0
					for _,z:=range v.Actions{
					    if z.Action!="" {
						num++
					    }
					}
					if num>=4 {
					    userInfo.Showdown++
					}

					//总局数
					userInfo.TotalGame++

					//更新数据
					err=models.UserGameSts.Upsert(int(v.Id),userInfo)
					if err != nil {
					    log.Errorf("models.UserGameSts.Upsert",err)
					}
					log.Debug("到这里了")
				    }

				    log.Debug("sts channel out end..")
				}(_gameInfo)
			}
		}
	}()
}
