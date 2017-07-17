package weibo

type TokenInfo struct {
	UID       int    `json:"uid"`
	AppKey    string `json:"appkey"`
	Scope     string `json:"scope"`
	CreatedAt int    `json:"create_at"`
	ExpireIn  int    `json:"expire_in"`
}
