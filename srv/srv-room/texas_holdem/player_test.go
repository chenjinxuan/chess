package texas_holdem_test

import (
	th "chess/srv/srv-room/texas_holdem"
	"testing"
)

func TestPlayers_ToProtoMessage(t *testing.T) {
	players := th.Players{
		&th.Player{Pos: 1, Id: 1},
		&th.Player{Pos: 2, Id: 2},
	}
	t.Log(players.ToProtoMessage())
}
