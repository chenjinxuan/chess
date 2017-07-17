package qqsdk

import (
	"fmt"
	"github.com/going/toolkit/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func extractDataByRegex(content, query string) (string, error) {
	rx := regexp.MustCompile(query)
	value := rx.FindStringSubmatch(content)

	if len(value) == 0 {
		return "", fmt.Errorf("正则表达式没有匹配到内容:(%s)", query)
	}

	return strings.TrimSpace(value[1]), nil
}

func qqGet(reqUrl string) ([]byte, error) {
	log.Println(reqUrl)
	var err error
	if resp, err := http.Get(reqUrl); err == nil {
		defer resp.Body.Close()

		if content, err := ioutil.ReadAll(resp.Body); err == nil {
			//先测试返回的是否是ReturnError
			if values, err := url.ParseQuery(string(content)); err == nil {
				code := values.Get("code")
				msg := values.Get("msg")
				if len(code) > 0 && len(msg) > 0 {
					return nil, fmt.Errorf("Request %s failed with code %s. Error message is '%s'.",
						reqUrl, code, msg)
				}
			}

			return content, nil
		}
	}

	return nil, err
}
