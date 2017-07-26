package models

import (
	"time"
)

var UsersRegisterLog = new(UsersRegisterLogModel)

const (
	UserInitStatusFresh    = 1
	UserInitStatusNotFresh = 2
)

type UsersRegisterLogModel struct {
	Id              int
	UserId          int
	InitStatus      int
	InitSnapshot    string
	Channel         string
	From            string
	Ver             string
	DeviceUniqueKey string
	DeviceUniqueVal string
	Updated         time.Time
	Created         time.Time
}

func (m *UsersRegisterLogModel) GetLog(key, val string, userId int) (res []UsersRegisterLogModel, err error) {
	sqlStr := `SELECT id, user_id, init_status, init_snapshot, channel, device_from, ver, device_unique_key, device_unique_val 
				FROM users_register_log
				WHERE (device_unique_key = ? and device_unique_val = ?) or user_id = ?`
	rows, err := ChessMysql.Main.Query(sqlStr, key, val, userId)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		item := UsersRegisterLogModel{}
		err = rows.Scan(
			&item.Id,
			&item.UserId,
			&item.InitStatus,
			&item.InitSnapshot,
			&item.Channel,
			&item.From,
			&item.Ver,
			&item.DeviceUniqueKey,
			&item.DeviceUniqueVal,
		)

		if err != nil {
			return
		}
		res = append(res, item)
	}

	return
}

func (m *UsersRegisterLogModel) Insert() error {
	sqlStr := `INSERT INTO users_register_log (user_id,init_status,init_snapshot,channel ,device_from,ver,device_unique_key, device_unique_val)
				VALUES (?,?,?,?,?,?,?,?) `
	_, err := ChessMysql.Main.Exec(sqlStr, m.UserId, m.InitStatus, m.InitSnapshot, m.Channel, m.From, m.Ver, m.DeviceUniqueKey, m.DeviceUniqueVal)
	return err
}
