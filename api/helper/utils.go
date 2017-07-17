package helper

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

func StringInSlice(arr []string, s string) int {
	for k, v := range arr {
		if v == s {
			return k
		}
	}
	return -1
}

func Int64InSlice(arr []int64, s int64) int {
	for k, v := range arr {
		if v == s {
			return k
		}
	}
	return -1
}

// A data structure to hold key/value pairs
type Pair struct {
	Key   int
	Value float32
}

// A slice of pairs that implements sort.Interface to sort by values
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func SortMapByValue(m map[int]float32) PairList {
	pl := make(PairList, len(m))
	i := 0
	for k, v := range m {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

// 获取流量值 M -> G
func GenFlowValue(num int) string {
	var val string

	if num < 1024 {
		val = fmt.Sprintf("%dM", num)
	} else {
		val = fmt.Sprintf("%dG", num/1024)
	}

	return val
}

func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}

func RandInt64(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Int63n(max-min)
}
