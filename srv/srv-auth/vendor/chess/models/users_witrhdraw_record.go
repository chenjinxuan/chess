package models

type UsersWithDrawRecordModel struct {
	Id            int    `json:"id"`
	UserId        int    `json:"user_id"`
	Diamond       int    `json:"diamond"`
	DiamondBlance int    `json:"diamond_blance"`
	Count         int    `json:"count"`
	Type          int    `json:"type"`
	Status        int    `json:"status"`
	AppFrom       string `json:"app_from"`
	Note          string `json:"note"`
}

var UsersWithDrawRecord = new(UsersWithDrawRecordModel)

func (m *UsersWithDrawRecordModel) Exchange(data *UsersWithDrawRecordModel) error {
	tx, _ := Mysql.Chess.Begin()
	sqlStr := "UPDATE users_wallet SET balance = balance + ? ,total = total + ? ,diamond_balance = diamond_balance - ? WHERE user_id = ?"
	_, err := tx.Exec(sqlStr, data.Count, data.Count, data.Diamond, data.UserId)
	if err != nil {
		tx.Rollback()
		return err
	}
	sqlStr = `INSERT INTO users_withdraw_record(user_id,diamond,diamond_balance,count,type,status,app_from) VALUES(?,?,?,?,?,?,?)`
	_, err = tx.Exec(sqlStr, data.UserId, data.Diamond, data.DiamondBlance, data.Count, data.Type, data.Status, data.AppFrom)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
