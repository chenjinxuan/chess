package captcha

import (
	"fmt"
	"github.com/dchest/captcha"
	"github.com/satori/go.uuid"
	"io"
)

func init() {
	store := NewRedisStore()
	captcha.SetCustomStore(store)
}

func WriteImage(w io.Writer, id string, width, height int) error {
	return captcha.WriteImage(w, id, width, height)
}

func Verify(id, digits string) bool {
	isOk := captcha.VerifyString(id, digits)
	if isOk {
		Strore.SetRandom(id)
	}
	return isOk
}

func GenCaptcha() string {
	id := uuid.NewV4().String()
	Strore.SetRandom(id)
	return id
}

func GetMobileLoginId(mobile string) string {
	return fmt.Sprintf("mobile-captcha-%s", mobile)
}

func GenCaptchaByMobile(mobile string) string {
	id := GetMobileLoginId(mobile)
	Strore.SetRandom(id)
	return id
}

func VerifyByMobileLogin(mobile string, digits string) bool {
	id := GetMobileLoginId(mobile)
	return Verify(id, digits)
}
