package captcha

import (
	"github.com/dchest/captcha"
	"chess/api/config"
        "chess/common/db"
        "chess/models"
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
	models.ChessRedis.Captcha.Setex(id, string(digits), config.C.Captcha.ExpireTime)
}

func (r *RedisStore) Set(id string, digits []byte) {
    models.ChessRedis.Captcha.Setex(id, string(digits),config.C.Captcha.ExpireTime)
}

func (r *RedisStore) Get(id string, clear bool) (digits []byte) {
	// digits = captcha.RandomDigits(4)
	str, err := models.ChessRedis.Captcha.Get(id)
	if err != nil {
		// @todo
		digits = captcha.RandomDigits(4)
	}
	digits = []byte(str)
	return
}

func (r *RedisStore) Del(id string) error {
	err := models.ChessRedis.Captcha.Del(id)
	return err
}
