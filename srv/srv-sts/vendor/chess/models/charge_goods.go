package models

import (
	"time"
)

type ChargeGoodsModel struct {
	Id            int    `json:"id"`
	ChargeGoodsId string `json:"charge_goods_id"`
	Name          string `json:"name"`
	Price         int    `json:"price"`
	Number        int    `json:"number"`
	Note          string `json:"note"`
	Status        int    `json:"status"`
	Image         string
	Created       time.Time
}

const (
	ChargeGoodsStatusUnavailable = 0
	ChargeGoodsStatusAvailable   = 1
)

var ChargeGoods = new(ChargeGoodsModel)

func (m *ChargeGoodsModel) List() (list []ChargeGoodsModel, err error) {
	sqlStr := `SELECT charge_goods_id , name ,price ,number,image FROM charge_goods WHERE status = ?`
	rows, err := Mysql.Chess.Query(sqlStr, ChargeGoodsStatusAvailable)
	defer rows.Close()
	for rows.Next() {
		var l ChargeGoodsModel
		err = rows.Scan(&l.ChargeGoodsId, &l.Name, &l.Price, &l.Number, &l.Image)
		if err != nil {
			continue
		}
		list = append(list, l)
	}
	return
}
