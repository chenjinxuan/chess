package qqsdk

type UserInfo struct {
	Ret                int    //返回值
	Msg                string //错误消息
	Nickname           string //用户在QQ空间的昵称。
	Figureurl          string //大小为30×30像素的QQ空间头像URL。
	Figureurl_1        string //大小为50×50像素的QQ空间头像URL。
	Figureurl_2        string //大小为100×100像素的QQ空间头像URL。
	Figureurl_qq_1     string //大小为40×40像素的QQ头像URL。
	Figureurl_qq_2     string //大小为100×100像素的QQ头像URL。需要注意，不是所有的用户都拥有QQ的100x100的头像，但40x40像素则是一定会有。
	Gender             string //性别。 如果获取不到则默认返回"男"
	Is_yellow_vip      string //标识用户是否为黄钻用户（0：不是；1：是）。
	Vip                string //标识用户是否为黄钻用户（0：不是；1：是）
	Yellow_vip_level   string //黄钻等级
	Level              string //黄钻等级
	Is_yellow_year_vip string //标识是否为年费黄钻用户（0：不是； 1：是）
}

// access_token, expires_in, refresh_token string, err error
type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type ReturnDataForAddT struct {
	Id   int
	Time int64
}

type ReturnForAddT struct {
	Ret     int
	Msg     string
	Errcode int
	Data    ReturnDataForAddT
}
