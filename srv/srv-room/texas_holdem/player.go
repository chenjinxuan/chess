package texas_holdem

import (
	"chess/common/define"
	"chess/srv/srv-room/misc/packet"
	pb "chess/srv/srv-room/proto"
	"errors"
	"github.com/golang/protobuf/proto"
	"time"
	"chess/common/log"
)

const (
	ActStandup = "standup" // 站起
	ActSitdown = "sitdown" // 坐下
	ActReady = "ready" // 准备
	ActBetting = "betting" // 下注中
	ActCall  = "call"  // 跟注
	ActCheck = "check" // 让牌
	ActRaise = "raise" // 加注
	ActFold  = "fold"  // 弃牌
	ActAllin = "allin" // 全押
)

type Player struct {
	Id       int
	Nickname string
	Avatar   string
	Level    string
	Chips    int

	Pos    int
	Bet    int
	Action string
	Cards  Cards
	Hand   *Hand

	Table *Table

	ActBet chan *pb.RoomPlayerBetReq
	timer  *time.Timer // action timer

	Flag int // 会话标记
	stream pb.RoomService_StreamServer
	Ipc  chan *pb.Room_Frame
}

func NewPlayer(id int, stream pb.RoomService_StreamServer) *Player {
	player := &Player{
		Id:     id,
		Hand:   NewHand(),
		ActBet: make(chan *pb.RoomPlayerBetReq),
		//Ipc:    ipc,
		stream: stream,
	}
	player.Hand.Init()
	return player
}

func (p *Player) Broadcast(code int16, msg proto.Message) {
	if p.Table == nil {
		log.Debug("p.Broadcast Table nil")
		return
	}

	for _, oc := range p.Table.Players {
		if oc != nil && oc.Id != p.Id {
			oc.SendMessage(code, msg)
		}
	}
}

func (p *Player) BroadcastBystanders(code int16, msg proto.Message) {
	if p.Table == nil {
		log.Debug("p.Broadcast Table nil")
		return
	}

	for _, oc := range p.Table.Bystanders {
		if oc != nil && oc.Id != p.Id {
			oc.SendMessage(code, msg)
		}
	}
}

func (p *Player) BroadcastAll(code int16, msg proto.Message) {
	if p.Table == nil {
		log.Debug("p.Broadcast Table nil")
		return
	}

	for _, oc := range p.Table.Players {
		if oc != nil && oc.Id != p.Id {
			oc.SendMessage(code, msg)
		}
	}

	for _, oc := range p.Table.Bystanders {
		if oc != nil && oc.Id != p.Id {
			oc.SendMessage(code, msg)
		}
	}
}

func (p *Player) SendMessage(code int16, msg proto.Message) {
	log.Debugf("SendMessage to %d --- code:%d msg:%+v", p.Id, code, msg)
	message := &pb.Room_Frame{
		Type:    pb.Room_Message,
		Message: packet.Pack(code, msg),
	}
	if err := p.stream.Send(message); err != nil {
		log.Error("p.stream.Send ", err)
	}
	//p.Ipc <- message
}

func (p *Player) Betting(n int) (raised bool) {
	table := p.Table
	if table == nil {
		return
	}

	if n < 0 { // 弃牌
		p.Action = ActFold
		p.Cards = nil
		p.Hand.Init()
		n = 0
	} else if n == 0 { // 让牌
		p.Action = ActCheck
	} else if n+p.Bet <= table.Bet { // 跟注 或者 allin  table.Bet保持不变
		p.Action = ActCall
		p.Chips -= n
		p.Bet += n
	} else { // 加注
		p.Action = ActRaise
		p.Chips -= n
		p.Bet += n
		table.Bet = p.Bet
		raised = true
	}
	if p.Chips == 0 {
		p.Action = ActAllin
	}
	table.Chips[p.Pos-1] += int32(n)

	return
}

// 等待 获取玩家下注操作
func (p *Player) GetActionBet(timeout time.Duration) (*pb.RoomPlayerBetReq, error) {
	p.timer = time.NewTimer(timeout)

	select {
	case m := <-p.ActBet:
		return m, nil
	case <-p.Table.EndChan:
		return nil, nil
	case <-p.timer.C:
		return nil, errors.New("timeout")
	}
}

func (p *Player) Join(rid int, tid string) (table *Table) {
	table = GetTable(rid, tid)
	if table == nil {
		log.Debug("找不到牌桌")
		return
	}

	p.Bet = 0
	p.Cards = nil
	p.Hand.Init()
	p.Action = ActStandup
	p.Pos = 0
	p.Table = nil

	table.AddPlayer(p)

	log.Debugf("(%s)玩家%d加入牌桌, 位置%d, 当前牌桌有%d个玩家", table.Id, p.Id, p.Pos, table.N)

	// 2102, 通报加入游戏的玩家
	p.BroadcastAll(define.Code["room_player_join_ack"], &pb.RoomPlayerJoinAck{
		BaseAck: &pb.BaseAck{Ret: 1, Msg: "ok"},
		Player:  p.ToProtoMessage(),
	})

	log.Debug("推送房间信息")
	// 2006, 当玩家加入房间后，服务器会向此用户推送房间信息
	p.SendMessage(define.Code["room_get_table_ack"],  &pb.RoomGetTableAck{
		BaseAck: &pb.BaseAck{Ret: 1, Msg: "ok"},
		Table:   table.ToProtoMessage(),
	})

	return
}

// 站起
func (p *Player) Standup() {
	table := p.Table
	if table == nil {
		log.Errorf("玩家%d不在牌桌上！", p.Id)
		return
	}

	if table.N == 1 {
		log.Errorf("牌桌上只剩一位玩家，不允许站起操作！", p.Id)
		return
	}

	table.DelPlayer(p)
	table.AddBystander(p)

	// 2113, 广播玩家站起
	table.BroadcastAll(define.Code["room_player_standup_ack"], &pb.RoomPlayerStandupAck{
		BaseAck: &pb.BaseAck{Ret: 1, Msg: "ok"},
		TableId:  table.Id,
		PlayerId: int32(p.Id),
		PlayerPos: int32(p.Pos),
	})

	p.Bet = 0
	p.Cards = nil
	p.Hand.Init()
	p.Action = ActStandup
	p.Pos = 0

}

// 坐下
func (p *Player) Sitdown() {
	table := p.Table
	if table == nil {
		log.Errorf("玩家%d不在牌桌上！", p.Id)
		return
	}

	if table.N == table.Max {
		log.Errorf("牌桌上无空位，不允许坐下操作！", p.Id)
		return
	}

	table.DelBystander(p)
	table.AddPlayer(p)

	// 2115, 广播玩家坐下
	table.BroadcastAll(define.Code["room_player_sitdown_ack"], &pb.RoomPlayerSitdownAck{
		BaseAck: &pb.BaseAck{Ret: 1, Msg: "ok"},
		TableId:  table.Id,
		PlayerId: int32(p.Id),
		PlayerPos: int32(p.Pos),
	})
}

func (p *Player) Leave() (table *Table) {
	table = p.Table
	if table == nil {
		return
	}

	log.Debugf("(%s)玩家%d离开牌桌", table.Id, p.Id)

	if p.Action != ActStandup {
		// 2104, 广播离开房间的玩家
		table.BroadcastAll(define.Code["room_player_gone_ack"], &pb.RoomPlayerGoneAck{
			BaseAck: &pb.BaseAck{Ret: 1, Msg: "ok"},
			Player:  p.ToProtoMessage(),
		})
		table.DelPlayer(p)
	}

	p.Bet = 0
	p.Cards = nil
	p.Hand.Init()
	p.Action = ""
	p.Pos = 0
	p.Table = nil
	if p.timer != nil {
		p.timer.Reset(0)
	}

	return
}

func (p *Player) Next() *Player {
	table := p.Table
	if table == nil {
		return nil
	}

	for i := (p.Pos) % table.Cap(); i != p.Pos-1; i = (i + 1) % table.Cap() {
		if table.Players[i] != nil && table.Players[i].Action == ActReady {
			return table.Players[i]
		}
	}

	return nil
}

func (p *Player) ToProtoMessage() *pb.PlayerInfo {
	return &pb.PlayerInfo{
		Pos:      int32(p.Pos),
		Id:       int32(p.Id),
		Nickname: p.Nickname,
		Avatar:   p.Avatar,
		Chips:    int32(p.Chips),
		Bet:      int32(p.Bet),
		Action:   p.Action,
		Cards:    p.Cards.ToProtoMessage(),
		HandLevel: int32(p.Hand.Level),
	}
}

type Players []*Player

func (ps Players) ToProtoMessage() []*pb.PlayerInfo {
	_players := make([]*pb.PlayerInfo, len(ps), MaxN)
	for k, v := range ps {
		if v != nil {
			tmp := *v
			_player := tmp.ToProtoMessage()
			_players[k] = _player
		} else {
			_players[k] = &pb.PlayerInfo{}
		}
	}
	return _players
}
