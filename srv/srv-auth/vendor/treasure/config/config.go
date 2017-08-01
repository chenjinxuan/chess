package config

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
	"time"
)

var (
	C  *Config
	Cs map[string]*Config

	NotFound    = errors.New("Config file not found")
	FormatError = errors.New("Config format error")
)

type Databases struct {
	MySQL *MySQL `json:"mysql"`
	Mongo *Mongo `json:"mongo"`
	Redis *Redis `json:"redis"`
}

type MySQL struct {
	Setting *MySQLSetting           `json:"setting"`
	Server  map[string]*MySQLServer `json:"server"`
}

type MySQLSetting struct {
}

type MySQLServer struct {
	Db       string `json:"db"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Mongo struct {
	Setting *MongoSetting `json:"setting"`
	Server  *MongoServer  `json:"server"`
}

type MongoSetting struct {
	DialTimeout time.Duration `json:"dial_timeout"`
}

type MongoServer struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Redis struct {
	Setting *RedisSetting           `json:"setting"`
	Server  map[string]*RedisServer `json:"server"`
}

type RedisSetting struct {
	MaxIdle     int           `json:"max_idle"`
	MaxActive   int           `json:"max_active"`
	IdleTimeout time.Duration `json:"idle_timeout"`
	WaitIdle    bool          `json:"wait_idle"`
}

type RedisServer struct {
	Db       int    `json:"db"`
	Host     string `json:"host"`
	Password string `json:"password"`
}

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

type Config struct {
	Databases       *Databases       `json:"databases"`
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
	CouponEventList *CouponEventList `json:"coupon_event"`
	TokenInfoCheck  *TokenInfoCheck  `json:"token_info_check"`

	Filter   *FilterSet `json:"filter_set"`
	Shiwan   Shiwan     `json:"shiwan"`
	UserInit *UserInit  `json:"user_init"`

	Pay         []map[string]interface{}       `json:"pay"`
	Broker      Broker                         `json:"broker"`
	Feedback    Feedback                       `json:"feedback"`
	Imitation   Imitation                      `json:"imitation"`
	Steer       *Steer                         `json:"steer"`
	MidAutumn   *MidAutumn                     `json:"mid-autumn"`
	NationalDay map[string]MidAutumnPoolConfig `json:"national-day"`
	Hunt        *Hunt                          `json:"hunt"`
}

type Steer struct {
	RealBetRate float64 `json:"real_bet_rate"`
	GoodsPrice  int     `json:"goods_price"`
}

func Get(key string) *Config {
	if config, ok := Cs[key]; ok {
		return config
	}

	return C
}

func New(path string) (*Config, error) {
	config := new(Config)

	// get config from file
	err := config.Load(path)
	if err != nil {
		return config, err
	}

	return config, nil
}

func (c *Config) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return NotFound
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return FormatError
	}

	return nil
}

func SetContextConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		from := strings.ToLower(c.Query("from"))
		subfrom := strings.ToLower(c.Query("subfrom"))

		if from == "shiwan" {
			from = "web"
		}

		config := Get(from)

		// web内嵌页
		if subfrom == "web" {
			configBytes, _ := json.Marshal(config)
			subConfig := new(Config)
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
