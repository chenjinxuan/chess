package helper

import (
	"fmt"
	"strings"
	"testing"
)

func TestGenmobileStr(t *testing.T) {
	str := "13094816493"
	res := GenMobileNickname(str)
	t.Log(res)
}

func TestUrl(t *testing.T) {
	urlStr := "http://7xo4rm.com1.z0.glb.clouddn.com/10000021_1447144128.jpg"
	if !IsUrl(urlStr) {
		t.Error("check wrong")
	}
	fileStr := "10000021_1447144128.jpg"
	if IsUrl(fileStr) {
		t.Error("file check wrong")
	}
}

func TestConverUnsupportChar(t *testing.T) {
	str := "囧测试,!#2测试,abce_-1234"
	strAry := strings.Split(str, ",")
	for _, v := range strAry {
		t.Log(ConverUnsupportStr(v))
	}
}

func TestConvertChinese(t *testing.T) {
	// str := "漳州市品迪斯特贸易"
	str := "测测强强強強烈餓死好像很好吃？在"
	t.Log(str)
	//if len(str) > 30 {
	//	str = str[:30]
	// }
	res := ConvertNickname(str)
	t.Log(res)
}

func TestConvertCnStr(t *testing.T) {
	str := "??？ok测试？在"
	t.Log(str)
	res := ConvertNickname(str)
	t.Log(res)

	str2 := "测测强强"
	t.Log(str2)
	res2 := ConvertNickname(str2)
	t.Log(res2)
}

func TestShuffleString(t *testing.T) {
	fmt.Println(ShuffleString("12345678"))

	i := 0
	for i < 10 {
		chars := "0123456789"
		// 字符串乱序
		tmpchars := ShuffleString(chars)
		fmt.Println(ShuffleString(tmpchars[:3] + tmpchars[:1]))
		i++
	}
}

func TestIsNumeric(t *testing.T) {
	t.Log(IsNumeric("23032222934332"))
}

func TestRandStr(t *testing.T) {
	i := 0
	for i < 10 {
		t.Log(RandStr(6, "NUMBER"))
		i++
	}
}
