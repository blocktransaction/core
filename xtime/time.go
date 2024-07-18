package xtime

import (
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

// 获取毫秒时间戳
func MilliSecond() int64 {
	return time.Now().UnixNano() / 1e6
}

// 获取秒
func Second() int64 {
	return time.Now().Unix()
}

// 将时间戳转成yyyy-dd-mm hh:mm:ss格式
func TimestampToDatetime(timestamp string) string {
	d, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return ""
	}

	t1 := time.Unix(d, 0)
	// x<=-2 || x>=2
	if decimal.NewFromFloat(time.Since(t1).Abs().Minutes()).Cmp(decimal.NewFromInt(2)) >= 0 {
		d = Second()
	}
	return time.Unix(d, 0).Format(time.DateTime)
}

// 检查时间是否过期
func CheckTimeExpired(expireAt string) bool {
	if expireAt == "" {
		return false
	}

	t, _ := time.ParseInLocation(time.DateTime, expireAt, time.Local)
	//跟当前时期比较,如果大于等于代表已过期
	return decimal.NewFromFloat(time.Since(t).Seconds()).Cmp(decimal.NewFromInt(0)) >= 0
}

// 当前时间
func CurrentTimeStr() string {
	return time.Now().Format(time.DateTime)
}

// 截止到指定天，并转成秒的时间戳
func ToDay(day int) int64 {
	return time.Now().Add(time.Hour * time.Duration(24*day)).Unix()
}

// 至今天 23:59:59剩余秒数
func ToTodaySecond() int64 {
	now := time.Now()

	t := time.DateOnly + " 23:59:59"
	today, _ := time.Parse(time.DateTime, t)

	return today.Unix() - 28799 - now.Unix()
}
