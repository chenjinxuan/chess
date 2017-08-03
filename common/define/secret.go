package define

import (
	"treasure/config"
)

var (
	jwtSecretPrefix = "D3iURy2W"
	JwtSecret       string
	SmsSecret       = "&HI*^$H)"
)

func Init() {
	JwtSecret = jwtSecretPrefix + config.C.Secret.JwtSecret
}

func GetSecret(secret string) string {
	return jwtSecretPrefix + secret
}
