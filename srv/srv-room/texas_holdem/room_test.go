package texas_holdem_test

import (
	th "chess/srv/srv-room/texas_holdem"
	"testing"
)

func TestRoom_GetTable(t *testing.T) {
	th.InitRoomList()

	t1 := th.GetTable(1, "")
	t.Log(t1)

	t2 := th.GetTable(1, "")
	t.Log(t2)
}
