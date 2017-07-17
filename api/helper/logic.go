package helper

import (
	"fmt"
	"strconv"
	"strings"
)

func PackPoolKey(goodsid, period int) string {
	return fmt.Sprintf("%d_%d", goodsid, period)
}

func UnpackPoolKey(key string) (goodsid, period int, err error) {
	res := strings.Split(key, "_")
	if len(res) >= 2 {
		goodsid, err = strconv.Atoi(res[0])
		period, err = strconv.Atoi(res[1])
		return
	}
	return 0, 0, fmt.Errorf("could not unpack pool key(%s)", key)
}

// 模拟投注
func PackImitPoolKey(goodsid, period int) string {
	return fmt.Sprintf("imit_%d_%d", goodsid, period)
}

func UnpackImitPoolKey(key string) (goodsid, period int, err error) {
	res := strings.Split(key, "_")
	if len(res) >= 3 {
		goodsid, err = strconv.Atoi(res[1])
		period, err = strconv.Atoi(res[2])
		return
	}
	return 0, 0, fmt.Errorf("could not unpack pool key(%s)", key)
}
