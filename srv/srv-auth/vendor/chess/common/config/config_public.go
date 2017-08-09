package config

import (
	. "chess/common/consul"
)

var Public = new(PublicConfig)

type PublicConfig struct {
	QiniuAccessKey          string
	QiniuSecretKey          string
        SmsSendLimit            int
        SmsCheckLimit           int
        SmsTime                 int
        TokenSecret             string
}

func (c *PublicConfig) Import() error {
	var err error
	c.QiniuAccessKey, err = ConsulClient.Key("public/qiniu/access_key", "")
	if err != nil {
		return err
	}

	c.QiniuAccessKey, err = ConsulClient.Key("public/qiniu/secret_key", "")
	if err != nil {
		return err
	}

	c.SmsSendLimit, err = ConsulClient.KeyInt("public/sms_send", 10)
	if err != nil {
	    return err
	}
	c.SmsCheckLimit, err = ConsulClient.KeyInt("public/sms_check", 5)
	if err != nil {
	    return err
	}
	c.SmsTime, err = ConsulClient.KeyInt("public/sms_time", 600)
	if err != nil {
	    return err
	}
	c.TokenSecret, err = ConsulClient.Key("public/tokensecret", "rewrqwerq")
	if err != nil {
	    return err
	}
	return nil
}
