package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type GamblingModel struct {
	RoomId  int       `bson:"room_id" json:"room_id" description:"房间id"`
	TableId string    `bson:"table_id" json:"table_id" description:"桌子id"`
	Max     int       `bson:"max" json:"max" description:"最多玩家数"`
	Start   int64     `bson:"start" json:"start" description:"开始时间"`
	End     int64     `bson:"end" json:"end" description:"结束时间"`
	Cards   []*Card   `bson:"cards" json:"cards" description:"牌"`
	Button  int       `bson:"button" json:"button" description:"庄家"`
        Sb      int       `bson:"sb" json:"sb" description:"小盲注"`
        Bb      int       `bson:"bb" json:"bb" description:"大盲注"`
	SbPos   int       `bson:"sb_pos" json:"sb_pos" description:"小盲注位置"`
	BbPos   int       `bson:"bb_pos" json:"bb_pos" description:"大盲注位置"`
	Pot     []int32   `bson:"pot" json:"pot" description:"奖池"`
	Players []*Player `bson:"players" json:"players" description:"玩家"`
}
type Card struct {
	Suit  int `bson:"suit" json:"suit" description:"花色"`
	Value int `bson:"value" json:"value" description:"牌值"`
}

type Player struct {
	Id             int     `bson:"id" json:"id" description:"玩家id"`
	Nickname       string  `bson:"nickname" json:"nickname" description:"昵称"`
	Avatar         string  `bson:"avatar" json:"avatar" description:"头像"`
	Pos            int     `bson:"pos" json:"pos" description:"位置"`
	Bet            int     `bson:"bet" json:"bet" description:"投注金额"`
	Win            int     `bson:"win" json:"win" description:"赢取金额"`
	FormerChips    int     `bson:"former_chips" json:"former_chips" description:"开始筹码"`
	CurrentChips   int     `bson:"current_chips" json:"current_chips" description:"当前筹码"`
	Action         string  `bson:"action" json:"action" description:"动作"`
	Cards          []*Card `bson:"cards" json:"cards" description:"手牌"`
	HandLevel      int     `bson:"hand_level" json:"hand_level" description:"牌等级"`
	HandFinalValue int     `bson:"hand_final_value" json:"hand_final_value" description:"牌值"`
}
var Gambling = new(GamblingModel)
func (m *GamblingModel) Upsert() error {
	return Mongo.Chess.M(MongoDBStr, MongoColGambling, func(c *mgo.Collection) error {
		query := bson.M{
			"table_id": m.TableId,
			"start":    m.Start,
		}
		_, err := c.Upsert(query, m)
		return err
	})
}
//取出最近一个牌局
func (m *GamblingModel) GetByTableId(tableId string) (info GamblingModel,err error) {
    err = Mongo.Chess.M(MongoDBStr, MongoColGambling, func(c *mgo.Collection) error {
	query := bson.M{
	    "table_id": tableId,
	}
	return c.Find(query).Sort("-end").One(&info)
    })
    return
}