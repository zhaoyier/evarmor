package utils

import (
	"fmt"
	"strconv"
	"strings"

	"git.ezbuy.me/ezbuy/evarmor/common/log"
)

//AToI32 字符串转整数
func AToI32(str string) int32 {
	str = strings.TrimSpace(str)
	str = strings.Replace(str, ",", "", -1)
	num, err := strconv.Atoi(str)
	if err != nil {
		ret, _ := strconv.ParseFloat(str, 64)
		return int32(ret)
	}
	return int32(num)
}

//AToI64 字符串转整数
func AToI64(str string) int64 {
	str = strings.TrimSpace(str)
	str = strings.Replace(str, ",", "", -1)
	num, _ := strconv.Atoi(str)
	return int64(num)
}

//GetTurnover 获取成交额
func GetTurnover(str string) int64 {
	str = strings.Replace(strings.TrimSpace(str), ",", "", -1)
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Errorf("get turnover failed: %s", err.Error())
		return 0
	}
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return int64(value * 10000)
}

func AtoF64(str string) string {
	if str == "" {
		return "-"
	}
	val := float64(AToI64(str)) / 100
	return fmt.Sprintf("%.2f", val)
}

func GetPercent(numerator, denominator int32) int32 {
	val := float64(numerator) / float64(denominator)
	return int32(val * 100)
}
