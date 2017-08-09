package helper

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func ClientIP(c *gin.Context) string {
	var clientIp string
	var err error
	clientIp, _, err = net.SplitHostPort(c.ClientIP())
	if err != nil {
		clientIp = c.ClientIP()
	}

	return clientIp
}

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

func HttpGetWithHeader(urlStr string, headers map[string]string) (string, error) {
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

func HttpPost(url string, reqBodyBytes []byte) ([]byte, error) {
	reqBody := bytes.NewBuffer([]byte(reqBodyBytes))

	var body []byte
	timeout := time.Duration(10) * time.Second

	transport := &http.Transport{
		ResponseHeaderTimeout: timeout,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		},
		DisableKeepAlives: true,
	}

	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Post(url, "application/json;charset=utf-8", reqBody)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	return body, nil
}

func SetCookie(w gin.ResponseWriter, domain, name, val string, expire int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    val,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Domain:   domain,
		Expires:  time.Now().Add(time.Duration(expire) * time.Second),
	}
	http.SetCookie(w, cookie)
}

func EchoResult(c *gin.Context, res interface{}) {
	callBack := c.Query("callback")
	if callBack == "" {
		c.JSON(http.StatusOK, res)
	} else {
		resBytes, _ := json.Marshal(res)
		resStr := string(resBytes)
		c.String(http.StatusOK, "%s(%s)", callBack, resStr)
	}
}
