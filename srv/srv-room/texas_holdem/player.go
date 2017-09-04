package texas_holdem

import (
	"chess/common/define"
	"chess/common/log"
	"chess/common/services"
	"chess/models"
	"chess/srv/srv-room/misc/packet"
	pb "chess/srv/srv-room/proto"
	"chess/srv/srv-room/registry"
	"errors"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"io"
	"time"
)

const (
	ActStandup = "standup" // 站起
	ActSitdown = "sitdown" // 坐下
	ActReady   = "ready"   // 准备
	ActBetting = "betting" // 下注中
	ActCall    = "call"    // 跟注
	ActCheck   = "check"   // 让牌
	ActRaise   = "raise"   // 加注
	ActFold    = "fold"    // 弃牌
	ActFlee    = "flee"    //逃跑
	ActAllin   = "allin"   // 全押
)

type Player struct {
	Id         int
	Nickname   string
	Avatar     string
	Level      string
	Chips      int // 带入牌桌的筹码
	TotalChips int // 原有筹码
	CurrChips  int // 牌桌上以外的筹码

	Pos    int
	Bet    int
	Action string
	Cards  Cards
	Hand   *Hand

	Table *Table

	ActBet       chan *pb.RoomPlayerBetReq
	timer        *time.Timer // action timer
	notOperating int         // 未操作次数

	Flag   int // 会话标记
	Stream pb.RoomService_StreamServer

	chatCancleFunc context.CancelFunc
}

func NewPlayer(id int, stream pb.RoomService_StreamServer) *Player {
	player := &Player{
		Id:     id,
		Hand:   NewHand(),
		ActBet: make(chan *pb.RoomPlayerBetReq),
		//Ipc:    ipc,
		Stream: stream,
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
	if p.Flag&define.PLAYER_DISCONNECT == 0 {
		log.Debugf("SendMessage to %d --- code:%d msg:%+v", p.Id, code, msg)
		message := &pb.Room_Frame{
			Type:    pb.Room_Message,
			Message: packet.Pack(code, msg),
		}
		if err := p.Stream.Send(message); err != nil {
			log.Error("p.stream.Send ", err)
		}
	}

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
		p.notOperating = 0
		return m, nil
	case <-p.Table.EndChan:
		return nil, nil
	case <-p.timer.C:
		p.notOperating++
		if p.notOperating >= 2 {
			p.Standup(true)
		}
		p.notOperating = 0
		return nil, errors.New("timeout")
	}
}

// @todo chat服务异常处理
func (p *Player) SubscribeChat() {
	table := p.Table
	if table == nil {
		return
	}

	conn, sid := services.GetService2(define.SRV_NAME_CHAT)
	if conn == nil {
		log.Error("cannot get centre service:", sid)
		return
	}
	cli := pb.NewChatServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := cli.Subscribe(ctx, &pb.Chat_Consumer{Id: table.Id, From: -1})
	if err != nil {
		log.Error("c.Subscribe error: ", err)
		return
	}

	p.chatCancleFunc = cancel
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error("Chat subscribe stream.Recv error: ", err)
			return
		}

		p.SendMessage(define.Code["room_table_chat_ack"], &pb.RoomTableChatAck{
			Id:     message.Id,
			Body:   message.Body,
			Offset: message.Offset,
		})
	}
}

func (p *Player) UnsubscribeChat() {
	if p.chatCancleFunc != nil {
		p.chatCancleFunc() // close chat subscribe stream
	}
}

func (p *Player) SendChatMessage(msg *pb.RoomTableChatReq) {
	table := p.Table
	if table == nil {
		return
	}

	if msg.Id != table.Id {
		log.Debug("SendChatMessage Invalid!!!")
		return
	}

	conn, sid := services.GetService2(define.SRV_NAME_CHAT)
	if conn == nil {
		log.Error("cannot get chat service:", sid)
		return
	}
	cli := pb.NewChatServiceClient(conn)
	_, err := cli.Send(context.Background(), &pb.Chat_Message{Id: msg.Id, Body: msg.Body})
	if err != nil {
		log.Errorf("Chat service cli.Send: %v", err)
	}
}

func (p *Player) Join(rid int, tid string) (table *Table) {
	if p.Table != nil {
		log.Debugf("玩家%d已在牌桌上", p.Id)
		return
	}

	//  获取玩家筹码
	var userWallet models.UsersWalletModel
	err := models.UsersWallet.Get(p.Id, &userWallet)
	if err != nil {
		log.Debugf("玩家%d获取筹码失败", p.Id)
		log.Error("models.UsersWallet.Get: ", err)
		return
	}
	p.TotalChips = int(userWallet.Balance)
	p.CurrChips = int(userWallet.Balance)

	table = GetTable(rid, tid)
	if table == nil {
		log.Debug("找不到牌桌")
		return
	}

	if p.CurrChips < table.MinCarry {
		log.Debugf("玩家%d筹码不足，要求筹码%d，当前筹码%d", p.Id, table.MinCarry, p.CurrChips)
		return nil
	}

	p.Bet = 0
	p.Cards = nil
	p.Hand.Init()
	p.Action = ActStandup
	p.Pos = 0
	p.Table = nil

	table.AddPlayer(p)

	// 订阅牌桌聊天室
	go p.SubscribeChat()

	// 带入筹码
	p.Chips = table.MaxCarry / 2
	p.CurrChips -= p.Chips

	log.Debugf("(%s)玩家%d加入牌桌, 位置%d, 当前牌桌有%d个玩家", table.Id, p.Id, p.Pos, table.N)

	// 2102, 通报加入游戏的玩家
	p.BroadcastAll(define.Code["room_player_join_ack"], &pb.RoomPlayerJoinAck{
		BaseAck: &pb.BaseAck{Ret: 1, Msg: "ok"},
		Player:  p.ToProtoMessage(),
	})

	log.Debug("推送房间信息")
	// 2006, 当玩家加入房间后，服务器会向此用户推送房间信息
	p.SendMessage(define.Code["room_get_table_ack"], &pb.RoomGetTableAck{
		BaseAck: &pb.BaseAck{Ret: 1, Msg: "ok"},
		Table:   table.ToProtoMessage(),
	})

	// 在线数+1
	go Pcounter(table.RoomId, 1)

	return
}

// 逃跑
func (p *Player) Flee() {
	table := p.Table
	if table == nil {
		log.Errorf("玩家%d不在牌桌上！", p.Id)
		return
	}

	// 牌局未结束
	if len(p.Cards) > 0 && table.gambling != nil {
		table.lock.Lock()
		defer table.lock.Unlock()
		if table.gambling.Players[p.Pos-1] != nil {
			table.gambling.Players[p.Pos-1].Action = ActFlee
			if table.gambling.Players[p.Pos-1].Actions[table.dealIdx] == nil {
				table.gambling.Players[p.Pos-1].Actions[table.dealIdx] = &models.ActionData{
					Action: ActFlee,
				}
			} else {
				table.gambling.Players[p.Pos-1].Actions[table.dealIdx].Action = ActFlee
			}
		}
	}
}

// 站起
func (p *Player) Standup(force bool) {
	table := p.Table
	if table == nil {
		log.Errorf("玩家%d不在牌桌上！", p.Id)
		return
	}

	if table.N == 1 && !force {
		log.Errorf("牌桌上只剩一位玩家，不允许站起操作！", p.Id)
		return
	}

	if p.Action == ActStandup {
		log.Errorf("玩家%d当前已是站起状态！", p.Id)
		return
	}
	if p.Action == ActBetting && !force { // 弃牌
		p.ActBet <- &pb.RoomPlayerBetReq{
			TableId: table.Id,
			Bet:     -1,
		}
	}

	p.Flee()
	table.DelPlayer(p, false)
	table.AddBystander(p)

	// 2113, 广播玩家站起
	table.BroadcastAll(define.Code["room_player_standup_ack"], &pb.RoomPlayerStandupAck{
		BaseAck:   &pb.BaseAck{Ret: 1, Msg: "ok"},
		TableId:   table.Id,
		PlayerId:  int32(p.Id),
		PlayerPos: int32(p.Pos),
		Force:     force,
	})
	// 2122, 通知自动坐下等待玩家数
	p.SendMessage(define.Code["room_player_autositdown_ack"], &pb.RoomPlayerAutoSitdownAck{
		Num:   int32(len(table.AutoSitdown)),
		Queue: table.AutoSitdown,
	})

	p.Bet = 0
	p.Cards = nil
	p.Hand.Init()
	p.Action = ActStandup
	p.Pos = 0
	p.CurrChips += p.Chips
	p.Chips = 0

}

// 自动坐下
func (p *Player) AutoSitdown() {
	table := p.Table
	if table == nil {
		log.Errorf("玩家%d不在牌桌上！", p.Id)
		return
	}

	if p.Action != ActStandup {
		log.Errorf("玩家%d当前不是站起状态！", p.Id)
		return
	}

	if table.N < table.Max {
		log.Errorf("玩家%d可直接坐下！", p.Id)
		return
	}

	table.AddAutoSitdown(p.Id)

	// 2122, 广播自动坐下等待玩家数
	table.BroadcastBystanders(define.Code["room_player_autositdown_ack"], &pb.RoomPlayerAutoSitdownAck{
		Num:   int32(len(table.AutoSitdown)),
		Queue: table.AutoSitdown,
	})
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

	if p.Action != ActStandup {
		log.Errorf("玩家%d当前不是站起状态！", p.Id)
		return
	}

	table.DelBystander(p, false)
	table.AddPlayer(p)

	// 带入筹码
	p.Chips = table.MaxCarry / 2
	p.CurrChips -= p.Chips

	// 2115, 广播玩家坐下
	table.BroadcastAll(define.Code["room_player_sitdown_ack"], &pb.RoomPlayerSitdownAck{
		BaseAck: &pb.BaseAck{Ret: 1, Msg: "ok"},
		Player:  p.ToProtoMessage(),
	})
}

// 换桌
func (p *Player) ChangeTable() {
	table := p.Table
	if table == nil {
		return
	}

	if p.Action == ActStandup || p.Action == ActSitdown || p.Action == ActFold {
		rid := table.RoomId
		tid := table.Id
		another := GetAnotherTable(rid, tid)
		if another != nil {
			p.Leave()
			p.Join(rid, another.Id)
		}
	} else {
		log.Debugf("(%s)玩家%d牌局未结束，不允许换桌", table.Id, p.Id)
	}
}

func (p *Player) Leave() (table *Table) {
	table = p.Table
	if table == nil {
		return
	}

	log.Debugf("(%s)玩家%d离开牌桌", table.Id, p.Id)

	// 持久化
	go func(total, curr int) {
		add := curr - total

		err := models.UsersWallet.Checkout(p.Id, add)
		if err != nil {
			log.Errorf("models.UsersWallet.Checkout(%d, %d) Error: %s", p.Id, add, err)
		}
	}(p.TotalChips, p.CurrChips+p.Chips)

	if p.Action != ActStandup {
		p.Flee()
		// 2104, 广播离开房间的玩家
		table.BroadcastAll(define.Code["room_player_gone_ack"], &pb.RoomPlayerGoneAck{
			BaseAck: &pb.BaseAck{Ret: 1, Msg: "ok"},
			Player:  p.ToProtoMessage(),
		})
		table.DelPlayer(p, true)
	} else {
		table.DelBystander(p, true)
	}

	p.Bet = 0
	p.Cards = nil
	p.Hand.Init()
	p.Action = ""
	p.Pos = 0
	p.CurrChips += p.Chips
	p.TotalChips = p.CurrChips
	p.Chips = 0
	p.Table = nil
	if p.timer != nil {
		p.timer.Reset(0)
	}

	// 在线数-1
	go Pcounter(table.RoomId, -1)
	go p.UnsubscribeChat()

	return
}

// 掉线处理
func (p *Player) Disconnect() {

	table := p.Table
	if table == nil { // 不在牌桌上
		log.Debugf("玩家%d掉线了(不在牌桌上)...", p.Id)
		registry.Unregister(p.Id, p)
	} else { // 在牌桌上
		log.Debugf("玩家%d掉线了(在牌桌上)...", p.Id)
		p.Flag |= define.PLAYER_DISCONNECT
	}
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
		Pos:       int32(p.Pos),
		Id:        int32(p.Id),
		Nickname:  p.Nickname,
		Avatar:    p.Avatar,
		Chips:     int32(p.Chips),
		Bet:       int32(p.Bet),
		Action:    p.Action,
		Cards:     p.Cards.ToProtoMessage(),
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
