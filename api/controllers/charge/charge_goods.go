package c_charge

import (
    "github.com/gin-gonic/gin"
    "chess/common/define"
    "chess/models"
    "net/http"
    "chess/common/log"
    "chess/common/config"
)

type ChargeGoodsResult struct {
    define.BaseResult
    List []ChargeGoodsInfo `json:"list"`
}
type ChargeGoodsInfo struct {
    ChargeGoodsId string `json:"charge_goods_id" description:"商品id"`
    Name string `json:"name" description:"商品名"`
    Price int `json:"price" description:"价格"`
    Number int `json:"number" description:"钻石数量"`
    Image string `json:"image" description:"图片"`
}
// @Title 充值商品列表
// @Description 充值商品列表
// @Summary 充值商品列表
// @Accept json
// @Success 200 {object} c_charge.ChargeGoodsResult
// @router /charge/charge_goods/list [get]
func ChargeGoodsList(c *gin.Context) {
    var result ChargeGoodsResult
    _conf, ok1 := c.Get("config")
    cConf, ok2 := _conf.(*config.ApiConfig)
    if !ok1 || !ok2 {
	result.Msg = "Get config fail."
	c.JSON(http.StatusOK, result)
	return
    }
    list,err := models.ChargeGoods.List()
    if err != nil {
	log.Errorf(" models.ChargeGoods.List err:(%s)" ,err)
	result.Msg= "get list fail ."
	c.JSON(http.StatusOK, result)
	return
    }
    for _,v := range list {
	var data ChargeGoodsInfo
	data.Price=v.Price
	data.Number=v.Number
	data.Name=v.Name
	data.ChargeGoodsId= v.ChargeGoodsId
	data.Image = cConf.Backend.ImageDomainPrefix + v.Image
	result.List=append(result.List,data)
    }
    result.Ret = 1
    c.JSON(http.StatusOK, result)
    return
}