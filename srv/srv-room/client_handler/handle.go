package client_handler

import (
	. "chess/common/define"
	"chess/common/log"
	"chess/srv/srv-room/misc/packet"
	pb "chess/srv/srv-room/proto"
	"chess/srv/srv-room/registry"
	. "chess/srv/srv-room/texas_holdem"
	"github.com/golang/protobuf/proto"
)

var Handlers map[int16]func(*Player, []byte) []byte

func init() {
	Handlers = map[int16]func(*Player, []byte) []byte{
		2001: P_room_ping_req,
		2005: P_room_get_table_req,

		2101: P_room_player_join_req,
		2103: P_room_player_gone_req,
		2105: P_room_player_bet_req,
		2112: P_room_player_standup_req,
		2114: P_room_player_sitdown_req,
		2116: P_room_player_change_table_req,
		2118: P_room_player_logout_req,
	}
}

//----------------------------------- ping
func P_room_ping_req(p *Player, data []byte) []byte {
	tbl := &pb.AutoId{}
	proto.Unmarshal(data, tbl)
	return packet.Pack(Code["room_ping_ack"], tbl)
}

// 查询牌桌信息
func P_room_get_table_req(p *Player, data []byte) []byte {
	ack := &pb.RoomGetTableAck{
		BaseAck: &pb.BaseAck{Ret: 0, Msg: ""},
	}

	req := &pb.RoomGetTableReq{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Errorf("proto.Unmarshal Error: %s", err)
		ack.BaseAck.Msg = "wrong data"
		return packet.Pack(Code["room_get_table_ack"], ack)
	}

	table := GetTable(int(req.RoomId), req.TableId)
	if err == nil {
		ack.BaseAck.Msg = "table not found"
		return packet.Pack(Code["room_get_table_ack"], ack)
	}

	ack.BaseAck.Ret = 1
	ack.Table = table.ToProtoMessage()
	return packet.Pack(Code["room_get_table_ack"], ack)
}

// 玩家加入游戏
func P_room_player_join_req(p *Player, data []byte) []byte {
	ack := &pb.RoomPlayerJoinAck{
		BaseAck: &pb.BaseAck{Ret: 0, Msg: ""},
	}

	//req := &pb.RoomGetTableAck{}
	req := &pb.RoomPlayerJoinReq{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Errorf("proto.Unmarshal Error: %s", err)
		return nil
	}
	log.Debug("P_room_player_join_req: ", req)
	table := p.Join(int(req.RoomId), req.TableId)
	if table == nil {
		log.Error("table not found")
		ack.BaseAck.Msg = "table not found"
		return packet.Pack(Code["room_player_join_ack"], ack)
	}

	return nil
}

// 玩家离开牌桌
func P_room_player_gone_req(p *Player, data []byte) []byte {
	req := &pb.RoomPlayerGoneReq{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Errorf("proto.Unmarshal Error: %s", err)
		return nil
	}

	p.Leave()

	return nil
}

// 玩家下注
func P_room_player_bet_req(p *Player, data []byte) []byte {
	req := &pb.RoomPlayerBetReq{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Errorf("proto.Unmarshal Error: %s", err)
		return nil
	}

	p.ActBet <- req

	return nil
}

// 玩家站起
func P_room_player_standup_req(p *Player, data []byte) []byte {
	req := &pb.RoomPlayerStandupReq{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Errorf("proto.Unmarshal Error: %s", err)
		return nil
	}
	log.Debug("P_room_player_standup_req", req)
	p.Standup()
	return nil
}

// 玩家坐下
func P_room_player_sitdown_req(p *Player, data []byte) []byte {
	req := &pb.RoomPlayerSitdownReq{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Errorf("proto.Unmarshal Error: %s", err)
		return nil
	}
	p.Sitdown()
	return nil
}

// 换桌
func P_room_player_change_table_req(p *Player, data []byte) []byte {
	req := &pb.RoomPlayerChangeTableReq{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Errorf("proto.Unmarshal Error: %s", err)
		return nil
	}
	p.ChangeTable()
	return nil
}

// 玩家登出
func P_room_player_logout_req(p *Player, data []byte) []byte {
	req := &pb.RoomPlayerLogoutReq{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Errorf("proto.Unmarshal Error: %s", err)
		return nil
	}
	log.Debug("P_room_player_logout_req", req)
	p.Leave()
	registry.Unregister(p.Id, p)
	return nil
}
