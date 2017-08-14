package client_handler

import (
	"crypto/rc4"
	"fmt"
	"io"
	"math/big"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"chess/agent/misc/crypto/dh"
	"chess/agent/misc/packet"
	. "chess/common/define"
	"chess/common/log"
	"chess/common/services"

	pb "chess/agent/proto"

	"github.com/golang/protobuf/proto"

	. "chess/agent/types"
)

var Handlers map[int16]func(*Session, []byte) []byte

func init() {
	Handlers = map[int16]func(*Session, []byte) []byte{
		0:  P_heart_beat_req,
		10: P_user_login_req,
		30: P_get_seed_req,
	}
}

// 心跳包 直接把数据包发回去
func P_heart_beat_req(sess *Session, data []byte) []byte {
	req := &pb.AutoId{}
	err := proto.Unmarshal(data[2:], req)
	if err != nil {
		log.Error("P_heart_beat_req Unmarshal ERROR", err)
	}
	return packet.Pack(Code["heart_beat_ack"], req)
}

// 密钥交换
// 加密建立方式: DH+RC4
// 注意:完整的加密过程包括 RSA+DH+RC4
// 1. RSA用于鉴定服务器的真伪(这步省略)
// 2. DH用于在不安全的信道上协商安全的KEY
// 3. RC4用于流加密
func P_get_seed_req(sess *Session, data []byte) []byte {
	tbl := &pb.SeedInfo{}
	err := proto.Unmarshal(data[2:], tbl)
	if err != nil {
		log.Error("P_get_seed_req Unmarshal ERROR", err)
	}

	log.Debug("P_get_seed_req", data)
	log.Debugf("P_get_seed_req Unmarshal %+v", *tbl)
	// KEY1
	X1, E1 := dh.DHExchange()
	KEY1 := dh.DHKey(X1, big.NewInt(int64(tbl.ClientSendSeed)))

	// KEY2
	X2, E2 := dh.DHExchange()
	KEY2 := dh.DHKey(X2, big.NewInt(int64(tbl.ClientReceiveSeed)))

	ret := pb.SeedInfo{int32(E1.Int64()), int32(E2.Int64())}
	// 服务器加密种子是客户端解密种子
	encoder, err := rc4.NewCipher([]byte(fmt.Sprintf("%v%v", SALT, KEY2)))
	if err != nil {
		log.Error(err)
		return nil
	}
	decoder, err := rc4.NewCipher([]byte(fmt.Sprintf("%v%v", SALT, KEY1)))
	if err != nil {
		log.Error(err)
		return nil
	}
	sess.Encoder = encoder
	sess.Decoder = decoder
	sess.Flag |= SESS_KEYEXCG

	log.Debug("seed_info ---", ret)
	return packet.Pack(Code["get_seed_ack"], &ret)
}

// 玩家登陆过程
func P_user_login_req(sess *Session, data []byte) []byte {
	req := &pb.UserLoginReq{}
	err := proto.Unmarshal(data[2:], req)
	if err != nil {
		log.Error("P_user_login_req Unmarshal ERROR ", err)
	}

	log.Debug("P_user_login_req: ", req)

	// 登陆鉴权
	// 简单鉴权可以在agent直接完成，通常公司都存在一个用户中心服务器用于鉴权
	authConn, authServiceId := services.GetService2(SRV_NAME_AUTH)
	if authConn == nil {
		log.Error("cannot get auth service:", authServiceId)
		return nil
	}
	authCli := pb.NewAuthServiceClient(authConn)
	authRes, err := authCli.Auth(context.Background(), &pb.AuthArgs{UserId: req.UserId, Token: req.Token})
	if err != nil {
		log.Error("authCli.Auth: ", err)
		return packet.Pack(Code["user_login_ack"], &pb.UserLoginAck{&pb.BaseAck{Ret: SYSTEM_ERROR, Msg: "system error."}})
	}
	if authRes.Ret != 1 {
		return packet.Pack(Code["user_login_ack"], &pb.UserLoginAck{&pb.BaseAck{Ret: AUTH_FAIL, Msg: "Auth fail."}})
	}

	sess.UserId = req.UserId
	sess.Token = req.Token

	// 选择Room服务器
	// 选服策略依据业务进行，比如小服可以固定选取某台，大服可以采用HASH或一致性HASH
	var serviceId string
	var conn *grpc.ClientConn
	if req.ConnectTo != "" { // 客户端指定连接的服务
		serviceId = req.ConnectTo
		conn = services.GetServiceWithId(serviceId, SRV_NAME_ROOM)
	} else {
		conn, serviceId = services.GetService2(SRV_NAME_ROOM)
	}
	if conn == nil {
		log.Error("cannot get room service:", serviceId)
		return nil
	}

	cli := pb.NewRoomServiceClient(conn)

	// 开启到游戏服的流
	ctx := metadata.NewContext(
		context.Background(),
		metadata.New(map[string]string{
			"userid":       fmt.Sprint(sess.UserId),
			"service_name": SRV_NAME_ROOM,
			"service_id":   serviceId,
			"unique_id":    req.UniqueId,
			"is_reconnect": fmt.Sprint(req.IsReconnect),
		}),
	)
	stream, err := cli.Stream(ctx)
	if err != nil {
		log.Error(err)
		return nil
	}
	sess.Stream = stream
	sess.GSID = serviceId

	// 读取GAME返回消息的goroutine
	fetcher_task := func(sess *Session) {
		for {
			in, err := sess.Stream.Recv()
			if err == io.EOF { // 流关闭
				log.Debug(err)
				return
			}
			if err != nil {
				log.Error(err)
				return
			}
			select {
			case sess.MQ <- *in:
			case <-sess.Die:
			}
		}
	}
	go fetcher_task(sess)

	// ping
	go func(sess *Session) {
		for {
			frame := &pb.Room_Frame{
				Type:    pb.Room_Ping,
				Message: []byte{},
			}

			// check stream
			if sess.Stream == nil {
				return
			}

			if err := sess.Stream.Send(frame); err != nil {
				log.Error("Send room ping frame error:", err)
				return
			}

			time.Sleep(5 * time.Second)
		}
	}(sess)

	return packet.Pack(Code["user_login_ack"], &pb.UserLoginAck{
		BaseAck: &pb.BaseAck{Ret: 1, Msg: "ok"},
		ServiceId: serviceId,
	})
}
