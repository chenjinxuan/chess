package models

type TaskTypeModel struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	ExpireType int    `json:"expire_type"`
	Status     int    `json:"status"`
}

var TaskType = new(TaskTypeModel)

func (m *TaskTypeModel) List() (list []TaskTypeModel, err error) {
	sqlStr := `SELECT id ,name ,expire_type FROM task_type WHERE status = 1`
	rows, err := Mysql.Chess.Query(sqlStr)
	if err != nil {

	}
	defer rows.Close()
	for rows.Next() {
		var t TaskTypeModel
		err = rows.Scan(&t.Id, &t.Name, &t.ExpireType)
		if err != nil {
			continue
		}
		list = append(list, t)
	}
	return
}
