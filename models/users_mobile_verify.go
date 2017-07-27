package models

import (
	"github.com/Sirupsen/logrus"
	"chess/common/log"
)

var UsersMobileVerify = new(UsersMobileVerifyModel)

type UsersMobileVerifyModel struct {
	Id           int
	MobileNumber string
	VerifyCode   string
	VerifyType   int
	Expire       int
	Status       int
}

func (m *UsersMobileVerifyModel) Insert(mobileNumber, verifyCode, ip string, verifyType, expire int) (int, error) {
	sqlString := `INSERT INTO users_mobile_verify (mobile_number, verify_type, verify_code,ip, expire, status,send_status) VALUES (?, ?, ?, ?, ?, ?,?)`
	result, err := Mysql.Chess.Exec(
		sqlString,
		mobileNumber,
		verifyType,
		verifyCode,
		ip,
		expire,
		0,
		0,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *UsersMobileVerifyModel) GetCode(userMobileVerify *UsersMobileVerifyModel) error {
	sqlString := `SELECT id, verify_code, expire, status FROM users_mobile_verify WHERE mobile_number = ? AND verify_type = ? ORDER BY id DESC`
	err := Mysql.Chess.QueryRow(
		sqlString,
		userMobileVerify.MobileNumber,
		userMobileVerify.VerifyType,
	).Scan(
		&userMobileVerify.Id,
		&userMobileVerify.VerifyCode,
		&userMobileVerify.Expire,
		&userMobileVerify.Status,
	)

	log.Debugf("UsersMobileVerifyModel.GetCode",logrus.Fields{
		"sql":   sqlString,
		"Error": err,
		"v":     userMobileVerify,
	})

	return err
}

func (m *UsersMobileVerifyModel) GetCodeVerified(userMobileVerify *UsersMobileVerifyModel) error {
	sqlString := `SELECT id, verify_code, expire,status FROM users_mobile_verify WHERE mobile_number = ? AND verify_type = ?  AND status = 1 ORDER BY id DESC`
	err := Mysql.Chess.QueryRow(
		sqlString,
		userMobileVerify.MobileNumber,
		userMobileVerify.VerifyType,
	).Scan(
		&userMobileVerify.Id,
		&userMobileVerify.VerifyCode,
		&userMobileVerify.Expire,
		&userMobileVerify.Status,
	)
	log.Debugf("UsersMobileVerifyModel.GetCode",logrus.Fields{
		"sql":   sqlString,
		"Error": err,
		"v":     userMobileVerify,
	})

	return err
}
func (m *UsersMobileVerifyModel) SetCodeStatus(id, status int) error {
	sqlString := `UPDATE users_mobile_verify SET status = ? WHERE id = ?`
	_, err := Mysql.Chess.Exec(
		sqlString,
		status,
		id,
	)
	return err
}

func (m *UsersMobileVerifyModel) SetSendStatus(id, status int) error {
	sqlString := `UPDATE users_mobile_verify SET send_status = ? WHERE id = ?`
	_, err := Mysql.Chess.Exec(
		sqlString,
		status,
		id,
	)
	return err
}

func (m *UsersMobileVerifyModel) CountRecentlySend(mobile string, verifyType, hours int) (cnt int, err error) {
	sqlStr := `SELECT COUNT(1)
		FROM users_mobile_verify
		WHERE mobile_number = ? AND verify_type = ? AND status = 1 AND created >= DATE_SUB(NOW(), INTERVAL ? HOUR)`

	err = Mysql.Chess.QueryRow(sqlStr, mobile, verifyType, hours).Scan(&cnt)
	return
}
