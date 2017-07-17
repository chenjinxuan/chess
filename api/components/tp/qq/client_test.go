package qqsdk

import (
	"testing"
)

var (
	client = new(Client)
	code   = "46DE48C2318348103BF17C88640C4642"

	appid        = "101282168"
	appsecret    = "7975396f2221ca5a0bf9a0df6e2a6e8f"
	redirect_url = "http://yyjinbao.com"
)

func TestClient(t *testing.T) {
	client.SetAppId(appid)
	client.SetSecret(appsecret)
	client.SetRedirectUrl(redirect_url)

	token, err := client.GetAccessToken(code, redirect_url)
	if err != nil {
		t.Error(err)
	}
	t.Log(token)

	openid, err := client.GetOpenId(token.AccessToken)
	if err != nil {
		t.Error(err)
	}
	user, err := client.GetUserInfo(token.AccessToken, openid)
	if err != nil {
		t.Error(err)
	}
	t.Log(user)
}
