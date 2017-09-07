package wechat

import (
	"chess/common/config"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	//"treasure/log"
	"chess/common/log"
	"fmt"
	"strings"
)

var (
	WxClient   = new(Client)
	H5WxClient = new(Client)
)

type Client struct {
	appId     string
	appSecret string
}

func NewDefaultWxCliant() *Client {
	//conf := config.C.Tp.Wechat
	client := new(Client)
	//if config.C != nil {
	//	log.Log.Info("config is nil")
	//}
	// client.SetAppId(config.C.Tp.Wechat.AppId)
	//client.SetAppSecret(conf.AppSecret)
	return client
}

func NewClient(appId, appSecret string) *Client {
	client := new(Client)
	client.SetAppId(appId)
	client.SetAppSecret(appSecret)
	return client
}

func Init() {
	conf := config.C.Tp.Wechat
	WxClient.SetAppId(conf.AppId)
	WxClient.SetAppSecret(conf.AppSecret)

	h5Conf := config.C.Tp.H5Wechat
	H5WxClient.SetAppId(h5Conf.AppId)
	H5WxClient.SetAppSecret(h5Conf.AppSecret)

}

func (c *Client) SetAppId(aid string) {
	c.appId = aid
}

func (c *Client) SetAppSecret(secret string) {
	c.appSecret = secret
}

func (c *Client) GetAppId() string {
	return c.appId
}

func (c *Client) ExchangeTokenURL(code string) string {
	urlStr := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + url.QueryEscape(c.appId) +
		"&secret=" + url.QueryEscape(c.appSecret) +
		"&grant_type=authorization_code&code=" + url.QueryEscape(code)
	//log.Log.Debug(urlStr)
	return urlStr
}

// 公众号token
func (c *Client) GetAccessTokenURL() string {
	urlStr := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + url.QueryEscape(c.appId) +
		"&secret=" + url.QueryEscape(c.appSecret)
	return urlStr
}

func (c *Client) CallGet(urlStr string) ([]byte, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	//log.Log.Debug(body)
	return body, nil
}

func (c *Client) GetTokenByCode(code string) (Token, error) {
	var tk Token
	urlStr := c.ExchangeTokenURL(code)
	resp, err := c.CallGet(urlStr)
	if err != nil {
		return tk, err
	}
	err = json.Unmarshal([]byte(resp), &tk)
	log.Debug(string(resp))
	log.Debug(fmt.Sprintf("wechat-login-debug: - token: %v, union id :  %v", tk.AccessToken, tk.UnionId))
	if tk.AccessToken == "" {
		err = errors.New("cant get token by this code")
	}
	return tk, err
}

func (cfg *Client) UserInfoURL(accessToken, openId, lang string) string {
	if lang == "" {
		return "https://api.weixin.qq.com/sns/userinfo?access_token=" + url.QueryEscape(accessToken) +
			"&openid=" + url.QueryEscape(openId)
	}
	return "https://api.weixin.qq.com/sns/userinfo?access_token=" + url.QueryEscape(accessToken) +
		"&openid=" + url.QueryEscape(openId) +
		"&lang=" + url.QueryEscape(lang)
}

func (c *Client) GetUserInfoByToken(tk Token) (UserInfo, error) {
	var user UserInfo
	urlStr := c.UserInfoURL(tk.AccessToken, tk.OpenId, "")
	resp, err := c.CallGet(urlStr)
	if err != nil {
		return user, err
	}
	err = json.Unmarshal([]byte(resp), &user)
	if tk.AccessToken == "" {
		err = errors.New("cant get userinfo by token")
	}
	user.HeadImageURL = strings.Replace(user.HeadImageURL, "http://", "https://", -1)
	return user, err
}

func (cfg *Client) JsApiTicketURL(accessToken string) string {
	return "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=" + url.QueryEscape(accessToken) + "&type=jsapi"
}

func (c *Client) GetJsApiTicketByToken(token string) (JsApiTicket, error) {
	var ticket JsApiTicket
	urlStr := c.JsApiTicketURL(token)
	resp, err := c.CallGet(urlStr)
	if err != nil {
		return ticket, err
	}
	log.Debug("js api ticket get resp debug: ", string(resp))
	err = json.Unmarshal([]byte(resp), &ticket)
	return ticket, err
}

func (c *Client) GetAccessToken() (AccessToken, error) {
	var accessToken AccessToken
	urlStr := c.GetAccessTokenURL()
	resp, err := c.CallGet(urlStr)
	if err != nil {
		return accessToken, err
	}
	log.Debug("get access token debug: ", string(resp))
	err = json.Unmarshal([]byte(resp), &accessToken)
	return accessToken, err
}
