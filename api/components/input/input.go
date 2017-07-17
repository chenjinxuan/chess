package input

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"chess/api/config"
	"chess/api/helper"
	"chess/api/log"
)

func BindJSON(c *gin.Context, params interface{}, cConf *config.Config) error {
	defer c.Request.Body.Close()
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	log.Log.Debugf("Parms Des String(%s) Key(%s)", string(body), cConf.Backend.ParamsDesKey)

	err = json.Unmarshal(body, params)
	if err == nil {
		log.Log.Debug("input.BindJSON json decode error. ", err)
		return binding.Validator.ValidateStruct(params)
	}

	// 解密
	text := helper.DesDecryptECB(cConf.Backend.ParamsDesKey, string(body))
	if text == "" {
		return errors.New("decrypt params fail")
	}

	err = json.Unmarshal([]byte(text), params)
	if err != nil {
		return err
	}

	return binding.Validator.ValidateStruct(params)
}
