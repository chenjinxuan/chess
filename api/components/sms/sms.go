package sms

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
	"treasure/components/auth"
	"treasure/config"
	"treasure/databases"
	"treasure/define"
	"treasure/log"
	"treasure/models"
	"treasure/rpc"
)

const (
	// 短信类型
	SMS_MIN = iota
	SMS_LOGIN
	SMS_REGISTER
	SMS_PWD_RESET
	SMS_BIND           // 登录号码的绑定
	SMS_CONTACT_MOBIlE // 联系电话绑定
	SMS_CHANGE_MOBILE  // 更换绑定手机
	SMS_MERGE          // 合并账户时绑定手机
	SMS_JOIN_DRAW      // 参与抽奖
	SMS_MAX

	// SMS Error Code
	Wrong_Mobile_Number           = -20
	Send_Limit                    = -21
	Check_Limit                   = -22
	Send_Fail                     = -23 //短信发送失败,可能是超过了限制.
	Sms_Reset_Pwd_Mobile_No_Exist = -24 //重置密码的时候输入了个不存在的手机号
	Sms_Reg_Mobile_Exist          = -26 // 注册的时候手机号已经存在
	Sms_Contact_Mobile_Exist      = -27 // 绑定号码的时候手机号已经存在

	// SMS Check Code
	Code_Wrong  = -25
	Code_Expire = -1

	Need_Change_Mobile_Verify        = -61 //未经过更换绑定手机安全验证
	Change_Mobile_Code_Check_Expired = -62 //上次跟换绑定手机安全验证已过期
)

var (
	VerifyCodeExpired            = errors.New("Verify code is expired.")
	NeedChangeMobileVerify       = errors.New("Need change mobile verify")
	ChangeMobileCodeCheckExpired = errors.New("Change mobile code check expired.")
	UnexpectedSigningMethod      = errors.New("Unexpected signing method.")
	AuthFailed                   = errors.New("Auth failed.")
	MsgContentWrong              = errors.New("Message content wrong.")
	SMS                          = new(SMSHelper)
)

type SendReponse struct {
	Code    int             `json:"Code"`
	Message string          `json:"Message"`
	Data    SendReponseData `json:"Data"`
}

type SendReponseData struct {
	StateCode int
	Message   string
}

type SMSHelper struct {
}

func genCode() string {
	tn := time.Now().UnixNano()
	tns := strconv.Itoa(int(tn))
	return tns[13:]
}

// 发送验证码，返回code
func SendCode(MobileNumber, ip string, vtype int, cSms *config.Sms, from string) (int, error) {
	var verifyCode string
	var err error

	if !CheckMobileNumber(MobileNumber, cSms.MobileNumberRegular) {
		err = errors.New("mobile number wrong")
		return Wrong_Mobile_Number, err
	}

	err = SendMobileLimit(MobileNumber, cSms)
	if err != nil {
		return Send_Limit, err
	}

	verifyCode = genCode()

	// Insert to mobile verify database
	expire := time.Now().Add(30 * time.Minute).Unix()
	mvID, err := models.UsersMobileVerify.Insert(MobileNumber, verifyCode, ip, vtype, int(expire))
	if err != nil {
		// result.Msg = "Could not insert verify code."
		return 0, err
	}

	// Send login verify SMS - call RPC
	var rpcArgs rpc.SendSmsArgs
	rpcArgs.From = from
	rpcArgs.Phone = MobileNumber
	rpcArgs.TplArgs = append(rpcArgs.TplArgs, verifyCode)
	rpcArgs.TplKey = fmt.Sprintf("verify-%d", vtype)
	//rpcArgs.Sign.Source = cSms.Source
	//rpcArgs.Sign.Suffix = cSms.ContentPrefix

	rpcCli := rpc.Notifier.Get()
	var ret bool
	err = rpcCli.Call(rpc.NotifierServiceSendSms, rpcArgs, &ret)
	if err != nil {
		return Send_Fail, err
	}

	models.UsersMobileVerify.SetSendStatus(mvID, 1)
	return 1, nil
}

func SendToPhone(mobileNumber, ip, content, source string) error {
	apiUrlFmt := "http://sms.api.tongbu.com/message.ashx?SMSgxmtgame&t=%v&ip=%v&token=%v&source=%v&phone=%v&content=%v&stime="

	t := strconv.Itoa(int(time.Now().Unix()))
	tokenSource := mobileNumber + t + define.SmsSecret
	h := md5.New()
	io.WriteString(h, tokenSource)
	token := fmt.Sprintf("%x", h.Sum(nil))
	apiUrl := fmt.Sprintf(apiUrlFmt, t, ip, token, source, mobileNumber, url.QueryEscape(content))
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(apiUrl)

	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return err
	}
	var data SendReponse
	err = json.Unmarshal(body, &data)
	fmt.Println(fmt.Sprintf("%s-sms-send-Log:  Content: %s Reponse Body: %s ", mobileNumber, content, string(body)))
	// fmt.Printf(string(body))
	if err != nil {
		return err
	}
	if data.Code != 0 {
		return errors.New("send sms fail")
	}
	if data.Data.StateCode != 1 {
		return errors.New("send sms fail")
	}
	return nil
}

func SendPromotionToPhone(mobileNumber, ip, content, source string) (err error, rspStr string) {
	apiUrlFmt := `http://sms.api.tongbu.com/message.ashx?SMSmt&t=%v&token=%v&source=%v&ip=%v&phone=%v&content=%v&stime=`

	t := strconv.Itoa(int(time.Now().Unix()))
	tokenSource := mobileNumber + t + define.SmsSecret
	h := md5.New()
	io.WriteString(h, tokenSource)
	token := fmt.Sprintf("%x", h.Sum(nil))
	apiUrl := fmt.Sprintf(apiUrlFmt, t, token, source, ip, mobileNumber, url.QueryEscape(content))
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(apiUrl)

	if err != nil {
		return err, ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	rspStr = string(body)
	if err != nil {
		// handle error
		return err, rspStr
	}
	var data SendReponse
	err = json.Unmarshal(body, &data)
	log.Log.Info(apiUrl)
	defer log.Log.Info(fmt.Sprintf("%s-sms-send-Log:  Content: %s Reponse Body: %s ", mobileNumber, content, string(body)))
	// fmt.Printf(string(body))
	if err != nil {
		return err, rspStr
	}
	if data.Code != 0 {
		return errors.New("send sms fail"), rspStr
	}
	if data.Data.StateCode != 1 {
		return errors.New("send sms fail"), rspStr
	}
	return nil, rspStr
}

func CheckCode(mobile, code string, Type int, cSms *config.Sms) (int, string, error) {
	// check limit
	//err := CheckMobileLimit(mobile)
	count, err := GetCheckMobileCount(mobile)
	// @todo
	if err != nil {
		ret := Check_Limit
		msg := "system error"
		return ret, msg, err
	}
	if count > cSms.CheckLimit {
		ret := Check_Limit
		msg := "over check limit"
		return ret, msg, errors.New("over check limit")
	}
	// Get code
	var userMobileVerify = new(models.UsersMobileVerifyModel)
	userMobileVerify.MobileNumber = mobile
	userMobileVerify.VerifyType = Type
	if err := models.UsersMobileVerify.GetCode(userMobileVerify); err != nil {
		msg := auth.AuthFailed.Error()
		ret := Code_Wrong
		MobileCheckLimitCountPlus(mobile, cSms)
		return ret, msg, auth.AuthFailed
	}

	// Verify code is true
	if userMobileVerify.Status == 1 || userMobileVerify.VerifyCode != code {
		msg := AuthFailed.Error()
		ret := Code_Wrong
		MobileCheckLimitCountPlus(mobile, cSms)
		return ret, msg, AuthFailed
	}

	// Verify code is expired
	if int64(userMobileVerify.Expire) < time.Now().Unix() {
		ret := Code_Expire
		msg := VerifyCodeExpired.Error()
		return ret, msg, VerifyCodeExpired
	}

	// Set code used
	if err := models.UsersMobileVerify.SetCodeStatus(userMobileVerify.Id, 1); err != nil {
		ret := 0
		msg := "Could not update code status"
		return ret, msg, errors.New(msg)
	}

	return 1, "ok", nil
}

func CheckMobileNumber(mobile string, mobileNumberRegular string) bool {
	//regular := `^1([378][0-9]|4[57]|5[^4])\d{8}$`
	//regular := `^1(3[0-9]|4[57]|5[0-35-9]|8[0-9]|70)\d{8}$`
	regular := mobileNumberRegular
	reg := regexp.MustCompile(regular)
	res := reg.MatchString(mobile)
	return res
}

func getMobileCheckKey(mobile string) string {
	return "mobile-check-" + mobile
}

func GetCheckMobileCount(mobile string) (int, error) {
	key := getMobileCheckKey(mobile)
	isExist, err := databases.Redis.Sms.Exists(key)
	if err != nil {
		return 0, errors.New("system error")
	}
	// 检查key是否存在
	if !isExist {
		return 0, nil
	}

	// 获取当前count
	countStr, err := databases.Redis.Sms.Get(key)
	if err != nil {
		return 0, errors.New("system error")
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, err
	}
	return count, nil

}

func MobileCheckLimitCountPlus(mobile string, cSms *config.Sms) error {
	key := getMobileCheckKey(mobile)
	isExist, err := databases.Redis.Sms.Exists(key)
	if err != nil {
		return err
	}
	// 检查key是否存在
	if !isExist {
		err = databases.Redis.Sms.Setex(key, "1", cSms.Time)
		return err
	}

	// 获取当前count
	countStr, err := databases.Redis.Sms.Get(key)
	if err != nil {
		return err
	}

	count, err := strconv.Atoi(countStr)
	//
	c := count + 1
	err = databases.Redis.Sms.Setex(key, strconv.Itoa(c), cSms.Time)
	return nil
}

func SendMobileLimit(mobile string, cSms *config.Sms) error {
	key := "mobile-send-" + mobile
	// sms := databases.Redis.Sms
	isExist, err := databases.Redis.Sms.Exists(key)
	if err != nil {
		return errors.New("system error")
	}
	// 检查key是否存在
	if !isExist {
		err = databases.Redis.Sms.Setex(key, "1", cSms.Time)
		if err != nil {
			return err
		}
		return nil
	}

	// 获取当前count
	countStr, err := databases.Redis.Sms.Get(key)
	if err != nil {
		return errors.New("system error")
	}

	count, err := strconv.Atoi(countStr)
	//
	c := count + 1
	err = databases.Redis.Sms.Setex(key, strconv.Itoa(c), cSms.Time)
	if err != nil {
		return errors.New("system error,cant update count")
	}

	if count > cSms.SendLimit {
		return errors.New("over send code limte")
	}
	return nil
}

func TypeCheck(smsType int) bool {
	if smsType >= SMS_MAX || smsType <= SMS_MIN {
		return false
	}
	return true
}

func CheckCodeNotUpdateStatus(mobile, code string, Type int, cSms *config.Sms) (int, string, error) {
	// check limit
	count, err := GetCheckMobileCount(mobile)
	// @todo
	if err != nil {
		ret := Check_Limit
		msg := "system error"
		return ret, msg, err
	}
	if count > cSms.CheckLimit {
		ret := Check_Limit
		msg := "over check limit"
		return ret, msg, errors.New("over check limit")
	}
	// Get code
	var userMobileVerify = new(models.UsersMobileVerifyModel)
	userMobileVerify.MobileNumber = mobile
	userMobileVerify.VerifyType = Type
	if err := models.UsersMobileVerify.GetCode(userMobileVerify); err != nil {
		msg := auth.AuthFailed.Error()
		ret := Code_Wrong
		MobileCheckLimitCountPlus(mobile, cSms)
		return ret, msg, auth.AuthFailed
	}

	// Verify code is true
	if userMobileVerify.Status == 1 || userMobileVerify.VerifyCode != code {
		msg := AuthFailed.Error()
		ret := Code_Wrong
		MobileCheckLimitCountPlus(mobile, cSms)
		return ret, msg, AuthFailed
	}

	// Verify code is expired
	if int64(userMobileVerify.Expire) < time.Now().Unix() {
		ret := Code_Expire
		msg := VerifyCodeExpired.Error()
		return ret, msg, VerifyCodeExpired
	}
	return 1, "ok", nil
}

func CheckChangeMobileCode(mobile string) (int, string, error) {
	var userMobileVerify = new(models.UsersMobileVerifyModel)
	userMobileVerify.MobileNumber = mobile
	log.Log.Debug(userMobileVerify.MobileNumber)
	userMobileVerify.VerifyType = SMS_CHANGE_MOBILE
	log.Log.Debug(userMobileVerify.VerifyType)
	if err := models.UsersMobileVerify.GetCodeVerified(userMobileVerify); err != nil {
		ret := Need_Change_Mobile_Verify
		msg := NeedChangeMobileVerify.Error()
		return ret, msg, NeedChangeMobileVerify
	}
	// 上次验证5分钟后失效
	if int64(userMobileVerify.Expire) < time.Now().Unix() {
		ret := Change_Mobile_Code_Check_Expired
		msg := ChangeMobileCodeCheckExpired.Error()
		return ret, msg, ChangeMobileCodeCheckExpired
	}
	return 1, "ok", nil
}
