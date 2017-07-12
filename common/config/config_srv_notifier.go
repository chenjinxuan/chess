package config

import (
	"encoding/json"
	"chess/common/log"
	. "chess/common/consul"
)

var SrvNotifier = new(SrvNotifierConfig)

type SrvNotifierConfig struct {
	PublicConfig
	DbConfig

	ServerAliasName string //服务别名
	RPCPort         string

	Sms *SrvNotifierSmsConfig
}

type SrvNotifierSmsConfig struct {
	EndPoint  string            `json:"end_point"`  // SMS服务的地址，默认为（https://sms.aliyuncs.com）
	AccessId  string            `json:"access_id"`  // 访问SMS服务的accessid，通过官方网站申请或通过管理员获取
	AccessKey string            `json:"access_key"` // 访问SMS服务的accesskey，通过官方网站申请或通过管理员获取
	Signs     map[string]string `json:"signs"`      // 应用对应的签名
}

func (c *SrvNotifierConfig) Import() error {
	var err error

	err = c.PublicConfig.Import()
	if err != nil {
		return err
	}

	err = c.DbConfig.Import("srv_notifier")
	if err != nil {
		return err
	}

	c.RPCPort, err = ConsulClient.Key("srv_notifier/rpc_port", ":11121")
	if err != nil {
		return err
	}
	//ConsulClient.KeyWatch("user/test", &c.RPCPort)

	c.ServerAliasName, err = ConsulClient.Key("srv_notifier/server_alias_name", "srv-notifier")
	if err != nil {
		return err
	}

	smsStr, err := ConsulClient.Key("srv_notifier/sms", "")
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(smsStr), c.Sms)
	if err != nil {
		return err
	}

	log.Debugf("SrvNotifier config import success! [%+v]", *c)
	return nil
}
