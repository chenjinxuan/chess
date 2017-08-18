package texas_holdem_test

import (
	th "chess/srv/srv-room/texas_holdem"
	"testing"
)

func TestTable_ToProtoMessage(t *testing.T) {
	table := th.GetTable(1, "")
	t.Log(table.ToProtoMessage())
}
