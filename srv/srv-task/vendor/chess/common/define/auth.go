package define

const (
	AuthFailedStatus = -998 //token错误
	AuthALreadyLogin = -997 //token已刷新
	AuthExpire       = -996 //token过期
	AuthKickedOut    = -995 //踢出
        AuthReToken      = -994 //应刷新
)

var AuthMsgMap = make(map[int]string)

func init() {
	AuthMsgMap[AuthFailedStatus] = "token错误"
	AuthMsgMap[AuthALreadyLogin] = "token已被刷新,请重新登录"
	AuthMsgMap[AuthExpire] = "token已过期"
	AuthMsgMap[AuthKickedOut] = "另一台设备登录,您已被踢出"
        AuthMsgMap[AuthReToken] = "请刷新token"
}
