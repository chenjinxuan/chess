package models

import "time"

var Rooms = new(RoomsModel)

type RoomsModel struct {
	Id int
	BigBlind int
	SmallBlind int
	MinCarry int
	MaxCarry int
	Max int
	Updated time.Time
	Created time.Time
}

func(m *RoomsModel) GetAll() (list []RoomsModel, err error) {
	sqlStr := `SELECT id, big_blind, small_blind, min_carry, max_carry, max
		FROM rooms`

	rows, err := Mysql.Chess.Query(sqlStr)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var data RoomsModel
		err = rows.Scan(
			&data.Id,
			&data.BigBlind,
			&data.SmallBlind,
			&data.MinCarry,
			&data.MaxCarry,
			&data.Max,
		)
		if err != nil {
			continue
		}
		list = append(list, data)
	}
	return
}