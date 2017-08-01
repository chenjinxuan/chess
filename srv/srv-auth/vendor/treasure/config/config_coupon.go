package config

type CouponEventList struct {
	NewUser     CouponEvent      `json:"new_user"`
	GroupCoupon GroupCouponEvent `json:"group_coupon"`
}

type CouponEvent struct {
	Tag      string `json:"tag"`
	ShareGag int64  `json:"share_gap"`
}

type GroupCouponEvent struct {
	Tag       string  `json:"tag"`
	DailyLimt int64   `json:"daily_limit"`
	SentTimes []int64 `json:"sent_times"`
}
