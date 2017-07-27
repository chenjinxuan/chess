package controllers


import (
    "chess/models"
    "github.com/gin-gonic/gin"
    "chess/api/define"
    "net/http"
    "fmt"
)

type Result struct {
    define.BaseResult
    Data *models.UsersModel `json:"data"`
}

func Get(c *gin.Context) {
    var result Result
    var user =new(models.UsersModel)
    err:=models.Users.GetByMobileNumber("15345929485", user)
    if err != nil {
	fmt.Println(err)
    }
    result.Data = user
    result.Ret =1
    c.JSON(http.StatusOK, result)
    return
}