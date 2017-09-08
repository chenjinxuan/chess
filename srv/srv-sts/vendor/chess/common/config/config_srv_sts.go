package config

import (
	. "chess/common/consul"
	"chess/common/log"
)

var SrvSts = new(SrvStsConfig)
var CSts *SrvStsConfig

type SrvStsConfig struct {
	PublicConfig
	DbConfig

	ServerAliasName string //服务别名
	RPCPort         string
}

func (c *SrvStsConfig) Import() error {
	var err error

	err = c.DbConfig.Import("srv_sts")
	if err != nil {
		return err
	}

	c.RPCPort, err = ConsulClient.Key("srv_sts/rpc_port", ":11288")
	if err != nil {
		return err
	}
	//ConsulClient.KeyWatch("user/test", &c.RPCPort)

	c.ServerAliasName, err = ConsulClient.Key("srv_sts/server_alias_name", "srv-sts")
	if err != nil {
		return err
	}

	CSts = c
	log.Debugf("SrvSts config import success! [%+v]", *c)
	return nil
}
