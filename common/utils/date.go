package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//GetCurrentDate 获取当前日期20191001
func GetCurrentDate(timestamp int64) (string, string) {
	now := time.Now()
	if timestamp != 0 {
		now = time.Unix(timestamp, 0)
	}
	year := strconv.Itoa(now.Year())
	date := strings.Split(now.Format("2006-01-02 15:04:05"), " ")
	return year, strings.Join(strings.Split(date[0], "-"), "")
}

//GetCurrentDate2 2020-01-02
func GetCurrentDate2(timestamp int64) string {
	now := time.Now()
	if timestamp != 0 {
		now = time.Unix(timestamp, 0)
	}
	date := strings.Split(now.Format("2006-01-02 15:04:05"), " ")
	return date[0]
}

func GetCurrentDate4(ts int64) string {
	now := time.Now()
	if ts != 0 {
		now = time.Unix(ts, 0)
	}

	return now.Format("2006-01-02")
}

//GetCurrentDate3 20200102
func GetCurrentDate3(timestamp int64) int32 {
	now := time.Now()
	if timestamp != 0 {
		now = time.Unix(timestamp, 0)
	}
	date := strings.Split(now.Format("2006-01-02 15:04:05"), " ")
	dd := strings.Join(strings.Split(date[0], "-"), "")
	return AToI32(dd)
}

func GetDateTime(timestamp int64) int32 {
	now := time.Now()
	if timestamp != 0 {
		now = time.Unix(timestamp, 0)
	}
	date := strings.Split(now.Format("2006-01-02 15:04:05"), " ")
	datetime := strings.Join(strings.Split(date[0], "-"), "")
	return AToI32(datetime)
}

func GetDateTime2(ts int64) string {
	now := time.Unix(ts, 0)
	return now.Format("2006-01-02 15:04:05")
}

//TransDateFormat 2019-10-01->20191001
func TransDateFormat(date string) string {
	return strings.Join(strings.Split(date, "-"), "")
}

//DateToTimestamp 日期转时间戳
func DateToTimestamp(date string) int64 {
	tm, _ := time.Parse("2006-01-02", date)
	return tm.Unix()
}

//IsValidContract 合约有效性
func IsValidContract(date int32) bool {
	now := time.Now()
	tm, _ := time.Parse("20061", fmt.Sprintf("%d", date))
	if tm.Year() > now.Year() {
		return true
	}
	if tm.Year() < now.Year() {
		return false
	}
	if tm.Month()-now.Month() <= 0 {
		return false
	}
	return true
}

// GetTimestamp 时间戳
func GetTimestamp() int64 {
	return time.Now().Unix()
}

func GetPrefixYear() string {
	year := fmt.Sprintf("%d", time.Now().Year())
	return year[0:2]
}

//GetDateDiff  20200102->20191229
func GetDateDiff(current, base int32) int32 {
	ctm, _ := time.Parse("20060102", fmt.Sprintf("%d", current))
	btm, _ := time.Parse("20060102", fmt.Sprintf("%d", base))
	subM := ctm.Sub(btm)
	return int32(subM.Hours() / 24)
}

func GetMonthDay(date int32) int32 {
	tm, _ := time.Parse("20060102", fmt.Sprintf("%d", date))
	return AToI32(fmt.Sprintf("%d%0d%0d", tm.Year()-2000, tm.Month(), tm.Day()))
}

func GetWeekDay(date int32) int32 {
	tm, _ := time.Parse("20060102", fmt.Sprintf("%d", date))
	return int32(tm.Weekday())
}

func IsWeekend(ts int64) bool {
	day := time.Now().Weekday()
	if day == time.Sunday || day == time.Saturday {
		return true
	}
	return false
}

func IsMonday(ts int64) bool {
	day := time.Now().Weekday()
	if day == time.Monday {
		return true
	}
	return false
}
