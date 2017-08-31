package config

import (
	. "chess/common/consul"
	"chess/common/log"
)

var SrvCentre = new(SrvCentreConfig)

type SrvCentreConfig struct {
	PublicConfig
	DbConfig

	ServerAliasName string //服务别名
}

func (c *SrvCentreConfig) Import() error {
	var err error

	err = c.PublicConfig.Import()
	if err != nil {
		return err
	}

	err = c.DbConfig.Import("srv_centre")
	if err != nil {
		return err
	}

	c.ServerAliasName, err = ConsulClient.Key("srv_centre/server_alias_name", "centre")
	if err != nil {
		return err
	}

	log.Debugf("SrvCentre config import success! [%+v]", *c)
	return nil
}
