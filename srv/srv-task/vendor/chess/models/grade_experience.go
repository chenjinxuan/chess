package models

type GradeExperienceModel struct {
	Id            int    `json:"id"`
	Grade         int    `json:"grade"`
	GradeDescribe string `json:"grade_describe"`
	Experience    int    `json:"experience"`
}

var GradeExperience = new(GradeExperienceModel)

func (m *GradeExperienceModel) GetAll() (list []GradeExperienceModel, err error) {
	sqlStr := `SELECT id,grade,grade_describe,experience FROM grade_experience`
	rows, err := Mysql.Chess.Query(sqlStr)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var g GradeExperienceModel
		err = rows.Scan(&g.Id, &g.Grade, &g.GradeDescribe, &g.Experience)
		if err != nil {
			continue
		}
		list = append(list, g)
	}
	return
}
