package config

import (
	"chess/common/log"
)

var SrvUser = new(SrvUserConfig)

type SrvUserConfig struct {
	PublicConfig
	DbConfig

	ServerAliasName string //服务别名
	RPCPort         string
}

func (c *SrvUserConfig) Import() error {
	var err error

	err = c.PublicConfig.Import()
	if err != nil {
		return err
	}

	err = c.DbConfig.Import("srv_user")
	if err != nil {
		return err
	}

	c.RPCPort, err = ConsulClient.Key("srv_user/rpc_port", ":11121")
	if err != nil {
		return err
	}
	//ConsulClient.KeyWatch("user/test", &c.RPCPort)

	c.ServerAliasName, err = ConsulClient.Key("srv_user/server_alias_name", "srv-user")
	if err != nil {
		return err
	}

	log.Debugf("SrvUser config import success! [%+v]", *c)
	return nil
}
