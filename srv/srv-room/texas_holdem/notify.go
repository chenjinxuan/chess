package texas_holdem

import (
	. "chess/common/define"
	"chess/common/log"
	"chess/common/services"
	"chess/models"
	pb "chess/srv/srv-room/proto"
	"golang.org/x/net/context"
)

func Convert2GameTableInfoArgs(g *models.GamblingModel) *pb.GameTableInfoArgs {
	args := &pb.GameTableInfoArgs{
		RoomId:  int32(g.RoomId),
		TableId: g.TableId,
		Max:     int32(g.Max),
		Start:   int32(g.Start),
		End:     int32(g.End),
		Button:  int32(g.Button),
		Sb:      int32(g.Sb),
		Bb:      int32(g.Bb),
		SbPos:   int32(g.SbPos),
		BbPos:   int32(g.BbPos),
		Pot:     g.Pot,
	}

	for _, v := range g.Cards {
		if v != nil {
			args.Cards = append(args.Cards, &pb.CardsInfo{
				Suit:  int32(v.Suit),
				Value: int32(v.Value),
			})
		} else {
			args.Cards = append(args.Cards, nil)
		}
	}

	for _, v := range g.Players {
		if v != nil {
			player := &pb.Player{
				Id:             int32(v.Id),
				Nickname:       v.Nickname,
				Avatar:         v.Avatar,
				Pos:            int32(v.Pos),
				Bet:            int32(v.Bet),
				Win:            int32(v.Win),
				FormerChips:    int32(v.FormerChips),
				CurrentChips:   int32(v.CurrentChips),
				Action:         v.Action,
				HandLevel:      int32(v.HandLevel),
				HandFinalValue: int32(v.HandLevel),
			}
			for _, v1 := range v.Cards {
				if v1 != nil {
					player.Cards = append(player.Cards, &pb.CardsInfo{
						Suit:  int32(v1.Suit),
						Value: int32(v1.Value),
					})
				} else {
					player.Cards = append(player.Cards, nil)
				}
			}
			for _, v2 := range v.Actions {
				if v2 != nil {
					player.Actions = append(player.Actions, &pb.PlayerAction{
						Action: v2.Action,
						Bet:    int32(v2.Bet),
					})
				} else {
					player.Actions = append(player.Actions, &pb.PlayerAction{})
				}
			}
			args.Players = append(args.Players, player)
		} else {
			args.Players = append(args.Players, &pb.Player{})
		}
	}

	return args
}

func NotifyGameOver(g *models.GamblingModel) {
	args := Convert2GameTableInfoArgs(g)

	taskConn, taskServiceId := services.GetService2(SRV_NAME_TASK)
	if taskConn == nil {
		log.Error("cannot get task service:", taskServiceId)
	} else {
		taskCli := pb.NewTaskServiceClient(taskConn)
		_, err := taskCli.GameOver(context.Background(), args)
		if err != nil {
			log.Error("taskCli.GameOver: ", err)
		}
	}

	stsConn, stsServiceId := services.GetService2(SRV_NAME_STS)
	if stsConn == nil {
		log.Error("cannot get sts service:", stsServiceId)
	} else {
		stsCli := pb.NewStsServiceClient(stsConn)
		_, err := stsCli.GameInfo(context.Background(), args)
		if err != nil {
			log.Error("stsCli.GameInfo: ", err)
		}
	}

}
