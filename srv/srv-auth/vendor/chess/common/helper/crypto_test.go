package helper

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

/**
func Test_Encrypt(t *testing.T) {
	test := DesEncrypt("testkey1", `{"shiwan_id":"test","mobile":"18805050324","gender":1,"nickname":"test","avatar":"test","idfa":"testidfa","token":"11111","timestamp":11}`)
	t.Error(test)
	t.Error(DesDecrypt("testkey1", test))
}**/

func Test_Decrypt(t *testing.T) {
	source := "Adpz5yDGJy5AeHPgn2r8VKFMjuJ43e6rmSWoejaOnrUxQ17TvsNKbezcEA0We0dLIy7c8kz3qbhM7VvT5JnR55SgnRPRerz71FE5GL6jCygxfXbcgfE1aYTkVjb7HfLQh4fN_@_jXMJ5UxcTxM_*_2h54XLNhlLJ2d8JQdDj9_*_SERI6YD1NhtOuBGUpnQiIKfCK9kkEqzI7XNJGBoEyRqHNWf2S9jrF6kM_*_UeZLSfKfbamE_@_Ym9Hv2YQqmAdnqcbApw6NoWmxQIr_*_pXDwZP9orEtQGbDi2z0inD59mEd363KtsnlPxxPQUAcmpEiHojIPt_@_9FjCZcOAFO9LsPfUK84VLerktORgpdqrJncw_@_dYmQNSNUXGWPDxE7gh4pemN41SLNiPRsK8fxg8XfXTwuJNcBm8b3KnMm5lgYv95xefVocCXLvSrl9JZ7XT5PxpUk8yyUJ1Ksd5fTSYflyJYnuOdsRuYf19D2x9qgPm8n1sgcFbvBhgqYESjj9g=="
	source = strings.Replace(source, "_@_", "+", -1)
	source = strings.Replace(source, "_*_", "/", -1)
	_source, _ := base64.StdEncoding.DecodeString(source)
	__source := hex.EncodeToString(_source)
	res := DesDecrypt("ac68!3#1", "55%g7z!@", __source)
	t.Log(res)
}
