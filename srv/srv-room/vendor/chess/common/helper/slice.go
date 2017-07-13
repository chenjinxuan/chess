package helper

import (
	"math/rand"
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

func ShuffleStringSlice(a []string) {
	rand.Seed(time.Now().UnixNano())
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}
