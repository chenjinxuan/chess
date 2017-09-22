package models

import (
	"chess/common/log"
	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserBagMongoModel struct {
	UserId int            `bson:"user_id" json:"user_id"`
	List   []UserBagModel `bson:"list" "json:"list"`
}
type UserBagModel struct {
	Number  int `bson:"number" json:"number"`
	GoodsId int `bson:"goods_id" json:"goods_id"`
}

var UserBag = new(UserBagModel)

func (m *UserBagModel) Get(userId int) (bag UserBagMongoModel, err error) {
	err = Mongo.Chess.M(MongoDBStr, MongoColUserBag, func(c *mgo.Collection) error {
		query := bson.M{
			"user_id": userId,
		}
		return c.Find(query).One(&bag)
	})
	return
}

func (m *UserBagModel) Upsert(userId int, userBag UserBagMongoModel) error {
	return Mongo.Chess.M(MongoDBStr, MongoColUserBag, func(c *mgo.Collection) error {
		query := bson.M{
			"user_id": userId,
		}
		changeInfo, err := c.Upsert(query, userBag)

		// Debug
		log.Debugf("UserBagMongoModel.Upsert", logrus.Fields{
			"User ID":     userId,
			"Query":       query,
			"Change Info": changeInfo,
			"Error":       err,
		})

		return err
	})
}

func (m *UserBagModel) UpdateInc(userId, GoodsId, val int) error {
	return Mongo.Chess.M(MongoDBStr, MongoColUserBag, func(c *mgo.Collection) error {
		query := bson.M{
			"user_id":       userId,
			"list.goods_id": GoodsId,
		}
		update := bson.M{
			"list.$.number": val,
		}
		err := c.Update(query, bson.M{"$inc": update})

		return err
	})
}

func (m *UserBagModel) RemoveByGoodsId(userId, GoodsId int) error {
	return Mongo.Chess.M(MongoDBStr, MongoColUserBag, func(c *mgo.Collection) error {
		query := bson.M{
			"user_id": userId,
		}
		update := bson.M{
			"list": bson.M{
				"goods_id": GoodsId,
			},
		}
		err := c.Update(query, bson.M{"$pull": update})

		return err
	})
}
