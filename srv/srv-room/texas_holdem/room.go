package texas_holdem

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Tables struct {
	M       map[string]*Table
	counter int // 牌桌计数器
	lock    sync.Mutex
}

type Room struct {
	Id     int
	tables Tables
}

func NewRoom(rid int) *Room {
	if rooms[rid] == nil {
		rooms[rid] = &Room{
			Id: rid,
			tables: Tables{
				M:       make(map[string]*Table),
				counter: 0,
				lock:    sync.Mutex{},
			},
		}
	}
	return rooms[rid]
}

func (r *Room) SetTable(t *Table) {
	r.tables.lock.Lock()
	defer r.tables.lock.Unlock()

	r.setTable(t)
}

func (r *Room) setTable(t *Table) {
	if t.Id == "" {
		t.Id = fmt.Sprintf("%d-%d", r.Id, time.Now().Unix())
	}
	r.tables.M[t.Id] = t
	r.tables.counter++
}

func (r *Room) GetTable(tid string) *Table {
	r.tables.lock.Lock()
	defer r.tables.lock.Unlock()

	table := r.tables.M[tid]
	if table == nil {
		for _, v := range r.tables.M {
			if v.N < v.Max {
				return v
			}
		}
		table = NewTable(9, 5, 10)
		r.setTable(table)
	}

	return table
}

func (r *Room) DelTable(tid string) {
	r.tables.lock.Lock()
	defer r.tables.lock.Unlock()

	delete(r.tables.M, tid)
	r.tables.counter--
}

func (r *Room) Tables() map[string]*Table {
	return r.tables.M
}

var rooms = make(map[int]*Room)

// TODO 初始化房间列表
func init() {
	rooms[1] = &Room{
		Id:1,
		tables: Tables{
			M:       make(map[string]*Table),
			counter: 0,
			lock:    sync.Mutex{},
		},
	}
}

func DelTable(tid string) {
	tmp := strings.Split(tid, "-")
	if len(tmp) != 2 {
		return
	}

	rid, err := strconv.Atoi(tmp[0])
	if err != nil {
		return
	}

	if room, ok := rooms[int(rid)]; ok {
		room.DelTable(tid)
	}
}

func GetTable(rid int, tid string) *Table {
	tmp := strings.Split(tid, "-")
	if len(tmp) == 2 { // 根据tid来获取房间ID
		_rid, err := strconv.Atoi(tmp[0])
		if err == nil {
			rid = int(_rid)
		}
	}

	if room, ok := rooms[rid]; ok {
		return room.GetTable(tid)
	}
	return nil
}
