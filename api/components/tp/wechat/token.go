package wechat

type Token struct {
	AccessToken  string `json:"access_token"`
	CreateAt     int64  `json:"create_at"`  // 创建时间, unixtime, 分布式系统要求时间同步, 建议使用 NTP
	ExpiresIn    int64  `json:"expires_in"` // 超时时间, seconds
	RefreshToken string `json:"refresh_token"`

	OpenId  string   `json:"openid"`
	UnionId string   `json:"unionid,omitempty"`
	Scopes  []string `json:"scopes,omitempty"` // 用户授权的作用域
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
