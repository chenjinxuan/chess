package helper

import (
	"code.google.com/p/mahonia"
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GenMobileNickname(mobilenum string) string {
	if len(mobilenum) == 11 {
		return mobilenum[:3] + "***" + mobilenum[8:]
	} else {
		return mobilenum
	}

}

func IsMobileNumber(mobile string) bool {
	regular := `^1(3[0-9]|4[57]|5[0-35-9]|8[0-9]|7[0-9])\d{8}$`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobile)
}

func IsUrl(str string) bool {
	//if len(str) < 4 {
	//	return false
	//}
	//if str[:4] != "http" {
	//	return false
	//}

	// regular := `\bhttps?://[a-zA-Z0-9\-.]+(?:(?:/[a-zA-Z0-9\-._?,'+\&%$=~*!():@\\]*)+)?`
	regular := `\b(https?)://(?:(\S+?)(?::(\S+?))?@)?([a-zA-Z0-9\-.]+)(?::(\d+))?((?:/[a-zA-Z0-9\-._?,'+\&%$=~*!():@\\]*)+)?`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(str)
}

func IsChineseChar(str string) bool {
	regular := `[\u4e00-\u9fa5]`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(str)
}

func IsChar(str string) bool {
	// regular := `^[a-zA-Z0-9_u4e00-u9fa5]+$`
	// regular := "^[a-zA-Z0-9_-\u4e00-\u9fa5]+$"
	regular := `^[a-z0-9A-Z-_\p{Han}]*$`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(str)
}

func IsNumeric(str string) bool {
	regular := `^[1-9][0-9]*$`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(str)
}

func ConverUnsupportStr(str string) string {
	// srtAry := strings.Split(str, "")
	srtAry := []rune(str)
	var strResult []string
	for _, v := range srtAry {
		if IsChar(string(v)) {
			strResult = append(strResult, string(v))
		}
	}
	return strings.Join(strResult, "")

}

func ConverUnsupportStrForWechat(str string) string {
	srtAry := strings.Split(str, "")
	var strResult []string
	for _, v := range srtAry {
		if !IsChar(v) {
			return "微信用户"
		}
		if IsChar(v) {
			strResult = append(strResult, v)
		}
	}
	return str
	// return strings.Join(strResult, "")

}

func ConverUnsupportStrForQQ(str string) string {
	srtAry := strings.Split(str, "")
	var strResult []string
	for _, v := range srtAry {
		if !IsChar(v) {
			return "QQ用户"
		}
		if IsChar(v) {
			strResult = append(strResult, v)
		}
	}
	return strings.Join(strResult, "")
}

func ConvertNickname(str string) string {
	maxLength := 10
	str = ConverUnsupportStr(str)
	s := []rune(str)
	if len(s) > maxLength {
		return string(s[0:maxLength])
	}
	return string(str)
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func CovertGbkToUtf8(input string) (string, error) {
	dec := mahonia.NewDecoder("gbk")

	if output, ok := dec.ConvertStringOK(input); ok {
		return ReplaceGbk(output), nil
	}

	return "", errors.New("Covert Gbk To Utf8 Fail.")
}

func ReplaceGbk(input string) string {
	reg := regexp.MustCompile(`(?:gb2312|GB2312)`)
	output := reg.ReplaceAllString(input, "utf-8")
	return output
}

func ShuffleString(s string) string {
	a := []rune(s)
	rand.Seed(time.Now().UnixNano())
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
	return string(a[:])
}

// 获取随机字符串
func RandStr(l int, format string) string {
	var chars string
	switch format {
	case "ALL":
		chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-@#~"
	case "CHAR":
		chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-@#~"
	case "NUMBER":
		chars = "0123456789"
	case "NUMBEREXPZERO":
		chars = "123456789"
	default:
		chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-@#~"
	}

	if l > len(chars) {
		return ""
	}
	new := ShuffleString(chars)
	return new[:l]
}
