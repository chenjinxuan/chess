package nsqredo

import (
	"bytes"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	ENV_NSQD         = "NSQD_HOST"
	DEFAULT_PUB_ADDR = "http://172.17.42.1:4151/pub?topic=REDOLOG"
	MIME             = "application/octet-stream"
)

// a data change
type Change struct {
	Collection string // collection
	Field      string // field "a.b.c.d"
	Doc        interface{}
}

// a redo record represents complete transaction
type RedoRecord struct {
	API     string   // the api name
	UID     int32    // userid
	TS      uint64   // timestamp should get from snowflake
	Changes []Change // changes
}

var (
	_pub_addr string
	_prefix   string
	_ch       chan []byte
)

func init() {
	// get nsqd publish address
	_pub_addr = DEFAULT_PUB_ADDR
	if env := os.Getenv(ENV_NSQD); env != "" {
		_pub_addr = env + "/pub?topic=REDOLOG"
	}
	_ch = make(chan []byte, 4096)
	go publish_task()
}

// add a change with o(old value) and n(new value)
func (r *RedoRecord) AddChange(collection, field string, doc interface{}) {
	r.Changes = append(r.Changes, Change{Collection: collection, Field: field, Doc: doc})
}

func NewRedoRecord(uid int32, api string, ts uint64) *RedoRecord {
	return &RedoRecord{UID: uid, API: api, TS: ts}
}

func publish_task() {
	for {
		// post to nsqd
		bts := <-_ch
		resp, err := http.Post(_pub_addr, MIME, bytes.NewReader(bts))
		if err != nil {
			log.Println(err)
			continue
		}

		// read response
		if _, err := ioutil.ReadAll(resp.Body); err != nil {
			log.Println(err)
		}

		// close
		resp.Body.Close()
	}
}

// publish to nsqd (localhost nsqd is suggested!)
func Publish(r *RedoRecord) {
	// pack message
	if bts, err := bson.Marshal(r); err == nil {
		_ch <- bts
	} else {
		log.Println(err, r)
		return
	}
}
