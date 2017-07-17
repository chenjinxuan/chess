package qqsdk

import (
	"errors"
	"treasure/log"
)

func GetOpenId(accessToken string) (string, error) {
	reqUrl := UrlQQOAuth + "/me?access_token=" + accessToken

	var err error
	var respContent []byte

	if respContent, err = qqGet(reqUrl); err == nil {
		if openId, err := extractDataByRegex(string(respContent), `"openid":"(.*?)"`); err == nil {
			return openId, nil
		}
	}
	err = errors.New("cant get appid ")
	return "", err
}

func GetAppId(accessToken string) (string, error) {
	reqUrl := UrlQQOAuth + "/me?access_token=" + accessToken

	var err error
	var respContent []byte

	if respContent, err = qqGet(reqUrl); err == nil {
		if openId, err := extractDataByRegex(string(respContent), `"client_id":"(.*?)"`); err == nil {
			return openId, nil
		}
	}
	err = errors.New("cant get appid ")
	return "", err
}

func GetOpenIdAndUnionId(accessToken string) (openId string, unionId string, err error) {
	reqUrl := UrlQQOAuth + "/me?access_token=" + accessToken + "&unionid=1"

	var respContent []byte
	if respContent, err = qqGet(reqUrl); err == nil {
		log.Log.Debugf("qq get openid and unionid resp: %s", string(respContent))
		openId, err = extractDataByRegex(string(respContent), `"openid":"(.*?)"`)
		if err != nil {
			return
		}

		unionId, err = extractDataByRegex(string(respContent), `"unionid":"(.*?)"`)
		if err != nil {
			return
		}
		log.Log.Debugf("Get qq openid(%s) unionid(%s) via token(%s)", openId, unionId, accessToken)
		return
	}
	err = errors.New("cant get openid ")
	return
}
