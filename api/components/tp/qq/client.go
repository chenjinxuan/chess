package qqsdk

import (
	"chess/api/log"
	"chess/common/config"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var (
	DefaultClient = new(Client)
	H5Client      = new(Client)

	ErrCantGetToken = errors.New("cant get access token")
)

type Client struct {
	appId       string
	appSecret   string
	redirectUrl string
}

func (c *Client) SetAppId(aid string) {
	c.appId = aid
}

func (c *Client) GetAppId() string {
	return c.appId
}

func (c *Client) GetSecret() string {
	return c.appSecret
}

func (c *Client) SetSecret(secret string) {
	c.appSecret = secret
}

func (c *Client) SetRedirectUrl(redirectUrl string) {
	c.redirectUrl = redirectUrl
}

func NewClient(appId, appSecret, redirectUrl string) *Client {
	client := new(Client)
	client.SetAppId(appId)
	client.SetSecret(appSecret)
	client.SetRedirectUrl(redirectUrl)
	return client
}

func InitQQClient() {
	// init default client
	defaultConfig := config.C.Tp.QQ
	DefaultClient.SetAppId(defaultConfig.AppId)
	DefaultClient.SetSecret(defaultConfig.AppSecret)

	h5Config := config.C.Tp.H5QQ
	H5Client.SetAppId(h5Config.AppId)
	H5Client.SetSecret(h5Config.AppSecret)
	H5Client.SetRedirectUrl(h5Config.RedirectUrl)

}

func (c *Client) GetAuthorizationCodeUrl(redirectUrl, state, scope string) string {
	v := url.Values{}
	v.Add("response_type", "code")
	v.Add("client_id", c.appId)
	v.Add("redirect_uri", redirectUrl)
	v.Add("state", state)
	v.Add("scope", scope)

	return UrlQQOAuth + "/authorize?" + v.Encode()
}

func (c *Client) GetAccessToken(authCode, redirectUrl string) (Token, error) {
	var token Token
	v := url.Values{}
	v.Add("grant_type", "authorization_code")
	v.Add("client_id", c.appId)
	v.Add("client_secret", c.appSecret)
	v.Add("code", authCode)
	v.Add("redirect_uri", c.redirectUrl)

	reqUrl := UrlQQOAuth + "/token?" + v.Encode()

	if respContent, err := qqGet(reqUrl); err == nil {
		log.Log.Debug(string(respContent))
		if values, err := url.ParseQuery(string(respContent)); err == nil {
			token.AccessToken = values.Get("access_token")
			token.ExpiresIn = values.Get("expires_in")
			token.RefreshToken = values.Get("refresh_token")
			//errorDescription := values.Get("error_description")
		}
		if token.AccessToken == "" {
			return token, ErrCantGetToken
		}
	}

	return token, nil
}

func (c *Client) RefreshToken(appId, appKey, refreshToken string) (access_token, expires_in, refresh_token string, err error) {
	v := url.Values{}
	v.Add("grant_type", "refresh_token")
	v.Add("client_id", c.appId)
	v.Add("client_secret", c.appSecret)
	v.Add("refresh_token", refreshToken)

	reqUrl := UrlQQOAuth + "/token?" + v.Encode()

	if respContent, err := qqGet(reqUrl); err == nil {
		if values, err := url.ParseQuery(string(respContent)); err == nil {
			access_token = values.Get("access_token")
			expires_in = values.Get("expires_in")
			refresh_token = values.Get("refresh_token")
		}
	}

	return
}

func (c *Client) GetUserInfo(accessToken, openId string) (*UserInfo, error) {
	v := url.Values{}
	v.Add("access_token", accessToken)
	v.Add("oauth_consumer_key", c.appId)
	v.Add("openid", openId)
	v.Add("format", "json")

	reqUrl := UrlQQ + "/user/get_user_info?" + v.Encode()
	//reqUrl := UrlQQ + "/user/get_simple_userinfo?" + v.Encode()

	var err error
	var resp *http.Response

	if resp, err = http.Get(reqUrl); err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()

		var userInfo UserInfo
		json.NewDecoder(resp.Body).Decode(&userInfo)

		if userInfo.Ret != 0 {
			return nil, fmt.Errorf("Get %s failed. Ret:%d Msg:%s", reqUrl, userInfo.Ret, userInfo.Msg)
		}

		if userInfo.Nickname == "" {
			userInfo.Nickname = "QQ用户"
		}
		userInfo.Figureurl_qq_2 = strings.Replace(userInfo.Figureurl_qq_2, "http://", "https://", -1)

		return &userInfo, nil
	}

	return nil, fmt.Errorf("GetUserInfo failed with status code %d", resp.StatusCode)
}

func (c *Client) GetOpenId(accessToken string) (string, error) {
	reqUrl := UrlQQOAuth + "/me?access_token=" + accessToken

	var err error
	var respContent []byte

	if respContent, err = qqGet(reqUrl); err == nil {
		if openId, err := extractDataByRegex(string(respContent), `"openid":"(.*?)"`); err == nil {
			return openId, nil
		}
	}
	err = errors.New("cant get openid ")
	return "", err
}

func (c *Client) GetOpenIdAndUnionId(accessToken string) (openId string, unionId string, err error) {
	reqUrl := UrlQQOAuth + "/me?access_token=" + accessToken + "&unionid=1"

	var respContent []byte
	if respContent, err = qqGet(reqUrl); err == nil {
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
