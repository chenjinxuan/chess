package texas_holdem

import (
	pb "chess/srv/srv-room/proto"
)

const (
	SUITSIZE int = 4  //四种花色
	CARDRANK int = 13 //2 3 4....K A
)

type Card struct {
	Suit  int //程序统一标准：0是黑桃、1是红桃、2是梅花、3是方片
	Value int //0代表‘牌2’、1代表‘牌3’...etc

	Showtime int //just for sort
}

func (c Card) ToProtoMessage() *pb.CardInfo {
	return &pb.CardInfo{
		Suit:  int32(c.Suit),
		Val: int32(c.Value),
	}
}

//实现sort包中的排序接口
type Cards []*Card

func (c Cards) Len() int {
	return len(c)
}

func (c Cards) Less(i, j int) bool {
	if c[i].Showtime < c[j].Showtime {
		return true
	} else if c[i].Showtime > c[j].Showtime {
		return false
	} else {
		return c[i].Value < c[j].Value
	}
}

func (c Cards) Swap(i, j int) {
	tmp := c[i]
	c[i] = c[j]
	c[j] = tmp
}

func (c Cards) ToProtoMessage() []*pb.CardInfo {
	var _cards []*pb.CardInfo
	for _, v := range c {
		_card := v.ToProtoMessage()
		_cards = append(_cards, _card)
	}
	return _cards
}
