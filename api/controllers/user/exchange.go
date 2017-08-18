package c_user

import (
    "chess/common/define"
    "github.com/gin-gonic/gin"
    "strconv"
    "chess/models"
    "strings"
    "net/http"
)

type ExchangeParams struct {
    Diamonds int `json:"diamonds" form:"diamonds"`
    From string `json:"from" form:"from"`
}
type ExchangeResult struct {
    define.BaseResult
} 
func Exchange(c *gin.Context)  {
    var params ExchangeParams
    var result ExchangeResult
    var user_id int
    var err error
    if user_id,err=strconv.Atoi(c.Param("user_id")); err!=nil {
	result.Msg = "Get user fail."
	c.JSON(http.StatusOK, result)
	return
    }

    if err=c.Bind(&params); err != nil {
	result.Msg = "bind fail."
	c.JSON(http.StatusOK, result)
	return
    }
    //获取用户钱包信息
    _,diamond,err:=models.UsersWallet.GetBalance(user_id)
    if err != nil {
	result.Msg = "Get wallet fail."
	c.JSON(http.StatusOK, result)
	return
    }
    if params.Diamonds < 0 || diamond < params.Diamonds {
	result.Msg = "price fail."
	c.JSON(http.StatusOK, result)
	return
    }

    var model = new(models.UsersWithDrawRecordModel)
    model.AppFrom=strings.ToLower(params.From)
    model.UserId = user_id
    model.Diamond = params.Diamonds
    model.DiamondBlance = diamond - params.Diamonds
    model.Count = params.Diamonds * 1000
    model.Status = 1
    
   err= models.UsersWithDrawRecord.Exchange(model)
    if err != nil {
	result.Msg = "exchange fail."
	c.JSON(http.StatusOK, result)
	return
    }
    result.Ret = 1
    c.JSON(http.StatusOK, result)
    return
}