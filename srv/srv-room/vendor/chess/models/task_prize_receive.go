package models

type TaskPriceReceiveModel struct {
	Id         int
	UserId     int
	TaskId     int
	RewardType int
	RewardNum  int
}

var TaskPriceReceive = new(TaskPriceReceiveModel)

func (m *TaskPriceReceiveModel) Insert() error {
	sqlStr := `INSERT INTO task_prize_receive(user_id,task_id,reward_type,reward_num) VALUES(?,?,?,?)`
	_, err := Mysql.Chess.Exec(sqlStr, m.UserId, m.TaskId, m.RewardType, m.RewardNum)
	return err
}

func (m *TaskPriceReceiveModel) GetAllByUserId(userId int) (list []int,err error){
        sqlStr := `SELECT task_id FROM task_prize_receive WHERE user_id  =  ?`
       rows,err:=Mysql.Chess.Query(sqlStr,userId)
	if err != nil {
	    return
	}
	defer rows.Close()
    for rows.Next() {
	var i int
	err=rows.Scan(&i)
	list = append(list,i)
    }
    return
}