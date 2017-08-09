package define

// 网关协议号定义
var Code = map[string]int16{
	"heart_beat_req": 0,  // 心跳包..
	"heart_beat_ack": 1,  // 心跳包回复
	"user_login_req": 10, // 登陆
	"user_login_ack": 11, // 登陆成功
	"get_seed_req":   30, // socket通信加密使用
	"get_seed_ack":   31, // socket通信加密使用
	"kicked_out_ack": 40, // 踢出通知

	"centre_ping_req": 1001, //  ping
	"centre_ping_ack": 1002, //  pong

	"room_ping_req":                2001, //  ping
	"room_ping_ack":                2002, //  ping回复
	"room_set_table_req":           2003, // 创建牌桌
	"room_set_table_ack":           2004, // 创建牌桌回复
	"room_get_table_req":           2005, // 查询牌桌信息
	"room_get_table_ack":           2006, // 查询牌桌信息回复 (当玩家加入牌桌后，服务器会向此用户推送牌桌信息)
	"room_get_player_req":          2007, // 查询玩家信息
	"room_get_player_ack":          2008, // 查询玩家信息回复
	"room_player_join_req":         2101, // 玩家加入游戏
	"room_player_join_ack":         2102, // 通报加入游戏的玩家
	"room_player_gone_req":         2103, // 玩家离开牌桌
	"room_player_gone_ack":         2104, // 服务器广播离开房间的玩家
	"room_player_bet_req":          2105, // 玩家下注
	"room_player_bet_ack":          2106, // 玩家下注结果
	"room_button_ack":              2107, // 通报本局庄家 (服务器广播此消息，代表游戏开始并确定本局庄家)
	"room_deal_ack":                2108, // 发牌 - 共有四轮发牌，按顺序分别为：preflop (底牌), flop (翻牌), turn (转牌), river(河牌)
	"room_pot_ack":                 2109, // 通报奖池
	"room_action_ack":              2110, // 通报当前下注玩家
	"room_showdown_ack":            2111, // 摊牌和比牌
	"room_player_standup_req":      2112, // 玩家站起
	"room_player_standup_ack":      2113, // 玩家站起通报
	"room_player_sitdown_req":      2114, // 玩家坐下
	"room_player_sitdown_ack":      2115, // 玩家坐下通报
	"room_player_change_table_req": 2116, // 玩家换桌
	"room_shutdown_table_ack":      2117, // 关闭牌桌，服务进行维护时通报
}

var RCode = map[int16]string{
	0:  "heart_beat_req", // 心跳包..
	1:  "heart_beat_ack", // 心跳包回复
	10: "user_login_req", // 登陆
	11: "user_login_ack", // 登陆成功
	30: "get_seed_req",   // socket通信加密使用
	31: "get_seed_ack",   // socket通信加密使用
	40: "kicked_out_ack", // 踢出通知

	1001: "centre_ping_req", //  ping
	1002: "centre_ping_ack", //  pong

	2001: "room_ping_req",                //  ping
	2002: "room_ping_ack",                //  ping回复
	2003: "room_set_table_req",           // 创建牌桌
	2004: "room_set_table_ack",           // 创建牌桌回复
	2005: "room_get_table_req",           // 查询牌桌信息
	2006: "room_get_table_ack",           // 查询牌桌信息回复 (当玩家加入牌桌后，服务器会向此用户推送牌桌信息)
	2007: "room_get_player_req",          // 查询玩家信息
	2008: "room_get_player_ack",          // 查询玩家信息回复
	2101: "room_player_join_req",         // 玩家加入游戏
	2102: "room_player_join_ack",         // 通报加入游戏的玩家
	2103: "room_player_gone_req",         // 玩家离开牌桌
	2104: "room_player_gone_ack",         // 服务器广播离开房间的玩家
	2105: "room_player_bet_req",          // 玩家下注
	2106: "room_player_bet_ack",          // 玩家下注结果
	2107: "room_button_ack",              // 通报本局庄家 (服务器广播此消息，代表游戏开始并确定本局庄家)
	2108: "room_deal_ack",                // 发牌 - 共有四轮发牌，按顺序分别为：preflop (底牌), flop (翻牌), turn (转牌), river(河牌)
	2109: "room_pot_ack",                 // 通报奖池
	2110: "room_action_ack",              // 通报当前下注玩家
	2111: "room_showdown_ack",            // 摊牌和比牌
	2112: "room_player_standup_req",      // 玩家站起
	2113: "room_player_standup_ack",      // 玩家站起通报
	2114: "room_player_sitdown_req",      // 玩家坐下
	2115: "room_player_sitdown_ack",      // 玩家坐下通报
	2116: "room_player_change_table_req", // 玩家换桌
	2117: "room_shutdown_table_ack",      // 关闭牌桌，服务进行维护时通报
}

const (
	SALT = "CHESS_DH"

	AUTH_FAIL    = -999
	SYSTEM_ERROR = -500
)
