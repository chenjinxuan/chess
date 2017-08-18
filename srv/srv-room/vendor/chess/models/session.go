package models

import (
	"chess/common/log"
	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

const (
	MongoDBStr      = "chess"
	MongoColSession = "session"
)

var Session = new(SessionModel)

type SessionModel struct {
	UserId       int           `bson:"user_id"`
	From         string        `bson:"from"`
	UniqueId     string        `bson:"unique_id"`
	Token        *SessionToken `bson:"token"`
	RefreshToken string        `bson:"refresh_token"`
	Updated      time.Time     `bson:"updated"`
	Created      time.Time     `bson:"created"`
}

type SessionToken struct {
	Data   string `bson:"data"`
	Expire int64  `bson:"expire"`
}

func (m *SessionModel) Get(userId int, from, uniqueId string) (*SessionModel, error) {
	var session = new(SessionModel)

	err := Mongo.Chess.M(MongoDBStr, MongoColSession, func(c *mgo.Collection) error {
		query := bson.M{
			"user_id":   userId,
			"from":      from,
			"unique_id": uniqueId,
		}
		return c.Find(query).One(&session)
	})
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (m *SessionModel) Upsert(userId int, from, uniqueId string, session *SessionModel) error {
	from = strings.ToLower(from) // @todo
	return Mongo.Chess.M(MongoDBStr, MongoColSession, func(c *mgo.Collection) error {
		query := bson.M{
			"user_id": userId,
			//"from":      from,
			//"unique_id": uniqueId,
		}
		changeInfo, err := c.Upsert(query, session)

		// Debug
		log.Debugf("SessionModel.Upsert", logrus.Fields{
			"User ID":     userId,
			"From":        from,
			"Query":       query,
			"Change Info": changeInfo,
			"Error":       err,
		})

		return err
	})
}

func (m *SessionModel) RemoveByUid(userId int) error {
	return Mongo.Chess.M(MongoDBStr, MongoColSession, func(c *mgo.Collection) error {
		query := bson.M{
			"user_id": userId,
		}
		changeInfo, err := c.RemoveAll(query)
		log.Debugf("SessionModel.RemoveByUid", logrus.Fields{
			"User ID":     userId,
			"Change Info": changeInfo,
			"Error":       err,
		})
		return err
	})
}
