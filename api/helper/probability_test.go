package helper

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestGetRand(t *testing.T) {
	proArr := make(map[string]int)
	proArrTmp := make(map[string]int)
	proArr["111"] = 10
	proArr["222"] = 10
	proArr["333"] = 10
	proArr["444"] = 10
	proArr["555"] = 10
	proArr["666"] = 10
	pro := NewProbability(proArr)

	for i := 0; i < 1000; i++ {
		fmt.Println((rand.Intn(100000000) + 1))
		proArrTmp[pro.GetRand()]++
		//fmt.Println(pro.GetRand())
	}

	fmt.Println(proArrTmp)
}
