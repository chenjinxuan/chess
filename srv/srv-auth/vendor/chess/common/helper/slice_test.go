package helper

import (
	"fmt"
	"testing"
)

func TestShuffle(t *testing.T) {
	s := []string{"1", "2", "3", "4"}

	fmt.Println(s)
	ShuffleStringSlice(s)
	fmt.Println(s)
}
