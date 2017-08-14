package models

var Device = new(DeviceModel)

type DeviceModel struct {
	Id       int    `json:"id"`
	UserId   int    `json:"aid"`
	UniqueId string `json:"unique_id" binding:"required"`
	Type     int    `json:"type" binding:"required"`
	Openudid string `json:"openudid"`
	Idfa     string `json:"idfa"`
	Idfv     string `json:"idfv"`
	Imei     string `json:"imei"`
	Imsi     string `json:"imsi"`
	Mac      string `json:"mac"`
	Language string `json:"language"`
	Manu     string `json:"manu"`
	Model    string `json:"model"`
	RomInfo  string `json:"rom_info"`
	OsVer    string `json:"os_ver"`
}

func (m *DeviceModel) Upsert() error {
	sqlStr := `INSERT INTO device
		(user_id, unique_id, type, openudid, idfa, idfv, imei, imsi, mac, language, manu, model, rom_info, os_ver)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		openudid = ?, idfa = ?, idfv = ?, imei = ?, imsi = ?, mac = ?, language = ?, manu = ?, model = ?, rom_info = ?, os_ver = ?`

	_, err := Mysql.Chess.Exec(sqlStr,
		m.UserId, m.UniqueId, m.Type, m.Openudid, m.Idfa, m.Idfv, m.Imei, m.Imsi, m.Mac, m.Language, m.Manu, m.Model, m.RomInfo, m.OsVer,
		m.Openudid, m.Idfa, m.Idfv, m.Imei, m.Imsi, m.Mac, m.Language, m.Manu, m.Model, m.RomInfo, m.OsVer,
	)
	return err
}

func (m *DeviceModel) CountByIdfa(idfa string) (cnt int, err error) {
	sqlStr := `SELECT COUNT(1)
		FROM device
		WHERE idfa = ?`

	err = Mysql.Chess.QueryRow(sqlStr, idfa).Scan(&cnt)
	return
}

func (m *DeviceModel) GetUserIdByUniqueId(uniqueId string) (res int, err error) {
	sqlStr := `SELECT user_id FROM device WHERE unique_id = ? AND user_id != 0 and user_id != -1`

	err = Mysql.Chess.QueryRow(sqlStr, uniqueId).Scan(&res)
	return
}

func (m *DeviceModel) GetUserIdByIdfv(idfv string) (res int, err error) {
	sqlStr := `SELECT user_id FROM device WHERE idfv = ? AND user_id != 0 and user_id != -1`

	err = Mysql.Chess.QueryRow(sqlStr, idfv).Scan(&res)
	return
}

func (m *DeviceModel) GetUserIdByIdfa(idfa string) (res int, err error) {
	sqlStr := `SELECT user_id FROM device WHERE idfa = ? AND user_id != 0 and user_id != -1`

	err = Mysql.Chess.QueryRow(sqlStr, idfa).Scan(&res)
	return
}
