package config

import (
	. "chess/common/consul"
	"chess/common/log"
)

var SrvRoom = new(SrvRoomConfig)

type SrvRoomConfig struct {
	PublicConfig
	DbConfig

	ServerAliasName string //服务别名
}

func (c *SrvRoomConfig) Import() error {
	var err error

	err = c.PublicConfig.Import()
	if err != nil {
		return err
	}

	err = c.DbConfig.Import("srv_room")
	if err != nil {
		return err
	}

	//ConsulClient.KeyWatch("user/test", &c.RPCPort)

	c.ServerAliasName, err = ConsulClient.Key("srv_room/server_alias_name", "room")
	if err != nil {
		return err
	}

	log.Debugf("SrvUser config import success! [%+v]", *c)
	return nil
}
