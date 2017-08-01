package auth

import (
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

const (
	//  /profile/password
	Password_OK       = 1
	Length_Too_Short  = -32
	Length_Too_Long   = -33
	All_Number        = -34
	All_Letter        = -35
	Include_Blank     = -36
	Not_Strong_Enough = -37
	INCLUDE_BAD_CHAR  = -38
)

var Passwords *Password

type Password struct {
}

func (p *Password) Hash(password string) (string, error) {
	pass := []byte(password)
	hp, err := bcrypt.GenerateFromPassword(pass, 0)
	return string(hp), err
}

func (p *Password) Check(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (p *Password) StrongCheck(password string) bool {
	regular := `^[a-zA-Z0-9]{6,20}$`
	//regular := `^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,20}$`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(password)
}

func (p *Password) IsAllLetter(str string) bool {
	regular := `^([a-zA-Z]+)$`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(str)
}

func (p *Password) IsAllNumber(str string) bool {
	regular := `^[1-9]\d*$`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(str)
}

func (p *Password) NoBlank(str string) bool {
	//regular := `^[^\s]*＄`
	regular := `^[^\s]{6,20}$`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(str)
}

// 判断密码强度,1为符合要求
func (p *Password) CheckPasswordStrong(password string) int {
	if len(password) < 6 {
		return Length_Too_Short
	}
	if len(password) > 20 {
		return Length_Too_Long
	}
	if p.IsAllNumber(password) {
		return All_Number
	}
	if p.IsAllLetter(password) {
		return All_Letter
	}
	if !p.NoBlank(password) {
		return Include_Blank
	}
	return Password_OK
}
