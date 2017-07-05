package db

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"time"
	"chess/common/config"
)

func InitMongo() {
	for dbName, server := range config.Db.Mongo.Server {
		mongoMap[dbName] = NewMongoDB(server.Host, server.Username, server.Password,
			config.Db.Mongo.Setting.DialTimeout)
	}
}

func M(name string) (*MongoDB, error) {
	db, ok := mongoMap[name]
	if !ok {
		return nil, fmt.Errorf("mongo %s not found", name)
	}
	return db, nil
}

type MongoDB struct {
	Host, User, Pass string
	session          *mgo.Session
	dialTimeout      time.Duration
}

func NewMongoDB(host, user, pass string, dialTimeout time.Duration) *MongoDB {
	m := &MongoDB{host, user, pass, nil, 0}
	m.SetDialTimeout(dialTimeout)
	return m
}

func (db *MongoDB) SetDialTimeout(timeout time.Duration) {
	db.dialTimeout = timeout
}

func (db *MongoDB) M(database, collection string, f func(*mgo.Collection) error) error {
	session, err := db.Session()
	if err != nil {
		return err
	}
	defer session.Close()

	mdb, err := db.DB(session, database)
	if err != nil {
		return err
	}

	c := mdb.C(collection)
	return f(c)
}

func (db *MongoDB) Session() (*mgo.Session, error) {
	if db.session == nil {
		session, err := mgo.DialWithTimeout(db.Host, db.dialTimeout*time.Second)
		if err != nil {
			return nil, err
		}
		db.session = session
	}
	return db.session.Copy(), nil
}

func (db *MongoDB) DB(session *mgo.Session, name string) (*mgo.Database, error) {
	mdb := session.DB(name)
	if db.User != "" && db.Pass != "" {
		if err := mdb.Login(db.User, db.Pass); err != nil {
			return nil, err
		}
	}

	return mdb, nil
}
