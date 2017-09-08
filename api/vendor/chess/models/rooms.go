package models

import "time"

var Rooms = new(RoomsModel)

type RoomsModel struct {
	Id          int       `json:"id"`
	BigBlind    int       `json:"big_blind" description:"小盲注"`
	SmallBlind  int       `json:"small_blind" description:"小盲注"`
	MinCarry    int       `json:"min_carry" description:"最小携带筹码"`
	MaxCarry    int       `json:"max_carry" description:"最大携带筹码"`
	Max         int       `json:"max" description:"最大人数"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	RoomsTypeId int       `json:"rooms_type_id"`
}

func (m *RoomsModel) GetAll() (list []RoomsModel, err error) {
	sqlStr := `SELECT id, big_blind, small_blind, min_carry, max_carry, max, rooms_type_id
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
			&data.RoomsTypeId,
		)
		if err != nil {
			continue
		}
		list = append(list, data)
	}
	return
}
