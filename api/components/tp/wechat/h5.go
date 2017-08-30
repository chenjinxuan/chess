package wechat

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	// "net/url"
	"chess/api/log"
	"chess/api/redis"
	"strconv"
	"time"
)

type TicketData struct {
	AppId     string `json:"app_id"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"noncestr"`
	Signature string `json:"signature"`
}

// 获取公众号的access token
func (c *Client) GetGlobalAccessToken() (string, error) {
	var accessToken string
	redis := api_redis.Redis.Main
	key := "GlobalAccessToken" + c.appId
	isExist, err := api_redis.Redis.Task.Exists(key)
	if err != nil {
		return accessToken, err
	}
	// token 有缓存
	if isExist {
		// 从redis读取token
		str, err := redis.Get(key)
		return str, err
	}

	token, err := c.GetAccessToken()
	if err != nil {
		return accessToken, err
	}
	// store ticket
	err = redis.Setex(key, token.AccessToken, 7200)
	if err != nil {
		log.Log.Error(err)
	}
	return token.AccessToken, nil
}

func (c *Client) GetTicketStr() (string, error) {
	redis := api_redis.Redis.Main
	key := fmt.Sprintf("jsapiTicket-%s", c.appId)
	// 获取公众号
	var ticketStr string
	tokenStr, err := c.GetGlobalAccessToken()
	if err != nil {
		return ticketStr, err
	}
	// 换取 js api_ticket
	ticket, err := c.GetJsApiTicketByToken(tokenStr)
	if err != nil || ticket.Ticket == "" {
		return ticketStr, err
	}
	// store ticket
	err = redis.Setex(key, ticket.Ticket, 7200)
	if err != nil {
		log.Log.Error(err)
	}
	ticketStr = ticket.Ticket
	return ticketStr, nil
}

func (c *Client) GetTicket(openid, fromUrl string) (TicketData, error) {

	var timeStr, nonceStr string
	var data TicketData
	timeInt := time.Now().Unix()
	timeStr = strconv.Itoa(int(timeInt))
	//
	var ticketStr string
	redis := api_redis.Redis.Main
	//fromUrl := c.Request.Header.Get("Referer")

	// 生成nonceStr
	nonceStr = fmt.Sprintf("%x", md5.Sum([]byte(openid)))

	key := fmt.Sprintf("jsapiTicket-%s-%s", c.appId, openid)
	// 根据openid 检查ticket是否已经缓存
	isExist, err := redis.Exists(key)
	if err != nil {
		return data, err
	}
	// ticket 有缓存
	if isExist {
		// 从redis读取ticket
		str, err := redis.Get(key)
		if err != nil {
			return data, err
		}
		ticketStr = str
	} else {
		// 获取公众号
		tokenStr, err := c.GetGlobalAccessToken()
		if err != nil {
			return data, err
		}
		// 换取 js api_ticket
		ticket, err := c.GetJsApiTicketByToken(tokenStr)
		if err != nil || ticket.Ticket == "" {
			return data, err
		}
		// store ticket
		err = redis.Setex(key, ticket.Ticket, 7200)
		if err != nil {
			log.Log.Error(err)
		}
		ticketStr = ticket.Ticket
	}

	//var val url.Values
	/**
	val := url.Values{}
	val.Add("jsapi_ticket", ticketStr)
	val.Add("noncestr", nonceStr)
	val.Add("timestamp", timeStr)
	val.Add("url", fromUrl)
	val.Add("params", "value")
	string1 := val.Encode()
	**/
	string1 := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", ticketStr, nonceStr, timeStr, fromUrl)
	log.Log.Debug("string1:", string1)
	signature := sha1.Sum([]byte(string1))

	//gen ticket data
	ticketData := TicketData{
		AppId:     c.GetAppId(),
		Timestamp: timeInt,
		NonceStr:  nonceStr,
		Signature: fmt.Sprintf("%x", signature),
	}
	return ticketData, nil
}
