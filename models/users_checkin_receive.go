package models

import "time"

type UsersCheckinReceiveModel struct {
	Id                  int
	UserId              int
	CheckinDaysRewardId int
	IsMore              int
	Updated             time.Time
	Created             time.Time
}

var UsersCheckinReceive = new(UsersCheckinReceiveModel)

func (m *UsersCheckinReceiveModel) GetIsMore(userId int) (isMore int, err error) {
	sqlStr := `SELECT COUNT(1) FROM users_checkin_receive WHERE user_id = ? AND is_more = 1 AND date(created) = date(now())`
	err = Mysql.Chess.QueryRow(sqlStr, userId).Scan(&isMore)
	return
}

func (m *UsersCheckinReceiveModel) Insert() error {
	sqlStr := `INSERT INTO users_checkin_receive(user_id,chekin_days_reward_id,is_more) VALUES(?,?,?)`
	_, err := Mysql.Chess.Exec(sqlStr, m.UserId, m.CheckinDaysRewardId, m.IsMore)
	return err
}
