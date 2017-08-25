package models

type TaskRewardTypeModel struct {
    Id int `json:"id"`
    Name string `json:"name"`
    Status int `json:"status"`
}

var TaskRewardType = new(TaskRewardTypeModel)
func (m *TaskRewardTypeModel) List() (list []TaskRewardTypeModel,err error) {
    sqlStr := `SELECT id ,name  FROM task_reward_type WHERE status = 1`
    rows,err:=Mysql.Chess.Query(sqlStr)
    if err != nil {

    }
    defer rows.Close()
    for rows.Next()  {
	var t TaskRewardTypeModel
	err = rows.Scan(&t.Id,&t.Name)
	if err != nil {
	    continue
	}
	list = append(list,t)
    }
    return
}