package models

import "time"

var UsersAward = new(UsersAwardModel)

type UsersAwardModel struct {
	Id int
	UserId int
	TableId string
	RoomId int
	Num int
	Created time.Time
}

func (m *UsersAwardModel) Insert() error {
	sqlStr := `INSERT INTO users_award
		(user_id, table_id, room_id, num)
		VALUES
		(?, ?, ?, ?)`

	_, err := Mysql.Chess.Exec(sqlStr, m.UserId, m.TableId, m.RoomId, m.Num)
	return err
}