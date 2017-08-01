package define

//for lotterygen notify pool full
type NotifyPoolFull struct {
	Key      string `json:"key"`
	Expect   string `json:"expect"`
	OpenCode string `json:"opencode"`
	EndPoint int64  `json:"endpoint"`
}

//for lotterygen notify pool progress
type NotifyPoolProgress struct {
	Key      string `json:"key"`
	Count    int    `json:"count"`
	Left     int    `json:"left"`
	Size     int    `json:"size"`
	Progress int    `json:"progress"`
}

//for lotterycentre notify new betting gen
type NotifyNewBetting struct {
	GoodsId int
	Period  int
}

const (
	CouponCaseMaskOriginFlag = 1 << iota
	CouponCaseMaskSeconaryFlag
)

const (
	GenCouponOrderCodeNormal = 0
	GenCouponOrderCodeIgnore = 1
)

type CouponCase struct {
	Trigger string                 `json:"trigger"`
	Mask    int                    `json:"mask"`
	KVMap   map[string]interface{} `json:"kvmap"`
}

type NotifyGenCoupon struct {
	UserId    int        `json:"user_id"`
	OrderCode int        `json:"order_code"`
	Created   int64      `json:"created"`
	Case      CouponCase `json:"case"`
}
