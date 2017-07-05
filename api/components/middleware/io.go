package middleware

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"chess/common/config"
	"chess/common/helper"
	"chess/common/log"
)

func BindJSON(c *gin.Context, params interface{}) error {
	defer c.Request.Body.Close()
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	log.Debugf("Post Data Des String(%s) Key(%s)", string(body), config.Api.PostDesKey)

	err = json.Unmarshal(body, params)
	if err == nil {
		return binding.Validator.ValidateStruct(params)
	}

	// 解密
	text := helper.DesDecryptECB(config.Api.PostDesKey, string(body))
	if text == "" {
		return errors.New("Decrypt post data fail")
	}

	err = json.Unmarshal([]byte(text), params)
	if err != nil {
		return err
	}

	log.Debugf("Decrypt post data success! %+v", params)

	return binding.Validator.ValidateStruct(params)
}
