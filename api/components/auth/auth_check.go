package auth

import (
	"errors"
	"strconv"
	"chess/common/config"
        "chess/api/redis"
)

func GetFailLoginKey(mobile string) string {
	return "mobile-login-fail-" + mobile
}

func GetFailLoginCount(mobile string) (int, error) {
	key := GetFailLoginKey(mobile)
	client := redis.Redis.Login
	isExist, err := client.Exists(key)
	if err != nil {
		return 0, errors.New("system error")
	}
	// 检查key是否存在
	if !isExist {
		return 0, nil
	}

	// 获取当前count
	countStr, err := client.Get(key)
	if err != nil {
		return 0, errors.New("system error")
	}
	count, err := strconv.Atoi(countStr)
	return count, nil
}

func FailCountPlusOne(mobile string) error {
	key := GetFailLoginKey(mobile)
	client := redis.Redis.Login
	isExist, err := client.Exists(key)
	if err != nil {
		return errors.New("system error")
	}
	if !isExist {
		err = client.Setex(key, "1" ,config.C.Login.LimitTime)
		return err
	}

	// 获取当前count
	countStr, err := client.Get(key)
	if err != nil {
		return errors.New("system error")
	}
	count, err := strconv.Atoi(countStr)
	c := count + 1
	err = client.Setex(key,strconv.Itoa(c), config.C.Login.LimitTime )
	if err != nil {
		return errors.New("system error,cant update count")
	}
	return nil
}
