package texas_holdem_test

import (
	th "chess/srv/srv-room/texas_holdem"
	"testing"
)

func TestAnalyseCards(t *testing.T) {
	//test have not Init
	h := th.NewHand()
	err := h.SetCard(&th.Card{Suit: 2, Value: 12})
	if err.Error() != "Hand must init first" {
		t.Error("test have not Init fail")
	}

	err = h.AnalyseHand()
	if err.Error() != "Hand must init first" {
		t.Error("test have not Init fail")
	}

	//test RoyalFlush : 黑桃10 J Q K A  红桃A  梅花A
	h.Init()

	h.SetCard(&th.Card{Suit: 0, Value: 12})
	h.SetCard(&th.Card{Suit: 0, Value: 11})
	h.SetCard(&th.Card{Suit: 0, Value: 10})
	h.SetCard(&th.Card{Suit: 0, Value: 9})
	h.SetCard(&th.Card{Suit: 0, Value: 8})
	h.SetCard(&th.Card{Suit: 1, Value: 12})
	h.SetCard(&th.Card{Suit: 2, Value: 12})

	h.AnalyseHand()

	if h.Level != 10 || h.FinalValue != -1 {
		t.Error("test RoyalFlush fail")
	}

	//test 8 cards
	err = h.SetCard(&th.Card{Suit: 2, Value: 12})
	if err.Error() != "after a game, you should init Hand again" {
		t.Error("test 8 cards fail")
	}

	//test straight flush
	h.Init()
	h.SetCard(&th.Card{Suit: 1, Value: 12})
	h.SetCard(&th.Card{Suit: 0, Value: 11})
	h.SetCard(&th.Card{Suit: 0, Value: 10})
	h.SetCard(&th.Card{Suit: 0, Value: 9})
	h.SetCard(&th.Card{Suit: 0, Value: 8})
	h.SetCard(&th.Card{Suit: 0, Value: 7})
	h.SetCard(&th.Card{Suit: 2, Value: 12})

	h.AnalyseHand()

	if h.Level != 9 || h.FinalValue != 13 {
		t.Error("test straight flush fail")
	}

	//test four of a kind
	h.Init()
	h.SetCard(&th.Card{Suit: 0, Value: 2})
	h.SetCard(&th.Card{Suit: 1, Value: 2})
	h.SetCard(&th.Card{Suit: 2, Value: 2})
	h.SetCard(&th.Card{Suit: 3, Value: 2})
	h.SetCard(&th.Card{Suit: 0, Value: 12})
	h.SetCard(&th.Card{Suit: 1, Value: 12})
	h.SetCard(&th.Card{Suit: 2, Value: 12})

	h.AnalyseHand()

	if h.Level != 8 || h.FinalValue != 2223332 {
		t.Error("test four of a kind fail")
	}

	//test full house
	h.Init()
	h.SetCard(&th.Card{Suit: 0, Value: 4})
	h.SetCard(&th.Card{Suit: 1, Value: 3})
	h.SetCard(&th.Card{Suit: 2, Value: 2})
	h.SetCard(&th.Card{Suit: 3, Value: 2})
	h.SetCard(&th.Card{Suit: 0, Value: 12})
	h.SetCard(&th.Card{Suit: 1, Value: 12})
	h.SetCard(&th.Card{Suit: 2, Value: 12})

	h.AnalyseHand()

	if h.Level != 7 || h.FinalValue != 13322243 {
		t.Error("test full house fail")
	}

	//test flush
	h.Init()
	h.SetCard(&th.Card{Suit: 0, Value: 4})
	h.SetCard(&th.Card{Suit: 0, Value: 3})
	h.SetCard(&th.Card{Suit: 0, Value: 2})
	h.SetCard(&th.Card{Suit: 0, Value: 5})
	h.SetCard(&th.Card{Suit: 0, Value: 7})
	h.SetCard(&th.Card{Suit: 1, Value: 12})
	h.SetCard(&th.Card{Suit: 2, Value: 6})

	h.AnalyseHand()

	if h.Level != 6 || h.FinalValue != 376 {
		t.Error("test flush fail")
	}

	//test straight
	h.Init()
	h.SetCard(&th.Card{Suit: 0, Value: 2})
	h.SetCard(&th.Card{Suit: 1, Value: 3})
	h.SetCard(&th.Card{Suit: 2, Value: 4})
	h.SetCard(&th.Card{Suit: 3, Value: 5})
	h.SetCard(&th.Card{Suit: 0, Value: 6})
	h.SetCard(&th.Card{Suit: 1, Value: 6})
	h.SetCard(&th.Card{Suit: 2, Value: 6})

	h.AnalyseHand()

	if h.Level != 5 || h.FinalValue != 8 {
		t.Error("test straight fail")
	}

	//test three of a kind
	h.Init()
	h.SetCard(&th.Card{Suit: 0, Value: 2})
	h.SetCard(&th.Card{Suit: 1, Value: 3})
	h.SetCard(&th.Card{Suit: 2, Value: 4})
	h.SetCard(&th.Card{Suit: 3, Value: 7})
	h.SetCard(&th.Card{Suit: 0, Value: 6})
	h.SetCard(&th.Card{Suit: 1, Value: 6})
	h.SetCard(&th.Card{Suit: 2, Value: 6})

	h.AnalyseHand()

	if h.Level != 4 || h.FinalValue != 6667432 {
		t.Error("test three of a kind fail")
	}

	//test two pairs
	h.Init()
	h.SetCard(&th.Card{Suit: 0, Value: 2})
	h.SetCard(&th.Card{Suit: 1, Value: 2})
	h.SetCard(&th.Card{Suit: 2, Value: 4})
	h.SetCard(&th.Card{Suit: 3, Value: 4})
	h.SetCard(&th.Card{Suit: 0, Value: 6})
	h.SetCard(&th.Card{Suit: 1, Value: 7})
	h.SetCard(&th.Card{Suit: 2, Value: 8})

	h.AnalyseHand()

	if h.Level != 3 || h.FinalValue != 4422876 {
		t.Error("test two pairs fail")
	}

	//test one pair
	h.Init()
	h.SetCard(&th.Card{Suit: 0, Value: 2})
	h.SetCard(&th.Card{Suit: 1, Value: 2})
	h.SetCard(&th.Card{Suit: 2, Value: 4})
	h.SetCard(&th.Card{Suit: 3, Value: 6})
	h.SetCard(&th.Card{Suit: 0, Value: 7})
	h.SetCard(&th.Card{Suit: 1, Value: 8})
	h.SetCard(&th.Card{Suit: 2, Value: 10})

	h.AnalyseHand()

	if h.Level != 2 || h.FinalValue != 2308764 {
		t.Error("test one pair fail")
	}

	//test high card
	h.Init()
	h.SetCard(&th.Card{Suit: 0, Value: 2})
	h.SetCard(&th.Card{Suit: 1, Value: 3})
	h.SetCard(&th.Card{Suit: 2, Value: 4})
	h.SetCard(&th.Card{Suit: 3, Value: 6})
	h.SetCard(&th.Card{Suit: 0, Value: 9})
	h.SetCard(&th.Card{Suit: 1, Value: 11})
	h.SetCard(&th.Card{Suit: 2, Value: 12})

	h.AnalyseHand()

	if h.Level != 1 || h.FinalValue != 13196432 {
		t.Error("test high card fail")
	}
}
