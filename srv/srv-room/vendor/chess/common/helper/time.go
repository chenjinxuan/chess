package helper

import (
	"strconv"
	"strings"
	"time"
)

const (
	FormatDate             = "2006-01-02"
	FormatDateHi           = "2006-01-02 15:04"
	FormatDateTime         = "2006-01-02 15:04:05"
	FormatDateTimeMilliSec = "2006-01-02 15:04:05.000"
)

// 获取毫秒级时间戳
func GetNowUnixMilli() int64 {
	return int64(time.Now().UnixNano() / 1000000)
}

// 获取秒级时间戳
func GetNowUnix() int64 {
	return time.Now().Unix()
}

// 纳秒级时间戳转毫秒级时间戳 19位->13位
func UnixNano2UnixMilli(nano int64) int64 {
	return int64(nano / 1000000)
}

// 纳秒级时间戳转秒级时间戳 19位->10位
func UnixNano2Unix(nano int64) int64 {
	return int64(nano / 1000000000)
}

// 毫秒级时间戳转秒级时间戳 13位->10位
func UnixMilli2Unix(milli int64) int64 {
	return int64(milli / 1000)
}

// 格式化毫秒级时间戳 13位->2006-01-02 15:04:05.000
func UnixMilli2TimeStr(milli int64) string {
	return time.Unix(0, int64(milli*1000000)).Format(FormatDateTimeMilliSec)
}

// 格式化毫秒级时间戳 13位->2006-01-02
func UnixMilli2DateStr(milli int64) string {
	return time.Unix(0, int64(milli*1000000)).Format("2006-01-02")
}

// 毫秒级时间戳转时间  13位->time.Time
func UnixMilli2Time(milli int64) time.Time {
	return time.Unix(int64(milli/1000), 0)
}

//2014-12-25 18:12:25  --->   1420511108210
func TimeStr2Unix(timeStr string) (int64, error) {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

//2014-12-25 18:12:25  --->   1420511108210
func TimeStr2UnixMilli(timeStr string) (int64, error) {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
	if err != nil {
		return 0, err
	}
	return int64(t.UnixNano() / 1000000), nil
}

//Local 2014-12-25 18:12:25  --->   time.Time
func TimeStr2Time(timeStr string) (time.Time, error) {
	local := time.Local
	return time.ParseInLocation("2006-01-02 15:04:05", timeStr, local)
}

func GetDate() string {
	return time.Now().Format(FormatDate)
}

func GetTime() string {
	return time.Now().Format(FormatDateTime)
}

// 2006-01-02 15:04:05.123 =>  150405123
func Time2Int64(t int64) int64 {
	datestring := UnixMilli2TimeStr(t)
	_time := datestring[11:]
	timestr := _time + strings.Repeat("0", 12-len(_time))
	timeintstr := strings.Replace(strings.Replace(timestr, ":", "", -1), ".", "", -1)
	timeToInt, _ := strconv.ParseInt(timeintstr, 10, 64)
	return timeToInt
}
