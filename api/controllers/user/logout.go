package c_user

import (
    "github.com/gin-gonic/gin"
    "chess/common/define"
    grpcServer "chess/api/grpc"
    pb "chess/api/proto"
    "golang.org/x/net/context"
    "net/http"
)
type TokenOutParams struct {
    Token    string `form:"token" binding:"required" description:"token"`
}
type TokenOutResult struct {
    define.BaseResult
}
// @Title 登出 (token过期)
// @Description 登出 (token过期)
// @Summary 登出 (token过期)
// @Accept json
// @Param   token     query    string   true        "token"
// @Param   user_id     path    int   true        "user_id"
// @Success 200 {object} c_user.TokenOutResult
// @router /user/{user_id}/logout [get]
func Logout (c *gin.Context) {
    var params TokenOutParams
    var result TokenOutResult
    // Generate a new login token
    if err:=c.Bind(&params) ;err!=nil {
	result.Msg = "params bind fail ."
	c.JSON(http.StatusOK, result)
	return
    }
    AuthClient := grpcServer.GetAuthGrpc()
    authResult, err := AuthClient.BlackToken(context.Background(), &pb.BlackTokenArgs{Code:define.AuthExpire,Token:params.Token})
    if err != nil {
	result.Msg = "logout fail ."
	c.JSON(http.StatusOK, result)
	return
    }
    result.Ret = int(authResult.Ret)
    c.JSON(http.StatusOK, result)
    return
}
