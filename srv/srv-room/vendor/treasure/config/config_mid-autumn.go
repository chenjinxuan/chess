package config

type MidAutumn struct {
	ChargePool MidAutumnPoolConfig `json:"charge_pool"`
	FreePool   MidAutumnPoolConfig `json:"free_pool"`
}

type MidAutumnPoolConfig struct {
	PoolKey  string             `json:"pool_key"`
	PoolSize int                `json:"pool_size"` //池子大小 必须为100的整数倍
	PoolConf map[string]float32 `json:"pool_conf"` // 元素=>比率
}
