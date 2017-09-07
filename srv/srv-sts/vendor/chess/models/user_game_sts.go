package models

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/Sirupsen/logrus"
    "chess/common/log"
)

type UserGameStsModel struct {
    UserId int `bson:"user_id" json:"user_id"`
    Win      int        `bson:"win" json:"win" description:"胜利数"`
    TotalGame    int             `jbson:"total_game" son:"total_game" description:"总局数"`
    Cards        []Card         `bson:"cards" json:"cards" description:"最大牌"`
    HandLevel      int     `bson:"hand_level" json:"hand_level" description:"牌等级"`
    HandFinalValue int     `bson:"hand_final_value" json:"hand_final_value" description:"牌值"`
    BestWinner   int             `bson:"bast_winner" json:"best_winner" description:"最大赢取筹码"`
    Showdown int         `bson:"showdown" json:"showdown" description:"摊牌数"`
    Inbound  int         `bson:"inbound" json:"inbound" description:"入局数"`
    Experience int       `bson:"experience" json:"experience" description:"经验"`
    Grade    int         `bson:"grade" json:"grade" description:"等级"`
    GradeDescribe string `bson:"grade_describe" json:"grade_describe" description:"等级描述"`
    NextExperience int   `bson:"next_experience" json:"next_experience" description:"下一级所要求经验"`
}

var UserGameSts = new(UserGameStsModel)

func (m *UserGameStsModel) Get(userId int) (task UserGameStsModel, err error) {
    err = Mongo.Chess.M(MongoDBStr, MongoColUserGameSts, func(c *mgo.Collection) error {
	query := bson.M{
	    "user_id": userId,
	}
	return c.Find(query).One(&task)
    })

    return
}
func (m *UserGameStsModel) Upsert(userId int, Task UserGameStsModel) error {

    return Mongo.Chess.M(MongoDBStr, MongoColUserGameSts, func(c *mgo.Collection) error {
	query := bson.M{
	    "user_id": userId,
	}
	changeInfo, err := c.Upsert(query, Task)
	// Debug
	log.Debugf("UserGameStsModel.Upsert", logrus.Fields{
	    "User ID":     userId,
	    "Query":       query,
	    "Change Info": changeInfo,
	    "Error":       err,
	})

	return err
    })
}
