package define



const (
    AuthFailedStatus = -998  //token错误
    AuthALreadyLogin = -997  //token已刷新
    AuthExpire       = -996  //token过期
)

var AuthMsgMap =make(map[int]string)

func init()  {
    AuthMsgMap[AuthFailedStatus]="token错误"
    AuthMsgMap[AuthALreadyLogin]="token已被刷新,请重新登录"
    AuthMsgMap[AuthExpire]="token已过期"
}