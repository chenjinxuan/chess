package c_user

import (
    "github.com/gin-gonic/gin"
    "strconv"
    "chess/common/define"
    "chess/models"
    "net/http"
)

type UserInfo struct { 
    Id int `json:"id"`
    NickName string `json:"nick_name" description:"昵称"`
    MobileNumber string `json:"mobile_number" description:"手机号"`
    Gender  int `json:"gender" description:"性♂别 // 0 未知  1男2女 "`
    Avatar string `json:"avatar" description:"头像"`
    Type int `json:"type" description:"用户类型"`
    Status int `json:"status" description:"状态"`
    IsFresh int `json:"is_fresh" description:"是否新用户"`
    Balance int `json:"balance" description:"金币余额"`
    DiamondBalance int `json:"diamond_balance" description:"钻石余额"`
}

type UserInfoResult struct {
    define.BaseResult
    Data UserInfo `json:"data"`

}

// @Title 获取用户基本信息
// @Description 获取用户基本信息
// @Summary 获取用户基本信息
// @Accept json
// @Param   token     query    string   true        "token"
// @Param   user_id     path    int   true        "user_id"
// @Success 200 {object} c_user.UserInfoResult
// @router /user/{user_id}/info [get]
func GetUserInfo(c *gin.Context)  {
    var result UserInfoResult
    UserId ,err:=strconv.Atoi(c.Param("user_id"))
    if err != nil {
	result.Msg="bind params fail ."
	c.JSON(http.StatusOK,result)
	return
    }

    var user = new(models.UsersModel)
    //获取信息
    err = models.Users.Get(UserId,user)
    if err != nil {
	result.Msg="get info fail 1"
	c.JSON(http.StatusOK,result)
	return
    }
    result.Data.Id=user.Id
    result.Data.NickName=user.Nickname
    result.Data.MobileNumber=user.MobileNumber
    result.Data.Gender=user.Gender
    result.Data.Avatar=user.Avatar
    result.Data.Type=user.Type
    result.Data.Status=user.Status
    result.Data.IsFresh=user.IsFresh
    //余额查询
    result.Data.Balance,result.Data.DiamondBalance,err = models.UsersWallet.GetBalance(UserId)
    if err != nil {
	result.Msg="get info fail 2"
	c.JSON(http.StatusOK,result)
	return
    }
    result.Ret = 1
    c.JSON(http.StatusOK,result)
    return
}