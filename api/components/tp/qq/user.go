package qqsdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func GetUserInfo(accessToken, openId, appId string) (*UserInfo, error) {
	v := url.Values{}
	v.Add("access_token", accessToken)
	v.Add("oauth_consumer_key", appId)
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
