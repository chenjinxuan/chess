package user_init

//import (
//	"database/sql"
//	"encoding/json"
//	"errors"
//	"fmt"
//	"strings"
//	"chess/common/config"
//	//"chess/common/db"
//	"chess/api/helper"
//	"chess/api/log"
//	"chess/models"
//)
//
//var ParamsInvalid = errors.New("params invalid")
//var FoundNoHandler = errors.New("found no handler")
//var UserHasLoginBefore = errors.New("user has login before")
//var DeviceNoWhite = errors.New("device no in white list")
//
//type UserInitSnapshot map[string]interface{}
//
//type UserInitInterface interface {
//	GetUserInitSnapshot(users models.UsersModel, keyCase map[string]interface{}, userInit *config.UserInit) (UserInitSnapshot, error)
//}
//
//var UserInitMap map[string]UserInitInterface
//
//func init() {
//	UserInitMap = make(map[string]UserInitInterface)
//
//	UserInitMap["android"] = new(AndroidUserInit)
//	UserInitMap["ios"] = new(IOSUserInit)
//	UserInitMap["shiwan"] = new(IOSUserInit)
//	UserInitMap["web"] = new(WebUserInit)
//}
//
//func Filter(extra map[string]interface{}, white map[string][]string) bool {
//	from, ok := extra["from"].(string)
//	if !ok {
//		return false
//	}
//
//	from = strings.ToLower(from)
//	var uni string
//
//	if strings.Contains(from, "android") {
//		uniqueId, ok := extra["unique_id"].(string)
//		if !ok || uniqueId == "" {
//			return false
//		}
//		uni = uniqueId
//	} else {
//		idfa, ok := extra["idfa"].(string)
//		if !ok || idfa == "" {
//			return false
//		}
//		uni = idfa
//	}
//
//	// 白名单
//	if helper.StringInSlice(white[from], uni) != -1 {
//		return true
//	}
//
//	return false
//}
//
//func GetUniqueVal(from string, extra map[string]interface{}) string {
//	var uni string
//
//	if strings.Contains(from, "android") {
//		uniqueId, ok := extra["unique_id"].(string)
//		if !ok || uniqueId == "" {
//			return ""
//		}
//		uni = uniqueId
//	} else {
//		idfa, ok := extra["idfa"].(string)
//		if !ok || idfa == "" {
//			return ""
//		}
//		uni = idfa
//	}
//
//	return uni
//}
//
//func UserInit(user models.UsersModel, extra map[string]interface{}, cConf *config.Config) (int, error) {
//	// 判断是否首次登陆
//	r_key := fmt.Sprintf("userinit_has_login_%d", user.Id)
//	hasLogin, _ := models.ChessRedis.Login.GetInt(r_key)
//        models.ChessRedis.Login.Set(r_key, "1")
//	if hasLogin == 1 {
//		return 0, UserHasLoginBefore
//	}
//
//	from, ok := extra["from"].(string)
//	if !ok {
//		return 0, ParamsInvalid
//	}
//
//	// todo filter
//	if cConf.UserInit.Filter.Enable && !Filter(extra, cConf.UserInit.Filter.White) {
//		return 0, DeviceNoWhite
//	}
//
//	from = strings.ToLower(from)
//	var initHandler UserInitInterface
//	if strings.Contains(from, "android") {
//		initHandler, _ = UserInitMap["android"]
//	} else {
//		initHandler, _ = UserInitMap["ios"]
//	}
//
//	res, err := initHandler.GetUserInitSnapshot(user, extra, cConf.UserInit)
//	log.Log.Debugf("init user snapshot %+v", res)
//	is_fresh, ok := res["is_fresh"].(int)
//	if !ok {
//		is_fresh = 0
//	}
//
//	// 插入红包生成队列，由后端服务判断是否生成红包
//	//kvmap := extra
//	//kvmap["user_id"] = user.Id
//	//kvmap["unique_val"] = GetUniqueVal(from, extra)
//	//kvmap["is_fresh"] = is_fresh
//	//go coupon.PushToCouponGenQueue("userinit", user.Id, kvmap, cConf)
//
//	return is_fresh, err
//}
//
//func getInitUserTemplate(isFresh bool, userInit *config.UserInit) (res UserInitSnapshot) {
//	res = make(UserInitSnapshot)
//
//	if userInit.Enable {
//		if isFresh {
//			res["init_balance"] = userInit.Balance
//			res["is_fresh"] = 1
//		} else {
//			res["init_balance"] = 0
//			res["is_fresh"] = 0
//		}
//	} else {
//		log.Log.Debug("not enable init user")
//	}
//	return
//}
//
//func saveRegisterLog(isfresh bool, userid int, initstatus int, snapshot UserInitSnapshot, channel, from, ver string, key, val string, userInit *config.UserInit) (err error) {
//	res := getInitUserTemplate(isfresh, userInit)
//
//	log.Log.Debugf("%+v", res)
//
//	//save log
//	registerLog := new(models.UsersRegisterLogModel)
//	registerLog.UserId = userid
//	registerLog.InitStatus = initstatus
//	_snapshot, _ := json.Marshal(res)
//	registerLog.InitSnapshot = string(_snapshot)
//	registerLog.Channel = channel
//	registerLog.From = from
//	registerLog.Ver = ver
//	registerLog.DeviceUniqueKey = key
//	registerLog.DeviceUniqueVal = val
//	err = registerLog.Insert()
//	if err != nil {
//		return err
//	}
//
//	//init balance
//	//balance, ok := res["init_balance"].(int)
//	//if ok && balance != 0 {
//	//	log.Log.Debugf("add cash (%d,%d)", userid, balance)
//	//	err = models.UsersAddCashLog.AddCash(userid, balance, 1, "userinit", "forfreshuser")
//	//	if err != nil {
//	//		log.Log.Error("add cash fail")
//	//		return
//	//	}
//	//}
//
//	//update set user is fresh
//	if isfresh {
//		log.Log.Debug("update fresh:", userid)
//		err = models.Users.UpdateFresh(userid, 1)
//		if err != nil {
//			log.Log.Error(err)
//			return
//		}
//	}
//	return
//}
//
////Android 用户初始化实现
//type AndroidUserInit struct{}
//
//func (i *AndroidUserInit) GetUserInitSnapshot(users models.UsersModel, keyCase map[string]interface{}, userInit *config.UserInit) (res UserInitSnapshot, err error) {
//	res = make(UserInitSnapshot)
//	//check unique from device
//	uniqueId, ok := keyCase["unique_id"].(string)
//	if !ok || uniqueId == "" {
//		err = ParamsInvalid
//		return
//	}
//
//	_, errDevice := models.Device.GetUserIdByUniqueId(uniqueId)
//	if errDevice != nil && errDevice != sql.ErrNoRows {
//		err = errDevice
//		return
//	}
//
//	//check unique from users_register_log
//	logs, _ := models.UsersRegisterLog.GetLog("unique_id", uniqueId, users.Id)
//
//	ver, ok := keyCase["ver"].(string)
//	if !ok {
//		ver = ""
//	}
//
//	if errDevice == sql.ErrNoRows && len(logs) == 0 {
//		res["is_fresh"] = 1
//		err = saveRegisterLog(true, users.Id, models.UserInitStatusFresh, res, users.Channel, users.AppFrom, ver, "unique_id", uniqueId, userInit)
//	} else {
//		res["is_fresh"] = 0
//		err = saveRegisterLog(false, users.Id, models.UserInitStatusFresh, res, users.Channel, users.AppFrom, ver, "unique_id", uniqueId, userInit)
//	}
//	return
//}
//
////IOS 用户初始化实现
//type IOSUserInit struct{}
//
//func (i *IOSUserInit) GetUserInitSnapshot(users models.UsersModel, keyCase map[string]interface{}, userInit *config.UserInit) (res UserInitSnapshot, err error) {
//	res = make(UserInitSnapshot)
//	//check unique from device
//	idfa, ok := keyCase["idfa"].(string)
//	if !ok || idfa == "" {
//		err = ParamsInvalid
//		return
//	}
//
//	_, errDevice := models.Device.GetUserIdByIdfa(idfa)
//	if errDevice != nil && errDevice != sql.ErrNoRows {
//		err = errDevice
//		return
//	}
//
//	//check unique from users_register_log
//	logs, _ := models.UsersRegisterLog.GetLog("idfa", idfa, users.Id)
//
//	ver, ok := keyCase["ver"].(string)
//	if !ok {
//		ver = ""
//	}
//
//	if errDevice == sql.ErrNoRows && len(logs) == 0 {
//		res["is_fresh"] = 1
//		err = saveRegisterLog(true, users.Id, models.UserInitStatusFresh, res, users.Channel, users.AppFrom, ver, "idfa", idfa, userInit)
//	} else {
//		res["is_fresh"] = 0
//		err = saveRegisterLog(false, users.Id, models.UserInitStatusFresh, res, users.Channel, users.AppFrom, ver, "idfa", idfa, userInit)
//	}
//	return
//}
//
////试玩 用户初始化实现
//type ShiwanUserInit struct{}
//
//func (i *ShiwanUserInit) GetUserInitSnapshot(users models.UsersModel, keyCase map[string]interface{}, userInit *config.UserInit) (res UserInitSnapshot, err error) {
//	return
//}
//
////Web 用户初始化实现
//type WebUserInit struct{}
//
//func (i *WebUserInit) GetUserInitSnapshot(users models.UsersModel, keyCase map[string]interface{}, userInit *config.UserInit) (res UserInitSnapshot, err error) {
//	return
//}
