package consul

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"chess/common/log"
)

var (
	ConsulCfg       ConsulConfig
	ConsulClient *ConsulCliWrap
)

var (
	ErrCliNil               = errors.New("consul api not init")
	ErrAddrOrDatacenterMiss = errors.New("consul need address and datacenter")
)

type ConsulCliWrap struct {
	cli        *api.Client
	datacenter string
	address    string
	token      string
}

type ConsulConfig struct {
	Address    string
	Datacenter string
	Token      string
	Proxy      string
}

func SetConsulCfg(addr, dc, token, proxy string) {
	ConsulCfg.Address = addr
	ConsulCfg.Datacenter = dc
	ConsulCfg.Token = token
	ConsulCfg.Proxy = proxy
}

func SetConsulCfgViaEnv() {
	ConsulCfg.Address = os.Getenv("CONSUL_ADDRESS")
	ConsulCfg.Datacenter = os.Getenv("CONSUL_DATACENTER")
	ConsulCfg.Token = os.Getenv("CONSUL_TOKEN")
	ConsulCfg.Proxy = os.Getenv("CONSUL_PROXY")
}

func InitConsulClient(addr, dc, token, proxy string) (err error) {
	SetConsulCfg(addr, dc, token, proxy)
	ConsulClient, err = NewConsulClient(ConsulCfg)
	return
}

func InitConsulClientViaEnv() (err error) {
	SetConsulCfgViaEnv()
	ConsulClient, err = NewConsulClient(ConsulCfg)
	return
}

func NewConsulClient(cfg ConsulConfig) (*ConsulCliWrap, error) {
	if cfg.Address == "" || cfg.Datacenter == "" {
		return nil, ErrAddrOrDatacenterMiss
	}

	var proxyClient *http.Client
	if cfg.Proxy != "" {
		proxyUrl, err := url.Parse(cfg.Proxy)
		if err != nil {
			proxyClient = &http.Client{}
		} else {
			proxyClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
		}
	} else {
		proxyClient = &http.Client{}
	}

	c := &api.Config{}
	c.Address = cfg.Address
	c.Datacenter = cfg.Datacenter
	c.Token = cfg.Token
	c.HttpClient = proxyClient
	cli, err := api.NewClient(c)

	if err != nil {
		return nil, err
	}

	wrap := ConsulCliWrap{cli: cli, address: cfg.Address, datacenter: cfg.Datacenter, token: cfg.Token}

	return &wrap, nil
}

//now only support get
func (c *ConsulCliWrap) Key(key string, def string) (string, error) {
	if c.cli != nil {
		kv := c.cli.KV()
		pair, _, err := kv.Get(key, nil)
		if pair == nil {
			return def, nil
		}
		log.Debugf("get key %s value %s", key, string(pair.Value))
		return string(pair.Value), err
	}

	return "", ErrCliNil
}

func (c *ConsulCliWrap) KeyInt(key string, def int) (int, error) {
	res, err := c.Key(key, "")
	if err != nil {
		return 0, err
	}

	if res == "" {
		return def, nil
	}
	resInt, err := strconv.Atoi(res)
	return resInt, err
}

func (c *ConsulCliWrap) KeyInt64(key string, def int64) (int64, error) {
	res, err := c.Key(key, "")
	if err != nil {
		return 0, err
	}

	if res == "" {
		return def, nil
	}

	resInt64, err := strconv.ParseInt(res, 10, 64)
	return resInt64, err
}

func (c *ConsulCliWrap) KeyBool(key string, def bool) (bool, error) {
	res, err := c.Key(key, "")
	if err != nil {
		return false, err
	}

	if res == "" {
		return def, nil
	}

	resBool, err := strconv.ParseBool(res)
	return resBool, err
}

func (c *ConsulCliWrap) KeySet(key string, def []string) ([]string, error) {
	res, err := c.Key(key, "")
	if err != nil {
		return []string{}, err
	}

	if res == "" {
		return def, nil
	}

	resSet := strings.Split(res, ",")
	return resSet, nil
}

func (c *ConsulCliWrap) KeyWatch(key string, bind *string) error {
	log.Debugf("Watching string key:%s", key)

	plan := c.WatchKeyPlanFactory(key)

	plan.Handler = c.WatchFuncFactory(key, func(raw interface{}) {
		if raw == nil {
			return
		}
		v, ok := raw.(*api.KVPair)
		if !ok || v == nil {
			log.Errorf("consul return invalid")
		} else {
			*(bind) = string(v.Value)
		}
	})

	go plan.Run(c.address)
	return nil
}

func (c *ConsulCliWrap) KeyIntWatch(key string, bind *int) error {
	log.Debugf("Watching int key:%s", key)

	plan := c.WatchKeyPlanFactory(key)

	plan.Handler = c.WatchFuncFactory(key, func(raw interface{}) {
		if raw == nil {
			return
		}
		v, ok := raw.(*api.KVPair)
		if !ok || v == nil {
			log.Errorf("consul return invalid")
		} else {
			valStr := string(v.Value)
			valInt, err := strconv.Atoi(valStr)
			if err != nil {
				log.Errorf("consul watch key:%s type error val:%s", key, valStr)
				return
			}
			*(bind) = valInt
		}
	})

	go plan.Run(c.address)
	return nil
}

func (c *ConsulCliWrap) KeyInt64Watch(key string, bind *int64) error {
	log.Debugf("Watching int64 key:%s", key)

	plan := c.WatchKeyPlanFactory(key)

	plan.Handler = c.WatchFuncFactory(key, func(raw interface{}) {
		if raw == nil {
			return
		}
		v, ok := raw.(*api.KVPair)
		if !ok || v == nil {
			log.Errorf("consul return invalid")
		} else {
			valStr := string(v.Value)
			valInt64, err := strconv.ParseInt(valStr, 10, 64)
			if err != nil {
				log.Errorf("consul watch key:%s type error val:%s", key, valStr)
				return
			}

			*(bind) = valInt64
		}
	})
	go plan.Run(c.address)
	return nil
}

func (c *ConsulCliWrap) KeyBoolWatch(key string, bind *bool) error {
	log.Debugf("Watching bool key:%s", key)

	plan := c.WatchKeyPlanFactory(key)

	plan.Handler = c.WatchFuncFactory(key, func(raw interface{}) {
		if raw == nil {
			return
		}
		v, ok := raw.(*api.KVPair)
		if !ok || v == nil {
			log.Errorf("consul return invalid")
		} else {
			valStr := string(v.Value)
			valBool, err := strconv.ParseBool(valStr)
			if err != nil {
				log.Errorf("consul watch key:%s type error val:%s", key, valStr)
				return
			}
			*(bind) = valBool
		}
	})

	go plan.Run(c.address)
	return nil
}

func (c *ConsulCliWrap) KeySetWatch(key string, bind *[]string) error {
	log.Debugf("Watching set key:%s", key)

	plan := c.WatchKeyPlanFactory(key)

	plan.Handler = c.WatchFuncFactory(key, func(raw interface{}) {
		if raw == nil {
			return
		}
		v, ok := raw.(*api.KVPair)
		if !ok || v == nil {
			log.Errorf("consul return invalid")
		} else {
			valStr := string(v.Value)
			resSet := strings.Split(valStr, ",")
			*(bind) = resSet
		}
	})

	go plan.Run(c.address)
	return nil
}

func (c *ConsulCliWrap) WatchFuncFactory(key string, handler func(raw interface{})) func(idx uint64, raw interface{}) {
	return func(idx uint64, raw interface{}) {
		log.Debugf("key %s changed:%d", key, idx)
		handler(raw)
	}
}

func (c *ConsulCliWrap) WatchKeyPlanFactory(key string) *watch.Plan {
	plan, _ := watch.Parse(map[string]interface{}{
		"datacenter": c.datacenter,
		"type":       "key",
		"token":      c.token,
		"key":        key,
	})
	return plan
}

func (c *ConsulCliWrap) RegisterService(registration *api.AgentServiceRegistration) error {
	return c.cli.Agent().ServiceRegister(registration)
}

func (c *ConsulCliWrap) Services() (map[string]*api.AgentService, error) {
	return c.cli.Agent().Services()
}

//func (c *ConsulCliWrap) ServiceWatch() error {
//
//}