package helper

import (
	"github.com/gin-gonic/gin"
	"net"
	"strings"
	"treasure/config"
	"unicode"
)

func ClientIP(c *gin.Context) string {
	var clientIp string
	var err error
	clientIp, _, err = net.SplitHostPort(c.ClientIP())
	if err != nil {
		clientIp = c.ClientIP()
	}

	for _, ipReplace := range config.C.IpReplace {
		if ipReplace == clientIp {
			// replace the ip to beijing
			f := func(c rune) bool {
				return !unicode.IsNumber(c)
			}
			ipSplit := strings.FieldsFunc(clientIp, f)
			if len(ipSplit) < 4 {
				break
			}
			ipSplit[0] = "121"
			ipSplit[1] = "4"
			ipSplit[2] = "56"

			clientIp = strings.Join(ipSplit, ".")
			break
		}
	}

	return clientIp
}
