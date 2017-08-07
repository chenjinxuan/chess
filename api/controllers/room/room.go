package c_room

import (
    "chess/models"
    "github.com/gin-gonic/gin"
    "chess/common/define"
    "net/http"
)

type RoomsResult struct {
    define.BaseResult
    Bb int `json:"bb" description:"大盲注"`
    List []*RoomsInfo `json:"list"`
}
type RoomsInfo struct {
    Id int `json:"id"`
    BigBlind int `json:"big_blind" description:"大盲注"`
    SmallBlind int `json:"small_blind" description:"小盲注"`
    MinCarry int `json:"min_carry" description:"最小携带筹码"`
    MaxCarry int `json:"max_carry" description:"最大携带筹码"`
    Max int `json:"max" description:"最大人数"`
    Online int `json:"online" description:"在线人数"`
}
// @Title 获取房间列表信息
// @Description 获取房间列表信息
// @Summary 获取房间列表信息
// @Accept json
// @Success 200 {object} c_room.RoomsResult
// @router /room/list [get]
func RoomsList(c *gin.Context)  {
       var result RoomsResult
       data,err:=models.Rooms.GetAll()
	if err != nil {
	    result.Msg="get fail."
	    c.JSON(http.StatusOK,result)
	    return

	}
    for _,v:=range data {
	var info = new(RoomsInfo)
	info.Id=v.Id
	info.BigBlind=v.BigBlind
	info.SmallBlind=v.SmallBlind
	info.MinCarry=v.MinCarry
	info.MaxCarry=v.MaxCarry
	info.Max=v.Max

	//去游戏付获取在线人数
	//v.Online
	result.List=append(result.List,info)
    }
    result.Ret = 1
    c.JSON(http.StatusOK,result)
    return
}