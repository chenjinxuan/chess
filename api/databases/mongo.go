package databases

import (
	"gopkg.in/mgo.v2"
	"time"
	"treasure/log"
)

var (
	Mongo *MongoDB
)

type MongoDB struct {
	Host, User, Pass string
	session          *mgo.Session
	dialTimeout      time.Duration
}

func NewMongoDB(host, user, pass string) *MongoDB {
	Mongo = &MongoDB{host, user, pass, nil, 0}
	return Mongo
}

func (db *MongoDB) SetDialTimeout(timeout time.Duration) {
	db.dialTimeout = timeout
}

func (db *MongoDB) SetMgoLogger(logger *log.WrapLog) {
	mgo.SetLogger(logger)
	mgo.SetDebug(true)
}

func (db *MongoDB) M(database, collection string, f func(*mgo.Collection) error) error {
	session, err := db.Session(database)
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

func (db *MongoDB) Session(database string) (*mgo.Session, error) {
	if db.session == nil {
		info := &mgo.DialInfo{
			Addrs:    []string{db.Host},
			Timeout:  db.dialTimeout * time.Second,
			Database: database,
			Username: db.User,
			Password: db.Pass,
		}
		session, err := mgo.DialWithInfo(info)
		if err != nil {
			return nil, err
		}
		db.session = session

		//session, err := mgo.DialWithTimeout(db.Host, db.dialTimeout*time.Second)
		//if err != nil {
		//	return nil, err
		//}
		//db.session = session

		//to be verified
		//db.session.SetSafe(&mgo.Safe{
		//	W:     1,
		//	WMode: "majority",
		//	J:     true,
		//})
	}
	return db.session.Copy(), nil
}

func (db *MongoDB) DB(session *mgo.Session, name string) (*mgo.Database, error) {
	mdb := session.DB(name)
	//if db.User != "" && db.Pass != "" {
	//	if err := mdb.Login(db.User, db.Pass); err != nil {
	//		return nil, err
	//	}
	//}

	return mdb, nil
}
