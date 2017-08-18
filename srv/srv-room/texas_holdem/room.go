package texas_holdem

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
	"chess/models"
	pb "chess/srv/srv-room/proto"
	"chess/common/log"
	"chess/common/define"
	"golang.org/x/net/context"
	"chess/common/services"
)

var serviceId string

type Tables struct {
	M       map[string]*Table
	counter int // 牌桌计数器
	pcounter int // 玩家计数器
	lock    sync.Mutex
}

type Room struct {
	Id     int
	BigBlind int
	SmallBlind int
	MinCarry int
	MaxCarry int
	Max int

	tables Tables
}

func NewRoom(rid, bb, sb, minC, maxC, max int) *Room {
	if RoomList[rid] == nil {
		RoomList[rid] = &Room{
			Id: rid,
			BigBlind: bb,
			SmallBlind: sb,
			MinCarry: minC,
			MaxCarry: maxC,
			Max: max,

			tables: Tables{
				M:       make(map[string]*Table),
				counter: 0,
				pcounter: 0,
				lock:    sync.Mutex{},
			},
		}
	}
	return RoomList[rid]
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
	log.Debugf("创建新牌桌(%s)", t.Id)
	r.tables.M[t.Id] = t
	r.tables.counter++
}

func (r *Room) pCounter(t int) {
	r.tables.lock.Lock()
	defer r.tables.lock.Unlock()

	r.tables.pcounter += t
	if r.tables.pcounter <0{
		r.tables.pcounter = 0
	}
}

func (r *Room) GetTableExists(tid string) *Table {
	r.tables.lock.Lock()
	defer r.tables.lock.Unlock()

	return r.tables.M[tid]
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
		table = NewTable(r.Id, r.Max, r.SmallBlind, r.BigBlind, r.MinCarry, r.MaxCarry)
		r.setTable(table)
	}

	return table
}

func (r *Room) GetAnotherTable(tid string) *Table {
	r.tables.lock.Lock()
	defer r.tables.lock.Unlock()

	for _, v := range r.tables.M {
		if v.Id != tid && v.N < v.Max {
			return v
		}
	}
	table := NewTable(r.Id, r.Max, r.SmallBlind, r.BigBlind, r.MinCarry, r.MaxCarry)
	r.setTable(table)

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

var RoomList = make(map[int]*Room)

// TODO 初始化房间列表
func InitRoomList(sid string) {
	serviceId = sid
	list, err := models.Rooms.GetAll()
	if err != nil {
		log.Errorf("models.Rooms.GetAll ERROR: %s", err)
	}

	for _, v := range list {
		log.Debugf("初始化游戏房间-%d", v.Id)
		NewRoom(v.Id, v.BigBlind, v.SmallBlind, v.MinCarry, v.MaxCarry, v.Max)
	}
}

func DelTable(tid string) {
	tmp := strings.Split(tid, "-")
	if len(tmp) < 2 {
		return
	}

	rid, err := strconv.Atoi(tmp[0])
	if err != nil {
		return
	}

	if room, ok := RoomList[int(rid)]; ok {
		room.DelTable(tid)
	}
}

func GetTableExists(rid int, tid string) *Table {
	tmp := strings.Split(tid, "-")
	if len(tmp) >= 2 { // 根据tid来获取房间ID
		_rid, err := strconv.Atoi(tmp[0])
		if err == nil {
			rid = int(_rid)
		}
	}
	if room, ok := RoomList[rid]; ok {
		return room.GetTableExists(tid)
	}
	return nil
}

func GetTable(rid int, tid string) *Table {
	tmp := strings.Split(tid, "-")
	if len(tmp) >= 2 { // 根据tid来获取房间ID
		_rid, err := strconv.Atoi(tmp[0])
		if err == nil {
			rid = int(_rid)
		}
	}
	if room, ok := RoomList[rid]; ok {
		return room.GetTable(tid)
	}
	return nil
}

// 获取其他牌桌
func GetAnotherTable(rid int, tid string) *Table {
	tmp := strings.Split(tid, "-")
	if len(tmp) >= 2 { // 根据tid来获取房间ID
		_rid, err := strconv.Atoi(tmp[0])
		if err == nil {
			rid = int(_rid)
		}
	}

	if room, ok := RoomList[rid]; ok {
		return room.GetAnotherTable(tid)
	}
	return nil
}

//
func Pcounter(rid, t int) {
	if room, ok := RoomList[rid]; ok {
		room.pCounter(t)

		conn, sid := services.GetService2(define.SRV_NAME_CENTRE)
		if conn == nil {
			log.Error("cannot get centre service:", sid)
			return
		}
		cli := pb.NewCentreServiceClient(conn)
		_, err := cli.UpdateRoomInfo(
			context.Background(),
			&pb.UpdateRoomInfoArgs{
				ServiceId: serviceId,
				RoomId: int32(rid),
				RoomInfo: &pb.RoomInfo{
					TableNumber: int32(room.tables.counter),
					PlayerNumber: int32(room.tables.pcounter),
				},
			},
		)
		if err != nil {
			log.Error("cli.UpdateRoomInfo: ", err)
		}

	}
}