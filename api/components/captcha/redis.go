package captcha

import (
	"chess/api/redis"
	"chess/common/config"
	"chess/common/db"
	"github.com/dchest/captcha"
)

//func init() {
//	store := NewRedisStore()
//	captcha.SetCustomStore(store)
//}

var (
	Strore = new(RedisStore)
)

type RedisStore struct {
	Client *db.Redis
}

func NewRedisStore() captcha.Store {
	//r := new(RedisStore)
	//r.Client = databases.Redis.Captcha
	return Strore
}

func (r *RedisStore) SetRandom(id string) {
	digits := captcha.RandomDigits(4)
	api_redis.Redis.Captcha.Setex(id, string(digits), config.C.Captcha.ExpireTime)
}

func (r *RedisStore) Set(id string, digits []byte) {
	api_redis.Redis.Captcha.Setex(id, string(digits), config.C.Captcha.ExpireTime)
}

func (r *RedisStore) Get(id string, clear bool) (digits []byte) {
	// digits = captcha.RandomDigits(4)
	str, err := api_redis.Redis.Captcha.Get(id)
	if err != nil {
		// @todo
		digits = captcha.RandomDigits(4)
	}
	digits = []byte(str)
	return
}

func (r *RedisStore) Del(id string) error {
	err := api_redis.Redis.Captcha.Del(id)
	return err
}
