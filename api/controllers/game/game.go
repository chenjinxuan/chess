package c_game

import (
	"chess/common/define"
	"chess/common/log"
	"chess/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
    "strconv"
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
func LastGame(c *gin.Context) {
	var result LastGameResult
	var params LastGameParams
	if err := c.Bind(&params); err != nil {
		result.Msg = "bind fail ."
		c.JSON(http.StatusOK, result)
		return
	}
	//查出最近一次牌局
	fmt.Println(params.TableId)
	data, err := models.Gambling.GetByTableId(params.TableId)
	if err != nil {
		if fmt.Sprint(err) == "not found" {
			result.Ret = 1
			c.JSON(http.StatusOK, result)
			return
		}
		log.Errorf("models.Gambing.GetByTableId", err)
		result.Msg = "get info fail 2"
		c.JSON(http.StatusOK, result)
		return
	}
	result.Data = data
	result.Ret = 1
	c.JSON(http.StatusOK, result)
	return
}

type GameListParams struct {
    PageSize int `form:"page_size"`
    PageNum int `form:"page_num"`
}

type GameListResult struct {
    define.BaseResult
    Count int `json:"count"`
    List []models.GamblingModel `json:"list"`
}
// @Title  牌局记录
// @Description 牌局记录 没传分页要求时默认取10条
// @Summary 牌局记录
// @Accept json
// @Param   user_id     path    int  true        "用户id"
// @Param   token     query    string  true        "token"
// @Param   page_size     query    string  true        "分页大小"
// @Param   page_num     query    string  true        "第几页"
// @Success 200 {object} c_game.GameListResult
// @router /game/{user_id}/game_list [get]
func GameList(c *gin.Context)  {
    var result GameListResult
    var params GameListParams
    userId, err := strconv.Atoi(c.Param("user_id"))
    if err != nil {
	result.Msg = "get userId fail ."
	c.JSON(http.StatusOK, result)
	return
    }
    //init
    params.PageNum=1
    params.PageSize=10
    if err := c.Bind(&params); err != nil {
	result.Msg = "bind fail ."
	c.JSON(http.StatusOK, result)
	return
    }
    result.Count,err=models.Gambling.GetCountByUserId(userId)
    if err != nil {
	log.Errorf("models.Gambling.GetCountByUserId", err)
	result.Msg = "get count info fail ."
	c.JSON(http.StatusOK, result)
	return
    }
    if result.Count == 0 {
	result.Ret = 1
	c.JSON(http.StatusOK, result)
	return
    }
    //查出数据
    result.List,err=models.Gambling.GetByUserId(userId,params.PageSize,params.PageNum)
    if err != nil {
	if fmt.Sprint(err) == "not found" {
	    result.Ret = 1
	    c.JSON(http.StatusOK, result)
	    return
	}
	log.Errorf("models.Gambling.GetByUserId", err)
	result.Msg = "get info fail ."
	c.JSON(http.StatusOK, result)
	return
    }
    result.Ret=1
    c.JSON(http.StatusOK, result)
    return
}