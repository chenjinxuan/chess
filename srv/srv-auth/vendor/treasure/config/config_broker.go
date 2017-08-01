package config

type Broker struct {
	LeaderboardLimit    int     `json:"leaderboard_limit"`
	BindSalt            string  `json:"bind_salt"`
	BindTimeLimit       float64 `json:"bind_time_limit"`        // 小时
	WithdrawMinPriceAli int     `json:"withdraw_min_price_ali"` // 支付宝最小提现金额
	WithdrawMinPriceDb  int     `json:"withdraw_min_price_db"`  // 夺宝币最小提现金额
}
