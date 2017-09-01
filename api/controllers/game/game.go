package c_game

import (
    "chess/common/define"
    "chess/models"
    "github.com/gin-gonic/gin"
    "net/http"
    "fmt"
    "chess/common/log"
)

type LastGameResult struct {
    define.BaseResult
    Data models.GamblingModel `json:"data"`
}

type LastGameParams struct {
    TableId string `form:"table_id"`
}
// @Title 上局牌局
// @Description 上局牌局
// @Summary 上局牌局
// @Accept json
// @Param   user_id     path    int  true        "用户id"
// @Param   token     query    string  true        "token"
// @Param   table_id     query    string  true        "table_id"
// @Success 200 {object} c_game.LastGameResult
// @router /game/{user_id}/last_game [get]
func LastGame(c *gin.Context)  {
    var result LastGameResult
    var params LastGameParams
    if err:=c.Bind(&params);err!=nil {
	result.Msg = "bind fail ."
	c.JSON(http.StatusOK, result)
	return
    }
    //查出最近一次牌局
    fmt.Println(params.TableId)
    data,err:=models.Gambling.GetByTableId(params.TableId)
    if err != nil {
	if fmt.Sprint(err)=="not found" {
	    result.Ret = 1
	    c.JSON(http.StatusOK, result)
	    return
	}
	log.Errorf("models.Gambing.GetByTableId",err)
	result.Msg = "get info fail 2"
	c.JSON(http.StatusOK, result)
	return
    }
    result.Data=data
    result.Ret = 1
    c.JSON(http.StatusOK, result)
    return
}