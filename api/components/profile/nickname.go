package profile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type NicknameResponse struct {
	Code    int          `json:"Code"`
	Message string       `json:"Message"`
	Data    NicknameData `json:"Data"`
}

type NicknameData struct {
	Result   string     `json:"Result"`
	Snippets []Snippets `json:"Snippets"`
}

type Snippets struct {
}

func getCheckUrl(nickname string) string {
	return fmt.Sprintf("http://other.wcf.tongbu.com/infos/FigLeaf.ashx?Match&content=%s&matchLevel=0&filterLevel=0", nickname)
}

func httpGet(urlStr string) (string, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		// handle error
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return "", err
	}

	return string(body), nil
}

// 返回 处理过的昵称, 是否含有敏感词, err
func CheckNickname(nickname string) (string, bool, error) {
	urlStr := getCheckUrl(url.QueryEscape(nickname))
	resp, err := httpGet(urlStr)
	if err != nil {
		return "", false, err
	}
	fmt.Println("CheckNickname Return:" + resp)
	var respData NicknameResponse
	err = json.Unmarshal([]byte(resp), &respData)
	if err != nil {
		return "", false, err
	}
	if len(respData.Data.Snippets) != 0 {
		return respData.Data.Result, false, nil
	}
	return respData.Data.Result, true, nil
}
