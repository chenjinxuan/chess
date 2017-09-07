package c_user

import (
	"chess/common/define"
	"chess/common/log"
	"chess/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	//grpcServer "chess/api/grpc"
	pb "chess/api/proto"
	//"golang.org/x/net/context"
	"encoding/json"
	"chess/api/redis"
)

type ExchangeParams struct {
	GoodsId int    `form:"goods_id" description:"购买兑换的商品id"`
	From    string `form:"from" description:"来源应用"`
}
type ExchangeResult struct {
	define.BaseResult
}

// @Title 钻石兑换金币
// @Description 钻石兑换金币 ret= -1 时 钻石不足
// @Summary 钻石兑换金币
// @Accept json
// @Param   token     query    string   true        "token"
// @Param   goods_id     query    string   true        "goods_id"
// @Param   from     query    string   true        "from"
// @Success 200 {object} c_user.ExchangeResult
// @router /user/{user_id}/exchange [get]
func Exchange(c *gin.Context) {
	var params ExchangeParams
	var result ExchangeResult
	var user_id int
	var err error
	if user_id, err = strconv.Atoi(c.Param("user_id")); err != nil {
		result.Msg = "Get user fail."
		c.JSON(http.StatusOK, result)
		return
	}

	if err = c.Bind(&params); err != nil {
		result.Msg = "bind fail."
		c.JSON(http.StatusOK, result)
		return
	}
	//获取商品信息
	goods, err := models.Goods.Get(params.GoodsId)
	if err != nil {
		log.Errorf("models.Goods.Get err(%s)", err)
		result.Msg = "get goods fail."
		c.JSON(http.StatusOK, result)
		return
	}
	//TODO 后续更加商品类型,属性做出相应的处理

	//获取用户钱包信息
	_, diamond, err := models.UsersWallet.GetBalance(user_id)
	if err != nil {
		log.Errorf("models.UsersWallet.GetBalance err(%s)", err)
		result.Msg = "Get wallet fail."
		c.JSON(http.StatusOK, result)
		return
	}
	if diamond < goods.Price {
		result.Ret = -1
		result.Msg = "price fail."
		c.JSON(http.StatusOK, result)
		return
	}

	var model = new(models.UsersWithDrawRecordModel)
	model.AppFrom = strings.ToLower(params.From)
	model.UserId = user_id
	model.Diamond = goods.Price
	model.DiamondBlance = diamond - goods.Price
	model.Count = goods.Price * goods.Rate
	model.Status = 1

	err = models.UsersWithDrawRecord.Exchange(model)
	if err != nil {
		log.Errorf("models.UsersWithDrawRecord.Exchange err(%s)", err)
		result.Msg = "exchange fail."
		c.JSON(http.StatusOK, result)
		return
	}
        //通知任务系统
    //通知是否要更新任务bag
	//TaskClient,ret := grpcServer.GetTaskGrpc()
	//if ret == 0{
	//
	//    result.Msg = "rpc fail"
	//    c.JSON(http.StatusOK, result)
	//    return
	//}
    //go TaskClient.IncrUserBag(context.Background(), &pb.UpdateBagArgs{UserId:int32(user_id),GoodsId:int32(params.GoodsId)})
    //因为加入背包事件比较重要 还是直接存redis吧
   strByte,err:= json.Marshal(pb.UpdateBagArgs{UserId:int32(user_id),GoodsId:int32(params.GoodsId)})
    go api_redis.Redis.Task.Lpush(define.TaskUserBagRedisKey,string(strByte))
	result.Ret = 1
	c.JSON(http.StatusOK, result)
	return
}
