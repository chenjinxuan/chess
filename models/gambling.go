package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type GamblingModel struct {
	RoomId  int       `bson:"room_id"`
	TableId string    `bson:"table_id"`
	Max     int       `bson:"max"`
	Start   int64     `bson:"start"`
	End     int64     `bson:"end"`
	Cards   []*Card   `bson:"cards"`
	Button  int       `bson:"button"`
	SbPos   int       `bson:"sb_pos"`
	BbPos   int       `bson:"bb_pos"`
	Pot     []int32   `bson:"pot"`
	Players []*Player `bson:"players"`
}

type Card struct {
	Suit  int `bson:"suit"`
	Value int `bson:"value"`
}

type Player struct {
	Id             int     `bson:"id"`
	Nickname       string  `bson:"nickname"`
	Avatar         string  `bson:"avatar"`
	Pos            int     `bson:"pos"`
	Bet            int     `bson:"bet"`
	Win            int     `bson:"win"`
	FormerChips    int     `bson:"former_chips"`
	CurrentChips   int     `bson:"current_chips"`
	Action         string  `bson:"action"`
	Cards          []*Card `bson:"cards"`
	HandLevel      int     `bson:"hand_level"`
	HandFinalValue int     `bson:"hand_final_value"`
}

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
