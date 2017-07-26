package config

type TpApp struct {
    AppId       string `json:"appid"`
    AppSecret   string `json:"appsecret"`
    RedirectUrl string `json:"redirect_url"`
}

type Tp struct {
    QQ       TpApp `json:"qq"`
    H5QQ     TpApp `json:"h5_qq"`
    Weibo    TpApp `json:"weibo"`
    Wechat   TpApp `json:"wechat"`
    H5Wechat TpApp `json:"h5_wechat"`
}
