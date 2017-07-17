package auth

import (
	"strings"
	"testing"
)

func TestPasswordStrong(t *testing.T) {
	pwdStr := "aaa,aaaaaaaaa,1111aaaa1111aaaaa1111111111111,213a4234&,234234234asdf,1234234234,abcded,you  balnk,google123"
	pwd := strings.Split(pwdStr, ",")
	for _, v := range pwd {
		ret := Passwords.CheckPasswordStrong(v)
		if ret != 1 {
			t.Log("check fail:", v, "  Ret", ret)
		} else {
			t.Log("pass: ", v)
		}
	}
}
