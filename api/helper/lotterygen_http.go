package helper

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func GetBettingCode(u, key, token string, count uint) ([]string, error) {
	u2, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	q := u2.Query()
	q.Set("key", key)
	q.Set("token", token)
	q.Set("count", strconv.Itoa(int(count)))

	u2.RawQuery = q.Encode()

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(u2.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)

	if len(body) == 0 {
		return result, nil
	}

	data := string(body)
	result = strings.Split(data, ",")

	return result, nil
}
