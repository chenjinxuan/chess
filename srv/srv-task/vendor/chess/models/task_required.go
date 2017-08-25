package models

type TaskRequiredModel struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	RoomType     int    `json:"room_type"`
	MatchType    int    `json:"match_type"`
	Status       int    `json:"status"`
	HandLevel    int    `json:"hand_level"`
	PlayerAction int    `json:"player_action"`
}

var TaskRequired = new(TaskRequiredModel)

func (m *TaskRequiredModel) List() (list []TaskRequiredModel, err error) {
	sqlStr := `SELECT id ,name,room_type,match_type,hand_level,player_action  FROM task_required WHERE status = 1`
	rows, err := Mysql.Chess.Query(sqlStr)
	if err != nil {

	}
	defer rows.Close()
	for rows.Next() {
		var t TaskRequiredModel
		err = rows.Scan(&t.Id, &t.Name, &t.RoomType, &t.MatchType, &t.HandLevel, &t.PlayerAction)
		if err != nil {
			continue
		}
		list = append(list, t)
	}
	return
}
