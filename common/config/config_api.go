package config

import (
	. "chess/common/consul"
	"chess/common/log"
	"encoding/json"
)

var Api = new(ApiConfig)

type ApiConfig struct {
	PublicConfig
	DbConfig

	Debug      bool          // debug 模式
	Port       string        // http 监听端口
	PostDesKey string        // post 数据加密 des key
	Rpc        *ApiRpcConfig // rpc server地址
}

type ApiRpcConfig struct {
	Notifier string `json:"notifier"`
}

func (c *ApiConfig) Import() error {
	var err error

	err = c.PublicConfig.Import()
	if err != nil {
		return err
	}

	err = c.DbConfig.Import("api")
	if err != nil {
		return err
	}

	c.Debug, err = ConsulClient.KeyBool("api/debug", true)
	if err != nil {
		return err
	}
	ConsulClient.KeyBoolWatch("api/debug", &c.Debug)

	c.Port, err = ConsulClient.Key("api/port", ":8888")
	if err != nil {
		return err
	}
	ConsulClient.KeyWatch("api/port", &c.Port)

	c.PostDesKey, err = ConsulClient.Key("api/post_des_key", "XQ1R1%8f")
	if err != nil {
		return err
	}
	ConsulClient.KeyWatch("api/post_des_key", &c.PostDesKey)

	rpcStr, err := ConsulClient.Key("api/rpc", "")
	if err != nil {
		return err
	}
	c.Rpc = new(ApiRpcConfig)
	err = json.Unmarshal([]byte(rpcStr), c.Rpc)
	if err != nil {
		return err
	}

	log.Debugf("Api config import success! [%+v]", *c)
	return nil
}
