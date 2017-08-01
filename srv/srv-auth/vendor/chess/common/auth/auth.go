package auth

import (
	"github.com/satori/go.uuid"
	"strconv"
	"time"
	"chess/common/config"
	"chess/common/define"
	"chess/models"
)

type AuthResult struct {
	define.BaseResult
	UserId       int    `json:"user_id"`
	Token        string `json:"token"`
	Expire       int64  `json:"expire"`
	RefreshToken string `json:"refresh_token"`
	// ExpireTime   int `json:"expire_time"`
}

func LoginUser(userid int, from, uniqueId string) (AuthResult, error) {
	var result AuthResult
	expire := time.Now().Add(time.Second * time.Duration(config.CAuth.Login.TokenExpire)).Unix()
	tokenString, err := CreateLoginToken(strconv.Itoa(userid), expire, define.JwtSecret)
	if err != nil {
		result.Msg = "Could not generate token."
		return result, err
	}

	// Create refresh token
	u := uuid.NewV4()
	refreshToken := u.String()

	// Create session and insert to database
	session := new(models.SessionModel)
	session.UserId = userid
	session.From = from
	session.UniqueId = uniqueId
	session.Token = &models.SessionToken{tokenString, expire}
	session.RefreshToken = refreshToken
	session.Updated = time.Now()
	session.Created = time.Now()
	err = models.Session.Upsert(userid, from, uniqueId, session)
	if err != nil {
		result.Msg = "Could not generate session."
		return result, err
	}

	result.Ret = 1
	result.UserId = userid
	result.Token = tokenString
	result.Expire = expire
	result.RefreshToken = refreshToken
	return result, nil
}

// verify login with password
func VerifyLogin(loginType uint8, user *models.UsersModel, password string) bool {
	switch loginType {
	case 0: // mobile number login
		err := models.Users.GetByMobileNumber(user.MobileNumber, user)
		if err != nil {
			return false
		}

		err = Passwords.Check(user.Pwd, password)
		if err != nil {
			return false
		}

		return true

	default:
	}

	return false
}
