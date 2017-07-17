package qqsdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func AddT(accessToken, appId, openId, content string) (*ReturnForAddT, error) {
	v := url.Values{}
	v.Add("access_token", accessToken)
	v.Add("oauth_consumer_key", appId)
	v.Add("openid", openId)
	v.Add("format", "json")
	v.Add("content", content)

	reqUrl := UrlQQ + "/t/add_t"

	var err error

	if resp, err := http.PostForm(reqUrl, v); err == nil {
		defer resp.Body.Close()

		var ret ReturnForAddT
		if err = json.NewDecoder(resp.Body).Decode(&ret); err == nil && ret.Ret != 0 {
			return nil, fmt.Errorf("Request %s failed with ret %d, msg %s.",
				reqUrl, ret.Ret, ret.Msg)
		}

		return &ret, nil
	}

	return nil, err
}
