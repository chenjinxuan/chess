package main

import (
	"testing"
)

func TestRankSet(t *testing.T) {
	rs := NewRankSet()
	for i := int32(0); i <= UPPER_THRESHOLD+1; i++ {
		rs.Update(i, i)
	}
	t.Log(rs.Count())

	for i := int32(0); i <= UPPER_THRESHOLD-LOWER_THRESHOLD+3; i++ {
		rs.Delete(i)
	}
	t.Log(rs.Count())

	for i := int32(0); i <= UPPER_THRESHOLD-LOWER_THRESHOLD+3; i++ {
		rs.Update(i, i)
	}
	t.Log(rs.Count())
}
