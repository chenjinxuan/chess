package models

const (
	UsersWalletStatusAvailable   = 1
	UsersWalletStatusUnavailable = 0
)

var UsersWallet = new(UsersWalletModel)

type UsersWalletModel struct {
	Id         int
	UserId     int
	Balance    uint
	Total      int
	VirBalance uint
	VirTotal   int
	VirIsNew   int
	Status     int
}

func (m *UsersWalletModel) Init(userId int) error {
	sqlString := `INSERT INTO users_wallet (user_id, status) VALUES  (?, ?)`

	_, err := ChessMysql.Main.Exec(sqlString, userId, 1)

	return err
}

func (m *UsersWalletModel) Get(userId int, data *UsersWalletModel) error {
	sqlString := `SELECT
					user_id, balance, total, vir_balance, vir_total, vir_is_new, status
				FROM users_wallet
				WHERE user_id = ?`

	return ChessMysql.Main.QueryRow(
		sqlString, userId,
	).Scan(
		&data.UserId,
		&data.Balance,
		&data.Total,
		&data.VirBalance,
		&data.VirTotal,
		&data.VirIsNew,
		&data.Status,
	)
}

func (m *UsersWalletModel) GetBalanceByMobile(mobile string) (balance int, err error) {
	sqlString := `SELECT balance 
		FROM users_wallet,users
		WHERE users_wallet.user_id = users.id AND users.mobile_number = ?`
	err = ChessMysql.Main.QueryRow(sqlString, mobile).Scan(&balance)
	return
}

func (m *UsersWalletModel) SendImitPresent(uid, amount int) error {
	tx, err := ChessMysql.Main.Begin()
	if err != nil {
		return err
	}

	sqlStr := `UPDATE users_wallet
		SET vir_balance = vir_balance + ?, vir_total = vir_total + ?, vir_is_new = 0
		WHERE user_id = ?`

	_, err = tx.Exec(sqlStr, amount, amount, uid)
	if err != nil {
		tx.Rollback()
		return err
	}

	sqlStr = `INSERT INTO users_addcash_log
		(user_id,amount,status,tag,comment)
		VALUES
		(?,?,?,?,?)`
	_, err = tx.Exec(sqlStr, uid, amount, 1, "add vir_balance", "imitation present")
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (m *UsersWalletModel) AddVirBalance(uid, amount int) error {
	tx, err := ChessMysql.Main.Begin()
	if err != nil {
		return err
	}

	sqlStr := `UPDATE users_wallet
		SET vir_balance = vir_balance + ?, vir_total = vir_total + ?
		WHERE user_id = ?`

	_, err = tx.Exec(sqlStr, amount, amount, uid)
	if err != nil {
		tx.Rollback()
		return err
	}

	sqlStr = `INSERT INTO users_addcash_log
		(user_id,amount,status,tag,comment)
		VALUES
		(?,?,?,?,?)`
	_, err = tx.Exec(sqlStr, uid, amount, 1, "add vir_balance", "charge")
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
