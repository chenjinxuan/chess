package captcha

import (
	"github.com/dchest/captcha"
	"chess/api/config"
	"chess/api/databases"
)

func init() {
	store := NewRedisStore()
	captcha.SetCustomStore(store)
}

var (
	Strore = new(RedisStore)
)

type RedisStore struct {
	Client *databases.RedisPool
}

func NewRedisStore() captcha.Store {
	//r := new(RedisStore)
	//r.Client = databases.Redis.Captcha
	return Strore
}

func (r *RedisStore) SetRandom(id string) {
	digits := captcha.RandomDigits(4)
	databases.Redis.Captcha.Setex(id, string(digits), config.C.Captcha.ExpireTime)
}

func (r *RedisStore) Set(id string, digits []byte) {
	databases.Redis.Captcha.Setex(id, string(digits), config.C.Captcha.ExpireTime)
}

func (r *RedisStore) Get(id string, clear bool) (digits []byte) {
	// digits = captcha.RandomDigits(4)
	str, err := databases.Redis.Captcha.Get(id)
	if err != nil {
		// @todo
		digits = captcha.RandomDigits(4)
	}
	digits = []byte(str)
	return
}

func (r *RedisStore) Del(id string) error {
	err := databases.Redis.Captcha.Del(id)
	return err
}
