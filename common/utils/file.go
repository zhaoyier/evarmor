package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func ReadFile(path string) ([]string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

func GetZhengZhouHoldingPath(timestamp int64) string {
	_, date := GetCurrentDate(timestamp)
	return fmt.Sprintf("%s/%s_%s", DownloadPath, date, ZhengZhouHoldingPath)
}

// GetZhengZhouFutureDataDailyPath 路径
func GetZhengZhouFutureDataDailyPath(timestamp int64) string {
	_, date := GetCurrentDate(timestamp)
	return fmt.Sprintf("%s/%s_%s", DownloadPath, date, ZhengZhouDailyPath)
}

func GetDaLianDailyPath(timestamp int64) string {
	_, date := GetCurrentDate(timestamp)
	return fmt.Sprintf("%s/%s_%s", DownloadPath, date, DaLianDataDailyPath)
}

//ShanghaiDailyPath 上海
func ShanghaiDailyPath(timestamp int64) string {
	_, date := GetCurrentDate(timestamp)
	return fmt.Sprintf("%s/%s_%s", DownloadPath, date, ShangHaiDailyPath)
}

func BillboardDailyPath(contract string, ts int64) string {
	_, date := GetCurrentDate(ts)
	return fmt.Sprintf("%s/%s_%s%s", DownloadPath, date, contract, BillboardPath)
}

//判断目录是否存在,不存在则创建
func GenerateDir() error {
	_, err := os.Stat(DownloadPath)
	if os.IsNotExist(err) {
		return os.Mkdir(DownloadPath, os.ModePerm)
	}
	return nil
}

func DownloadBaseCheck() error {
	//判断时间
	day := time.Now().Weekday()
	if day == time.Sunday || day == time.Saturday {
		return errors.New("today is weekend")
	}
	//判断目录是否存在
	if _, err := os.Stat(DownloadPath); os.IsNotExist(err) {
		os.Mkdir(DownloadPath, os.ModePerm)
	}
	return nil
}
