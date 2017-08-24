package profile

import (
	"strings"
	"testing"
)

func TestNickname(t *testing.T) {
	nicknameStr := "江泽民,囧囧囧,exo,吴世勋"
	nickname := strings.Split(nicknameStr, ",")
	for _, v := range nickname {
		str, isOK, err := CheckNickname(v)
		if err != nil {
			t.Error(err)
		}
		if !isOK {
			t.Log("有敏感词:", str)
		} else {
			t.Log("正常:", str)
		}
	}
}
