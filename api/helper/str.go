package helper

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func GenMobileNickname(mobilenum string) string {
	return mobilenum[:3] + "***" + mobilenum[8:]
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

func IsCharOK(str string) bool {
	// regular := `^[a-zA-Z0-9_u4e00-u9fa5]+$`
	// regular := "^[a-zA-Z0-9_-\u4e00-\u9fa5]+$"
	regular := `^[a-z0-9A-Z-_\p{Han}]*$`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(str)
}

func ConverUnsupportStr(str string) string {
	// srtAry := strings.Split(str, "")
	//srtAry := []rune(str)
	//var strResult []string
	//for _, v := range srtAry {
	//	if IsCharOK(string(v)) {
	//		strResult = append(strResult, string(v))
	//	}
	//}
	return strings.Trim(str, " ")

}

func ConverUnsupportStrForWechat(str string) string {
	srtAry := strings.Split(str, "")
	var strResult []string
	for _, v := range srtAry {
		if !IsCharOK(v) {
			return "微信用户"
		}
		if IsCharOK(v) {
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
		if !IsCharOK(v) {
			return "QQ用户"
		}
		if IsCharOK(v) {
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

// 过滤 emoji 表情
func FilterEmoji(content string) string {
	new_content := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			new_content += string(value)
		}
	}
	return new_content
}

//手机号验证
func FilterMobile(mobile string) bool {
	regular := `^1[3|4|5|7|8][0-9]{9}$`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobile)
}
