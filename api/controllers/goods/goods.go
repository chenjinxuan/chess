package c_goods

import (
    "github.com/gin-gonic/gin"
    "chess/common/define"
    "net/http"
    "chess/common/config"
    "chess/models"
    "chess/common/log"
)

type ListResult struct {
    define.BaseResult
    List []GoodsInfo `json:"list"`
}
type GoodsInfo struct {
    Id int `json:"id"`
    Name string `json:"name"`
    GoodsCategoryId int `json:"goods_category_id" description:"商品类别"`
    GoodsTypeId int `json:"goods_type_id" description:"商品属性"`
    Price int `json:"price" description:"商品价格->钻石"`
    GoodsDescribe string `json:"goods_describe" description:"商品描述"`
    Image string `json:"image" description:"商品图片"`
}
// @Title 商品列表
// @Description 商品列表
// @Summary 商品列表
// @Accept json
// @Success 200 {object} c_goods.ListResult
// @router /goods/list [get]
func List(c *gin.Context)  {
    var result ListResult
    _conf, ok1 := c.Get("config")
    cConf, ok2 := _conf.(*config.ApiConfig)
    if !ok1 || !ok2 {
	result.Msg = "Get config fail."
	c.JSON(http.StatusOK, result)
	return
    }
    list,err := models.Goods.List()
    if err != nil {
     log.Errorf("models.Goods.List err:(%s)",err)
     result.Msg = "get fail."
	c.JSON(http.StatusOK,result)
	return
    }

    for _,v :=range list {
	var goodsInfo GoodsInfo
	goodsInfo.Id=v.Id
	goodsInfo.Name=v.Name
	goodsInfo.GoodsCategoryId=v.GoodsCategoryId
	goodsInfo.GoodsTypeId=v.GoodsTypeId
	goodsInfo.Price=v.Price
	goodsInfo.GoodsDescribe=v.GoodsDescribe
	goodsInfo.Image = cConf.Backend.ImageDomainPrefix + v.Image
	result.List=append(result.List,goodsInfo)
    }
    result.Ret=1
    c.JSON(http.StatusOK,result)
    return
}
