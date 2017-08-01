package models

import (
	"database/sql"
	"github.com/Sirupsen/logrus"
	"time"
	"chess/api/log"
)

const (
	TYPE_DEFAULT            = 0
	TYPE_MOBILE_QUICK_LOGIN = 1
	TYPE_MOBILE_REG         = 2
	TYPE_QQ                 = 3
	TYPE_WEIBO              = 4
	TYPE_WECHAT             = 5

	IS_FRESH_FALSE = 0 // 走完新手引导
	IS_FRESH_TRUE  = 1
	IS_FRESH_JUMP  = 2 // 跳过
)

var Users = new(UsersModel)

type UsersModel struct {
	Id            int
	Email         string
	Pwd           string
	Nickname      string
	MobileNumber  string
	ContactMobile string
	RegIp         string
	Gender        int
	Avatar        string
	LastLoginIp   string
	Channel       string
	Type          int
	AppFrom          string
	Status        int
	IsFresh       int
	Updated       time.Time
	Created       time.Time
}

func (m *UsersModel) Get(id int, user *UsersModel) error {
	sqlString := `SELECT
					id, email, pwd, nickname,
					mobile_number,contact_mobile, gender,avatar, reg_ip,
					last_login_ip, channel,type,status, is_fresh,updated, created
				FROM users WHERE id = ?`
	return Mysql.Chess.QueryRow(
		sqlString, id,
	).Scan(
		&user.Id,
		&user.Email,
		&user.Pwd,
		&user.Nickname,
		&user.MobileNumber,
		&user.ContactMobile,
		&user.Gender,
		&user.Avatar,
		&user.RegIp,
		&user.LastLoginIp,
		&user.Channel,
		&user.Type,
		&user.Status,
		&user.IsFresh,
		&user.Updated,
		&user.Created,
	)
}

func (m *UsersModel) GetByMobileNumber(mobileNumber string, user *UsersModel) error {
	sqlString := `SELECT
					id, email, pwd, nickname,
					mobile_number, reg_ip,
					last_login_ip, status,is_fresh,
					updated, created
				FROM users WHERE mobile_number = ?`
	return Mysql.Chess.QueryRow(sqlString, mobileNumber).Scan(
		&user.Id,
		&user.Email,
		&user.Pwd,
		&user.Nickname,
		&user.MobileNumber,
		&user.RegIp,
		&user.LastLoginIp,
		&user.Status,
		&user.IsFresh,
		&user.Updated,
		&user.Created,
	)
}

func (m *UsersModel) GetContactMobileById(id int) (mobileNumber string, err error) {
	sqlStr := `SELECT contact_mobile FROM users WHERE id = ?`

	err = Mysql.Chess.QueryRow(sqlStr, id).Scan(&mobileNumber)
	return
}

func (m *UsersModel) Insert(user *UsersModel) (int, error) {
	sqlString := `INSERT INTO users
		(email, pwd, nickname, mobile_number,gender,avatar, reg_ip, last_login_ip,channel,app_from,type, status)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?,?,?,?,?)`
	result, err := Mysql.Chess.Exec(
		sqlString,
		user.Email,
		user.Pwd,
		user.Nickname,
		user.MobileNumber,
		user.Gender,
		user.Avatar,
		user.RegIp,
		user.LastLoginIp,
		user.Channel,
	        user.AppFrom,
		user.Type,
		1,
	)

	// Debug
	log.Log.WithFields(logrus.Fields{
		"sql":   sqlString,
		"Error": err,
	}).Debug("UsersModel.Insert")

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	return int(id), err
}

func (m *UsersModel) UpdateFresh(id int, is_fresh int) (err error) {
	sqlStr := `UPDATE users SET is_fresh = ? WHERE id = ?`
	_, err = Mysql.Chess.Exec(sqlStr, is_fresh, id)
	return err
}

func (m *UsersModel) GetDetail(id int, user *UsersModel) error {
	sqlString := `SELECT
					u.id, u.email, u.pwd, u.nickname,
					u.mobile_number,u.contact_mobile, u.gender, u.avatar, u.reg_ip,
					u.last_login_ip, u.channel, u.type, u.status, u.is_fresh, IFNULL(ui.device_from, ''), u.updated, u.created
				FROM users AS u LEFT JOIN users_info AS ui ON u.id = ui.user_id
				WHERE u.id = ?`
	return Mysql.Chess.QueryRow(
		sqlString, id,
	).Scan(
		&user.Id,
		&user.Email,
		&user.Pwd,
		&user.Nickname,
		&user.MobileNumber,
		&user.ContactMobile,
		&user.Gender,
		&user.Avatar,
		&user.RegIp,
		&user.LastLoginIp,
		&user.Channel,
		&user.Type,
		&user.Status,
		&user.IsFresh,
		&user.AppFrom,
		&user.Updated,
		&user.Created,
	)
}

func (m *UsersModel) GetPwdByMobile(mobile string) (pwd string, err error) {
	sqlStr := `SELECT pwd FROM users WHERE mobile_number = ?`
	err = Mysql.Chess.QueryRow(sqlStr, mobile).Scan(&pwd)
	return
}

func (m *UsersModel) MergeUser(pwd string, mobile string) error {
	tx, err := Mysql.Chess.Begin()
	if err != nil {
		return err
	}

	// 删除手机号
	sqlStr := `UPDATE users
		SET mobile_number = ?,contact_mobile = ? 
		WHERE mobile_number = ?`
	newMobile := mobile + "#"
	res, err := tx.Exec(sqlStr, newMobile, newMobile, mobile)
	log.Log.Debug(sqlStr)
	if err != nil {
		tx.Rollback()
		return err
	} else {
		affected, _ := res.RowsAffected()
		if affected != 1 {
			tx.Rollback()
			return sql.ErrNoRows
		}
	}

	sqlStr = `UPDATE users SET pwd = ? WHERE id = ? AND pwd = " "`
	res, err = tx.Exec(sqlStr, pwd, m.Id)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return err
	}

	sqlStr = `UPDATE users SET
		 mobile_number = ?,contact_mobile = ?
		WHERE id = ?`
	res, err = tx.Exec(sqlStr, mobile, mobile, m.Id)
	if err != nil {
		tx.Rollback()
		return err
	} else {
		affected, _ := res.RowsAffected()
		if affected != 1 {
			tx.Rollback()
			return sql.ErrNoRows
		}
	}
	return tx.Commit()
}

// 合并余额
func (m *UsersModel) MergeUserNew(uid, oldUid int, balance uint, pwd string, mobile string) error {
	tx, err := Mysql.Chess.Begin()
	if err != nil {
		return err
	}

	// 删除手机号
	sqlStr := `UPDATE users
		SET mobile_number = ?,contact_mobile = ?
		WHERE mobile_number = ?`
	newMobile := mobile + "#"
	res, err := tx.Exec(sqlStr, newMobile, newMobile, mobile)
	log.Log.Debug(sqlStr)
	if err != nil {
		tx.Rollback()
		return err
	} else {
		affected, _ := res.RowsAffected()
		if affected != 1 {
			tx.Rollback()
			return sql.ErrNoRows
		}
	}

	sqlStr = `UPDATE users SET pwd = ? WHERE id = ? AND pwd = " "`
	res, err = tx.Exec(sqlStr, pwd, uid)
	if err != nil {
		tx.Rollback()
		return err
	}

	sqlStr = `UPDATE users SET
		 mobile_number = ?,contact_mobile = ?
		WHERE id = ?`
	res, err = tx.Exec(sqlStr, mobile, mobile, uid)
	if err != nil {
		tx.Rollback()
		return err
	} else {
		affected, _ := res.RowsAffected()
		if affected != 1 {
			tx.Rollback()
			return sql.ErrNoRows
		}
	}

	sqlStr = `UPDATE users_wallet SET balance = balance - ? WHERE user_id = ? AND balance >= ?`
	_, err = tx.Exec(sqlStr, balance, oldUid, balance)
	if err != nil {
		tx.Rollback()
		return err
	}

	sqlStr = `INSERT INTO users_addcash_log (user_id, amount, tag, comment)
		VALUES
		(?, ?, ?, ?)`
	_, err = tx.Exec(sqlStr, uid, balance, "user merge", "user merge")
	if err != nil {
		tx.Rollback()
		return err
	}

	sqlStr = `UPDATE users_wallet SET balance = balance + ? WHERE user_id = ?`
	_, err = tx.Exec(sqlStr, balance, uid)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (m *UsersModel) UpdateLastLoginIp(uid int, ip string) error {
	sqlString := `UPDATE users
		SET last_login_ip = ?
		WHERE id = ?`
	_, err := Mysql.Chess.Exec(sqlString, ip, uid)
	return err

}

// 随机取一个ai用户
func (m *UsersModel) GetSimulateUserByRand() (user UsersModel, err error) {
	sqlStr := `SELECT id, nickname, avatar, reg_ip
		FROM users
		WHERE id = (SELECT user_id FROM simulate_charge ORDER BY id DESC LIMIT 1)`

	err = Mysql.Chess.QueryRow(sqlStr).Scan(
		&user.Id,
		&user.Nickname,
		&user.Avatar,
		&user.RegIp,
	)

	return
}

func (m *UsersModel) GetSimulateUser(limit int) (users []UsersModel, err error) {
	sqlStr := `SELECT id, nickname, avatar, reg_ip
		FROM users
		WHERE channel = 'simulate' AND status = 1
		LIMIT ?`

	rows, err := Mysql.Chess.Query(sqlStr, limit)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user UsersModel
		err = rows.Scan(&user.Id, &user.Nickname, &user.Avatar, &user.RegIp)
		if err != nil {
			continue
		}

		users = append(users, user)
	}

	return
}

func (m *UsersModel) GetUserCreatedDays(userId int) (days int, err error) {
	sqlStr := `SELECT DATEDIFF(CURRENT_DATE, DATE(created)) + 1
		FROM users WHERE id = ?`

	err = Mysql.Chess.QueryRow(sqlStr, userId).Scan(&days)
	return
}

type UserNickModel struct {
	NickName string `json:"nick_name"`
	Number   int64  `json:"number"`
}

func (m *UsersModel) GetUserNickAuto(num int) (n []UserNickModel, err error) {
	sqlStr := `SELECT nickname
FROM users  AS t1 JOIN (SELECT ROUND(RAND() * ((SELECT MAX(id) FROM users )-(SELECT MIN(id) FROM users ))+(SELECT MIN(id) FROM users )) AS id) AS t2
WHERE t1.id >= t2.id
ORDER BY t1.id LIMIT ?`
	rows, err := Mysql.Chess.Query(sqlStr, num)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u UserNickModel
		err = rows.Scan(&u.NickName)
		if err != nil {
			continue
		}
		n = append(n, u)
	}
	return
}

func (m *UsersModel) Save() error {
    sqlString := `UPDATE users SET
		email = ?, pwd = ?, nickname = ?, mobile_number = ?, reg_ip = ?, last_login_ip = ?, status = ?
		WHERE id = ?`
    result, err := Mysql.Chess.Exec(
	sqlString,
	m.Email,
	m.Pwd,
	m.Nickname,
	m.MobileNumber,
	m.RegIp,
	m.LastLoginIp,
	m.Status,
	m.Id,
    )

    // Debug
    log.Log.WithFields(logrus.Fields{
	"sql":   sqlString,
	"Error": err,
    }).Debug("UsersModel.Save", result)

    if err != nil {
	return err
    }
    return err
}

// Update Nickname
func (m *UsersModel) UpdateNickname(nickname string) error {
    sqlString := `UPDATE users SET
		 nickname = ?
		WHERE id = ?`
    result, err := Mysql.Chess.Exec(
	sqlString,
	nickname,
	m.Id,
    )

    // Debug
    log.Log.WithFields(logrus.Fields{
	"sql":   sqlString,
	"Error": err,
    }).Debug("UsersModel.UpdateProfile", result)

    if err != nil {
	return err
    }
    return err
}

// Update Mobile
func (m *UsersModel) UpdateMobile(mobile string) error {
    sqlString := `UPDATE users SET
		 mobile_number = ?
		WHERE id = ?`
    result, err := Mysql.Chess.Exec(
	sqlString,
	mobile,
	m.Id,
    )

    // Debug
    log.Log.WithFields(logrus.Fields{
	"sql":   sqlString,
	"Error": err,
    }).Debug("UsersModel.UpdateProfile", result)

    if err != nil {
	return err
    }
    return err
}

func (m *UsersModel) UpdateAllMobile(mobile string) error {
    sqlString := `UPDATE users 
		SET mobile_number = ?,contact_mobile = ? 
		WHERE id = ?`
    result, err := Mysql.Chess.Exec(
	sqlString,
	mobile,
	mobile,
	m.Id,
    )

    // Debug
    log.Log.WithFields(logrus.Fields{
	"sql":   sqlString,
	"Error": err,
    }).Debug("UsersModel.UpdateProfile", result)

    if err != nil {
	return err
    }
    return err
}

// Update Password
func (m *UsersModel) UpdatePassword(password string) error {
    sqlString := `UPDATE users SET
		 pwd = ?
		WHERE id = ?`
    result, err := Mysql.Chess.Exec(
	sqlString,
	password,
	m.Id,
    )

    // Debug
    log.Log.WithFields(logrus.Fields{
	"sql":   sqlString,
	"Error": err,
    }).Debug("UsersModel.UpdateProfile.Password", result)

    if err != nil {
	return err
    }
    return err
}

// Update Contact Mobile
func (m *UsersModel) UpdateContactMobile(mobile string) error {
    sqlString := `UPDATE users SET
		 contact_mobile = ?
		WHERE id = ?`
    result, err := Mysql.Chess.Exec(
	sqlString,
	mobile,
	m.Id,
    )

    // Debug
    log.Log.WithFields(logrus.Fields{
	"sql":   sqlString,
	"Error": err,
    }).Debug("UsersModel.UpdateContactMobile", result)

    if err != nil {
	return err
    }
    return err
}

// Update Avatar
func (m *UsersModel) UpdateAvatar(avatar string) error {
    sqlString := `UPDATE users SET
		 avatar = ?
		WHERE id = ?`
    result, err := Mysql.Chess.Exec(
	sqlString,
	avatar,
	m.Id,
    )

    // Debug
    log.Log.WithFields(logrus.Fields{
	"sql":   sqlString,
	"Error": err,
    }).Debug("UsersModel.UpdateAvatar", result)

    if err != nil {
	return err
    }
    return err
}

// Update Avatar
func (m *UsersModel) UpdateGender(gender int) error {
    sqlString := `UPDATE users SET
		 gender = ?
		WHERE id = ?`
    result, err := Mysql.Chess.Exec(
	sqlString,
	gender,
	m.Id,
    )

    // Debug
    log.Log.WithFields(logrus.Fields{
	"sql":   sqlString,
	"Error": err,
    }).Debug("UsersModel.UpdateGender", result)

    if err != nil {
	return err
    }
    return err
}