package models



type CheckinDaysRewardModel  struct {
    Id int `json:"id" description:"签到奖励的id"`
    Days int `json:"days" description:"签到的天数"`
    Type int `json:"type" description:"签到奖励的类型"`
    Number int `json:"number" description:"签到奖励的金额"`
    Image string `json:"image" description:"签到奖励的奖品图片"`
    ImageDescribe string `json:"image_describe" description:"描述"`
    Status int `json:"status"`
}

const (
    CheckinDaysRewardTypeGold = 1
)
var CheckinDaysReward = new(CheckinDaysRewardModel)

func (m *CheckinDaysRewardModel) GetAll() (list []CheckinDaysRewardModel,err error) {
    sqlStr :=`SELECT id,days,type,number,image,image_describe FROM checkin_days_reward WHERE status = 1`
    rows,err:=Mysql.Chess.Query(sqlStr)
    if err != nil {
	return
    }
    defer rows.Close()
    for rows.Next() {
	var c CheckinDaysRewardModel
	err = rows.Scan(&c.Id,&c.Days,&c.Type,&c.Number,&c.Image,&c.ImageDescribe)
	if err != nil {
	    continue
	}
	list = append(list,c)
    }
    return
}

func (m *CheckinDaysRewardModel) Get(days int) (c CheckinDaysRewardModel,err error) {
    sqlStr :=`SELECT id,days,type,number,image,image_describe FROM checkin_days_reward WHERE days = ? AND status = 1`
   err=Mysql.Chess.QueryRow(sqlStr,days).Scan(&c.Id,&c.Days,&c.Type,&c.Number,&c.Image,&c.ImageDescribe)
    return
}