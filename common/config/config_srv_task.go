package config

import (
    . "chess/common/consul"
    "chess/common/log"
)

var SrvTask = new(SrvTaskConfig)
var CTask *SrvTaskConfig

type SrvTaskConfig struct {
    PublicConfig
    DbConfig

    ServerAliasName string //服务别名
    RPCPort         string
}

func (c *SrvTaskConfig) Import() error {
    var err error

    err = c.DbConfig.Import("srv_task")
    if err != nil {
	return err
    }

    c.RPCPort, err = ConsulClient.Key("srv_task/rpc_port", ":11188")
    if err != nil {
	return err
    }
    //ConsulClient.KeyWatch("user/test", &c.RPCPort)

    c.ServerAliasName, err = ConsulClient.Key("srv_task/server_alias_name", "srv-user")
    if err != nil {
	return err
    }

    CTask = c
    log.Debugf("SrvAuth config import success! [%+v]", *c)
    return nil
}
