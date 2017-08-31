package c_user

import (
    "chess/common/define"
    "github.com/gin-gonic/gin"
    "strconv"
    "net/http"
    "chess/common/config"
    "chess/models"
    "fmt"
    "chess/common/log"
    "time"
)

type BagListResult struct {
    define.BaseResult
    UserId int `json:"user_id" description:"用户id"`
    List [] BagDetail  `json:"list"`
}

type BagDetail struct {
    GoodsId             int `json:"goods_id" description:"商品id"`
    Name            string `json:"name" description:"商品名"`
    GoodsCategoryId int `json:"goods_category_id" description:"商品类别"`
    GoodsTypeId     int `json:"goods_type_id" description:"商品属性"`
    IsExpire        int  `json:"is_expire" description:"0为永久,其他为过期时间 时间戳"`
    GoodsDescribe   string `json:"goods_describe" description:"描述"`
    Image           string `json:"image" description:"图片"`
    Number          int `json:"number" description:"数量"`
}
// @Title 用户背包
// @Description 用户背包
// @Summary 用户背包
// @Accept json
// @Param   user_id     path    int  true        "用户id"
// @Param   token     query    string  true        "token"
// @Success 200 {object} c_user.BagListResult
// @router /user/{user_id}/bag/list [get]
func BagList(c *gin.Context)  {
    var result BagListResult
    userId, err := strconv.Atoi(c.Param("user_id"))
    if err != nil {
	result.Msg = "bind fail ."
	c.JSON(http.StatusOK, result)
	return
    }
    _conf, ok1 := c.Get("config")
    cConf, ok2 := _conf.(*config.ApiConfig)
    if !ok1 || !ok2 {
	result.Msg = "Get config fail."
	c.JSON(http.StatusOK, result)
	return
    }
    //从mongo查出改用的背包
    bag,err:=models.UserBag.Get(userId)
    if err != nil {
	if fmt.Sprint(err) == "not found" {
	    result.Ret =1
	    c.JSON(http.StatusOK, result)
	    return
	}
	log.Errorf("models.UserBag.Get", err)
	result.Msg = "get bag fail"
	c.JSON(http.StatusOK, result)
	return
    }
    //循环查出背包的具体商品信息
    var goodsIdList []string
    var numberMap = make(map[int]int)
    for _,v:=range bag.List {
	//nunber = 0  时删除改商品
	if v.Number <= 0 {
	    go models.UserBag.RemoveByGoodsId(userId,v.GoodsId)
	    continue
	}
	goodsIdList=append(goodsIdList,strconv.Itoa(v.GoodsId))
	numberMap[v.GoodsId]=v.Number
    }
    goodsList,err:=models.Goods.GetBySlice(goodsIdList)
    if err != nil {
	log.Errorf("models.Goods.GetBySlice", err)
	result.Msg = "already checkin today"
	c.JSON(http.StatusOK, result)
	return
    }
    t:=time.Now().Unix()
    for _,v:=range goodsList  {
	if v.IsExpire !=0 {
	    if int(t)>v.IsExpire {
		//并且删除此物品
		go models.UserBag.RemoveByGoodsId(userId,v.Id)
		continue
	    }
	}
	var bagDetail BagDetail
	bagDetail.Number=numberMap[v.Id]
	bagDetail.GoodsId=v.Id
	bagDetail.GoodsDescribe=v.GoodsDescribe
	bagDetail.GoodsTypeId=v.GoodsTypeId
	bagDetail.GoodsCategoryId=v.GoodsCategoryId
	bagDetail.Image=cConf.Backend.ImageDomainPrefix + v.Image
	bagDetail.IsExpire=v.IsExpire
	bagDetail.Name=v.Name
	result.List=append(result.List,bagDetail)
    }
    result.Ret =1
    c.JSON(http.StatusOK, result)
    return
}

type BagUseParmas struct {
    GoodsId int `form:"goods_id"`
}
func BagUse(c *gin.Context)  {
    var result define.BaseResult
    var params BagUseParmas
    userId, err := strconv.Atoi(c.Param("user_id"))
    if err != nil {
	result.Msg = "bind fail ."
	c.JSON(http.StatusOK, result)
	return
    }
    if err:=c.Bind(&userId);err!= nil {

    }
    //查出bag
    bag,err:= models.UserBag.Get(userId)
    for _,v:=range bag.List {
	if v.GoodsId == params.GoodsId {//找到该商品
	    //判断数量
	    if v.Number <= 0 {
		go models.UserBag.RemoveByGoodsId(userId,v.GoodsId)
	    }
	}
    }
}