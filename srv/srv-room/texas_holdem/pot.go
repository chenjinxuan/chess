package texas_holdem

import (
	pb "chess/srv/srv-room/proto"
	"sort"
)

type handBet struct {
	Pos int
	Bet int
}

type handBets []handBet

func (p handBets) Len() int {
	return len(p)
}

func (p handBets) Less(i, j int) bool {
	return p[i].Bet < p[j].Bet
}

func (p handBets) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type handPot struct {
	Pot  int   // 池子筹码
	OPos []int // 参与池子的玩家pos
}

// main pot and side-pot calculation
func calcPot(bets []int32) (pots []handPot) {
	var obs []handBet

	//fmt.Println("amount bets:", bets)
	for i, bet := range bets {
		if bet > 0 {
			obs = append(obs, handBet{Pos: i + 1, Bet: int(bet)})
		}
	}
	sort.Sort(handBets(obs))

	//fmt.Println("amount bets(sorted):", obs)

	for i, ob := range obs {
		if ob.Bet > 0 {
			s := obs[i:]
			hpot := handPot{Pot: ob.Bet * len(s)}

			for j := range s {
				s[j].Bet -= ob.Bet
				hpot.OPos = append(hpot.OPos, s[j].Pos)
			}
			pots = append(pots, hpot)
		}
	}
	return
}

// 池子计算结果
type PotDetail struct {
	Pot int
	Ps  []int32
}

func (c PotDetail) ToProtoMessage() *pb.PotInfo {
	return &pb.PotInfo{
		Pot: int32(c.Pot),
		Ps:  c.Ps,
	}
}

type PotDetails []*PotDetail

func (pds PotDetails) ToProtoMessage() []*pb.PotInfo {
	_pds := make([]*pb.PotInfo, len(pds))
	for k, v := range pds {
		if v != nil {
			tmp := *v
			_pd := tmp.ToProtoMessage()
			_pds[k] = _pd
		} else {
			_pds[k] = &pb.PotInfo{}
		}
	}
	return _pds
}
