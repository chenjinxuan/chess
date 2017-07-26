package config

import (
//"fmt"
. "chess/common/consul"
	"chess/common/log"
	"encoding/json"

    "github.com/gin-gonic/gin"
    "strings"
)

var Api = new(ApiConfig)

type ApiConfig struct {
	PublicConfig
	DbConfig
	Debug      bool          // debug 模式
	Port       string        // http 监听端口
	PostDesKey string        // post 数据加密 des key
	Secret          *Secret          `json:"secret"`
	Frontend        *Frontend        `json:"frontend"`
	Backend         *Backend         `json:"backend"`
	Tp              *Tp              `json:"tp"`
	User            *User            `json:"user"`
	Storage         *Storage         `json:"storage"`
	Sms             *Sms             `json:"sms"`
	Geoip           *Geoip           `json:"geoip"`
	Ip17mon         *Ip17mon         `json:"ip17mon"`
	IpReplace       []string         `json:"ip_replace"`
	Captcha         *Captcha         `json:"captcha"`
	Login           *Login           `json:"login"`
	RPC             *RPC             `json:"rpc"`
	AppstoreVerify  *AppstoreVerify  `json:"appstore_verify"`
	TokenInfoCheck  *TokenInfoCheck  `json:"token_info_check"`

	Filter   *FilterSet `json:"filter_set"`
	UserInit *UserInit  `json:"user_init"`

	Pay         []map[string]interface{}       `json:"pay"`
	Feedback    Feedback                       `json:"feedback"`
	Steer       *Steer `json:"steer"`

}
var (
    C  *ApiConfig
    Cs map[string]*ApiConfig
)

type Secret struct {
    JwtSecret string `json:"jwt_secret"`
}

type RPC struct {
    Timeout    int            `json:"timeout"`
    LotteryGen *RPCLotteryGen `json:"lottery_gen"`
    CouponGen  *RPCCouponGen  `json:"coupon_gen"`
    Notifier   *RPCCouponGen  `json:"notifier"`
    Runner     *RPCCouponGen  `json:"runner"`
}

type RPCLotteryGen struct {
    Servers []string `json:"servers"`
}

type RPCCouponGen struct {
    Server string `json:"server"`
}

type Frontend struct {
    CouponJumpGoodsId    int    `json:"coupon_jump_goods_id"`
    AllowCrossDomainHost string `json:"allow_cross_domain_host"`
}

type Backend struct {
    LotteryUrl            string            `json:"lottery_url"`
    Keyword               map[string]string `json:"keyword"`
    ImageDomainPrefix     string            `json:"image_domain_prefix"`
    ImageMainSuffix       string            `json:"image_main_suffix"`
    ImageIntroBigSuffix   string            `json:"image_intro_big_suffix"`
    ImageIntroSuffix      string            `json:"image_intro_suffix"`
    ImageDescrSuffix      string            `json:"image_descr_suffix"`
    ImageBannerSuffix     string            `json:"image_banner_suffix"`
    ImageShareSuffix      string            `json:"image_share_suffix"`
    BettingCountLimit     int               `json:"betting_count_limit"`
    GoodsSpreadCategoryId []int             `json:"goods_spread_category_id"`
    CouponNotifyChannel   string            `json:"coupon_notify_channel"`
    ParamsDesKey          string            `json:"params_des_key"`
}

type User struct {
    DefaultAvatar   string   `json:"default_avatar"`
    AvatarHostLimit []string `json:"avatar_host_limit"`
    AvatarExtLimit  []string `json:"avatar_ext_limit"`
}

type Storage struct {
    QiniuAccessKey    string `json:"qiniu_access_key"`
    QiniuSecretKey    string `json:"qiniu_secret_key"`
    QiniuAvatarBucket string `json:"qiniu_avatar_bucket"`
    QiniuShareBucket  string `json:"qiniu_share_bucket"`
    QiniuShareUrl     string `json:"qiniu_share_url"`
    QiniuAvatarUrl    string `json:"qiniu_avatar_url"`
    CallbackUrl       string `json:"callback_url"`
}

type Sms struct {
    SendLimit           int    `json:"10minute_send_limit"`
    CheckLimit          int    `json:"10minute_check_limit"`
    Time                int64  `json:"time"`
    MobileNumberRegular string `json:"mobile_number_regular"`
    Source              string `json:"source"`
    ContentPrefix       string `json:"content_perfix"`
}

type Geoip struct {
    Lang string `json:"lang"`
    Path string `json:"path"`
}

type Ip17mon struct {
    Path string `json:"path"`
}

type Captcha struct {
    ExpireTime int64 `json:"expire_time"`
}

type Login struct {
    ShowCaptchaCount int   `json:"show_captcha_count"`
    LimitTime        int64 `json:"limit_time"`
    FailLimit        int   `json:"fail_limit"`
    TokenExpire      int64 `json:"token_expire_time"`
}

type AppstoreVerify struct {
    GoodsFilter *AppstoreVerifyGoodsFilter `json:"goods_filter"`
}

type AppstoreVerifyGoodsFilter struct {
    Enable  bool  `json:"enable"`
    Version int   `json:"version"`
    List    []int `json:"list"`
}

type FilterSet struct {
    Set     []string `json:"set"`
    IsAsync bool     `json:"is_async"`
    Action  string   `json:"action"`
}

type Feedback struct {
    IssueLength   int `json:"issue_length"`
    ContactLength int `json:"contact_length"`
}

type Shiwan struct {
    ApiUrl string `json:"api_url"`
    DesKey string `json:"deskey"`
    DesIv  string `json:"desiv"`
}
type UserInit struct {
    Enable  bool            `json:"enable"`
    Balance int             `json:"balance"`
    Filter  *UserInitFilter `json:"filter"`
}

type UserInitFilter struct {
    Enable bool                `json:"enable"`
    White  map[string][]string `json:"white"`
}

type TokenInfoCheck struct {
    Android TokenInfoCheckDetail `json:"android"`
    Ios     TokenInfoCheckDetail `json:"ios"`
}

type TokenInfoCheckDetail struct {
    Check bool `json:"check"`
    Min   int  `json:"min"`
    Max   int  `json:"max"`
}

type Base struct {
    Mod     string `json:"mod"`
    Key     string `json:"key"`
    ProdUrl string `json:"prod_url"`
}

type Steer struct {
    RealBetRate float64 `json:"real_bet_rate"`
    GoodsPrice  int     `json:"goods_price"`
}

type ApiRpcConfig struct {
	Notifier string `json:"notifier"`
}

func (c *ApiConfig) Import() error {
	var err error

	err = c.PublicConfig.Import()
	if err != nil {
		return err
	}

	err = c.DbConfig.Import("api")
	if err != nil {
		return err
	}

	c.Debug, err = ConsulClient.KeyBool("api/debug", true)
	if err != nil {
		return err
	}
	ConsulClient.KeyBoolWatch("api/debug", &c.Debug)

	c.Port, err = ConsulClient.Key("api/port", ":8888")
	if err != nil {
		return err
	}
	ConsulClient.KeyWatch("api/port", &c.Port)

	c.PostDesKey, err = ConsulClient.Key("api/post_des_key", "XQ1R1%8f")
	if err != nil {
		return err
	}
	ConsulClient.KeyWatch("api/post_des_key", &c.PostDesKey)


        //defaultStr,err:=ConsulClient.Key("api/default","")
	//if err != nil {
	//    return err
	//}
        //err = json.Unmarshal([]byte(defaultStr),&c)
	//if err != nil {
	//    return err
	//}
        log.Debugf("Api config import success! [%+v]", *c)
	return nil
}
//
func SetContextConfig() gin.HandlerFunc {
    return func(c *gin.Context) {
	from := strings.ToLower(c.Query("from"))
	subfrom := strings.ToLower(c.Query("subfrom"))


	config := Get(from)

	// web内嵌页
	if subfrom == "web" {
	    configBytes, _ := json.Marshal(config)
	    subConfig := new(ApiConfig)
	    err := json.Unmarshal(configBytes, subConfig)
	    if err != nil {
		subConfig = config
	    }

	    webConfig := Get("web")
	    subConfig.Backend.ParamsDesKey = webConfig.Backend.ParamsDesKey
	    c.Set("config", subConfig)
	} else {
	    c.Set("config", config)
	}
    }
}
func Get(key string) *ApiConfig {
    if conf, ok := Cs[key]; ok {
	return conf
    }
    return C
}
func InitConfig() {

    defaultStr,err:=ConsulClient.Key("api/default","")
    //if err != nil {
	//return err
    //}
    err = json.Unmarshal([]byte(defaultStr),&C)
    //if err != nil {
	//return err
    //}
    diffStr ,err := ConsulClient.KeyList("api/diff")
    Cs = make(map[string]*ApiConfig)
    defaultBytes, err := json.Marshal(C)
    if err != nil {
	panic(err)
    }

    for k, v := range diffStr {
	if k == "default" {
	    continue
	}

	//valueBytes, err := json.Marshal(v)
	//if err != nil {
	//    panic(err)
	//}
	conf := new(ApiConfig)
	err = json.Unmarshal(defaultBytes, conf)
	if err != nil {
	    panic(err)
	}
	err = json.Unmarshal(v, conf)
	if err != nil {
	    panic(err)
	}
	Cs[k] = conf
    }

}
