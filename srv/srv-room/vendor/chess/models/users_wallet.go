package models

const (
	UsersWalletStatusAvailable   = 1
	UsersWalletStatusUnavailable = 0
)

var UsersWallet = new(UsersWalletModel)

type UsersWalletModel struct {
	Id             int
	UserId         int
	Balance        uint
	Total          int
	DiamondBalance uint
	DiamondTotal   int
	VirIsNew       int
	Status         int
}

func (m *UsersWalletModel) Init(userId int) error {
	sqlString := `INSERT INTO users_wallet (user_id, status) VALUES  (?, ?)`

	_, err := Mysql.Chess.Exec(sqlString, userId, 1)

	return err
}

func (m *UsersWalletModel) Get(userId int, data *UsersWalletModel) error {
	sqlString := `SELECT
					user_id, balance, total, status
				FROM users_wallet
				WHERE user_id = ?`

	return Mysql.Chess.QueryRow(
		sqlString, userId,
	).Scan(
		&data.UserId,
		&data.Balance,
		&data.Total,
		&data.Status,
	)
}
func (m *UsersWalletModel) GetBalance(userId int) (balance, diamondBalance int, err error) {
	sqlString := `SELECT
					 balance,diamond_balance
				FROM users_wallet
				WHERE user_id = ?`

	err = Mysql.Chess.QueryRow(sqlString, userId).Scan(&balance, &diamondBalance)
	return
}

func (m *UsersWalletModel) GetBalanceByMobile(mobile string) (balance int, err error) {
	sqlString := `SELECT balance 
		FROM users_wallet,users
		WHERE users_wallet.user_id = users.id AND users.mobile_number = ?`
	err = Mysql.Chess.QueryRow(sqlString, mobile).Scan(&balance)
	return
}

func (m *UsersWalletModel) SendImitPresent(uid, amount int) error {
	tx, err := Mysql.Chess.Begin()
	if err != nil {
		return err
	}

	sqlStr := `UPDATE users_wallet
		SET diamond_balance = diamond_balance + ?, diamond_total = diamond_total + ?
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
	_, err = tx.Exec(sqlStr, uid, amount, 1, "add diamond_balance", "imitation present")
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (m *UsersWalletModel) AddVirBalance(uid, amount int) error {
	tx, err := Mysql.Chess.Begin()
	if err != nil {
		return err
	}

	sqlStr := `UPDATE users_wallet
		SET diamond_balance = diamond_balance + ?, diamond_total = diamond_total + ?
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
	_, err = tx.Exec(sqlStr, uid, amount, 1, "add diamond_balance", "charge")
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (m *UsersWalletModel) Checkout(uid, add int) error {
	sqlStr := `UPDATE users_wallet
		SET balance = IF(balance + ? < 0, 0, balance + ?)
		WHERE user_id = ?`

	_, err := Mysql.Chess.Exec(sqlStr, add, add, uid)
	return err
}
