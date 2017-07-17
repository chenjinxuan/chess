package weibo

import (
	"github.com/orvice/weigo"
	"strings"
)

func NewApi(token, appId, appSecret string) *weigo.APIClient {
	api := weigo.NewAPIClient(appId, appSecret, "", "code")
	api.SetAccessToken(token, 0)
	return api
}

func NewWeigoUser() *weigo.User {
	user := new(weigo.User)
	return user
}

func GetUserShowByUid(token, uid, appId string, api *weigo.APIClient) (weigo.User, error) {
	var result weigo.User

	kws := map[string]interface{}{
		"uid":    uid,
		"source": appId,
	}
	err := api.GET_users_show(kws, &result)
	result.Avatar_large = strings.Replace(result.Avatar_large, "http://", "https://", -1)
	return result, err
}

func GetTokenInfo(token string, api *weigo.APIClient) (weigo.TokenInfo, error) {
	var result weigo.TokenInfo
	kws := map[string]interface{}{
		"access_token": token,
	}
	err := api.GET_token_info(kws, &result)
	return result, err
}
