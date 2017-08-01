package config

import (
	. "chess/common/consul"
	"chess/common/log"
        "encoding/json"
)

var SrvAuth = new(SrvAuthConfig)
var CAuth *SrvAuthConfig
type SrvAuthConfig struct {
	PublicConfig
	DbConfig

	ServerAliasName string //服务别名
	RPCPort         string
        Login           *Login           `json:"login"`
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
	defaultStr,err:=ConsulClient.Key("srv_auth/default","")
	if err != nil {
	return err
	}
    	err = json.Unmarshal([]byte(defaultStr),&c)
	if err != nil {
	    return err
	}
	CAuth=c
	log.Debugf("SrvAuth config import success! [%+v]", *c)
	return nil
}
