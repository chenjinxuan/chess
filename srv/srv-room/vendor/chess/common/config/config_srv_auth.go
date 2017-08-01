package config

import (
	. "chess/common/consul"
	"chess/common/log"
)

var SrvAuth = new(SrvAuthConfig)

type SrvAuthConfig struct {
	PublicConfig
	DbConfig

	ServerAliasName string //服务别名
	RPCPort         string
}

func (c *SrvAuthConfig) Import() error {
	var err error

	err = c.PublicConfig.Import()
	if err != nil {
		return err
	}

	err = c.DbConfig.Import("srv_auth")
	if err != nil {
		return err
	}

	c.RPCPort, err = ConsulClient.Key("srv_auth/rpc_port", ":11121")
	if err != nil {
		return err
	}
	//ConsulClient.KeyWatch("user/test", &c.RPCPort)

	c.ServerAliasName, err = ConsulClient.Key("srv_auth/server_alias_name", "srv-user")
	if err != nil {
		return err
	}

	log.Debugf("SrvAuth config import success! [%+v]", *c)
	return nil
}
