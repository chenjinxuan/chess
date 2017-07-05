package main

import (
	"encoding/binary"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/yuin/gopher-lua"
	"gopkg.in/mgo.v2"
	"log"
	"path/filepath"
	"sort"
	"time"
)

const (
	BOLTDB_BUCKET = "REDOLOG"
	LAYOUT        = "2006-01-02T15:04:05"
)

type rec struct {
	db_idx int    // file
	key    uint64 // key of file
}

type ToolBox struct {
	L       *lua.LState // the lua virtual machine
	dbs     []*bolt.DB  // all opened boltdb
	recs    []rec
	mgo     *mgo.Session
	mgo_url string
}

type file_sort []string

func (a file_sort) Len() int      { return len(a) }
func (a file_sort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a file_sort) Less(i, j int) bool {
	layout := "REDO-2006-01-02T15:04:05.RDO"
	tm_a, _ := time.Parse(layout, a[i])
	tm_b, _ := time.Parse(layout, a[j])
	return tm_a.Unix() < tm_b.Unix()
}

func NewToolBox(dir string) *ToolBox {
	t := new(ToolBox)
	// lookup *.RDO
	files, err := filepath.Glob(dir + "/*.RDO")
	if err != nil {
		log.Println(err)
		return nil
	}
	// sort by creation time
	sort.Sort(file_sort(files))

	// open all db
	for _, file := range files {
		db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 2 * time.Second, ReadOnly: true})
		if err != nil {
			log.Println(err)
			continue
		}
		t.dbs = append(t.dbs, db)
	}

	// reindex all keys
	log.Println("loading database")
	for i := range t.dbs {
		t.dbs[i].View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(BOLTDB_BUCKET))
			c := b.Cursor()
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				t.recs = append(t.recs, rec{i, binary.BigEndian.Uint64(k)})
			}
			return nil
		})
	}

	// init lua machine
	log.Println("init lua machine")
	t.L = lua.NewState()
	// register
	t.register()
	log.Println("ready")

	t.L.DoString("help()")
	return t
}

func (t *ToolBox) Close() {
	t.L.Close()
	for _, db := range t.dbs {
		db.Close()
	}
	if t.mgo != nil {
		t.mgo.Close()
	}
}

func (t *ToolBox) register() {
	mt := t.L.NewTypeMetatable("mt_reclist")
	t.L.SetGlobal("mt_reclist", mt)
	t.L.SetField(mt, "__index", t.L.SetFuncs(t.L.NewTable(), map[string]lua.LGFunction{
		"get":    t.builtin_get,
		"length": t.builtin_length,
		"mgo":    t.builtin_mgo,
		"replay": t.builtin_replay,
	}))

	Int64(0).register(t.L)

	// global variable
	ud := t.L.NewUserData()
	ud.Value = t.recs
	t.L.SetGlobal("redo", ud)
	t.L.SetMetatable(ud, t.L.GetTypeMetatable("mt_reclist"))

	// global function
	t.L.SetGlobal("help", t.L.NewFunction(t.builtin_help))
}

func (t *ToolBox) exec(cmd string) {
	if err := t.L.DoString(cmd); err != nil {
		fmt.Println(err)
	}
}
