package models

type GoodsModel struct {
	Id              int
	Name            string
	GoodsCategoryId int
	GoodsTypeId     int
	Price           int
	Rate            int
	IsExpire        int
	GoodsDescribe   string
	Image           string
	Status          int
}

var Goods = new(GoodsModel)

func (m *GoodsModel) List() (list []GoodsModel, err error) {
	sqlStr := `SELECT id,name,goods_category_id,goods_type_id,price,goods_describe,image FROM goods WHERE status = 1`
	rows, err := Mysql.Chess.Query(sqlStr)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var g GoodsModel
		err = rows.Scan(&g.Id, &g.Name, &g.GoodsCategoryId, &g.GoodsTypeId, &g.Price, &g.GoodsDescribe, &g.Image)
		if err != nil {
			continue
		}
		list = append(list, g)
	}
	return
}

func (m *GoodsModel) Get(goodsId int) (g GoodsModel, err error) {
	sqlStr := `SELECT id,name,goods_category_id,goods_type_id,price,rate,is_expire,goods_describe,image FROM goods WHERE id = ? AND status = 1`
	err = Mysql.Chess.QueryRow(sqlStr, goodsId).Scan(&g.Id, &g.Name, &g.GoodsCategoryId, &g.GoodsTypeId, &g.Price, &g.Rate, &g.IsExpire, &g.GoodsDescribe, &g.Image)
	return
}
