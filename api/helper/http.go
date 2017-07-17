package helper

import (
	"code.google.com/p/mahonia"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

func HttpGet(urlStr string) (string, error) {
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

func CovertGbkToUtf8(input string) (string, error) {
	dec := mahonia.NewDecoder("gbk")

	if output, ok := dec.ConvertStringOK(input); ok {
		return ReplaceGbk(output), nil
	}

	return "", errors.New("Covert Gbk To Utf8 Fail.")
}

func ReplaceGbk(input string) string {
	reg := regexp.MustCompile(`(?:gb2312|GB2312)`)
	output := reg.ReplaceAllString(input, "utf-8")
	return output
}

func HttpGetPro(urlStr string, headers map[string]string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
