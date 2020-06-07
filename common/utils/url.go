package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http get failed: %s", resp.Status)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

//GetZhengZhouDataDailyURL http://www.czce.com.cn/cn/DFSStaticFiles/Future/2019/20190930/FutureDataDaily.txt
func GetZhengZhouDataDailyURL(timestamp int64) string {
	year, date := GetCurrentDate(timestamp)
	return fmt.Sprintf("%s/%s/%s/%s", ZhengZhouDomain, year, date, ZhengZhouFutureDataDaily)
}

func GetDaLianDailyURL(timestamp int64) string {
	return DaLianDataDailyURL
}

//ZhengZhouHoldingURL //http://www.czce.com.cn/cn/DFSStaticFiles/Future/2019/20190930/FutureDataHolding.txt
func ZhengZhouHoldingURL(timestamp int64) string {
	year, date := GetCurrentDate(timestamp)
	return fmt.Sprintf("%s/%s/%s/%s", ZhengZhouDomain, year, date, ZhengZhouFutureDataHolding)
}

//ShanghaiDailyURL 上海每日
func ShanghaiDailyURL(timestamp int64) string {
	//http://www.shfe.com.cn/data/dailydata/kx/kx20191017.dat
	_, date := GetCurrentDate(timestamp)
	return fmt.Sprintf("http://www.shfe.com.cn/data/dailydata/kx/kx%s.dat", date)
}

func GetMainContractURL(timestamp int64) string {
	date := GetCurrentDate4(timestamp)
	return fmt.Sprintf("http://m.data.eastmoney.com/api/futures/GetContract?market=069001008&date=%s", date)
	// return fmt.Sprintf("http://m.data.eastmoney.com/api/futures/GetContract?market=069001008&date=2020-05-27")
}

func BillboardURL(exchange, contract string, timestamp int64) string {
	//http://m.data.eastmoney.com/futures/index
	//http://m.data.eastmoney.com/api/futures/GetQhcjcc?market=069001007&date=2020-01-14&contract=P2005&name=%E5%A4%9A%E5%A4%B4%E6%8C%81%E4%BB%93%E9%BE%99%E8%99%8E%E6%A6%9C&page=1

	market := "069001006"
	switch exchange {
	case "shanghai":
		market = "069001005"
	case "zhengzhou":
		market = "069001008"
	case "dalian":
		market = "069001007"
	}

	date := GetCurrentDate2(timestamp)
	return fmt.Sprintf("http://m.data.eastmoney.com/api/futures/GetQhcjcc?market=%s&date=%s&contract=%s&name=%s&page=1", market, date, contract, "%E5%A4%9A%E5%A4%B4%E6%8C%81%E4%BB%93%E9%BE%99%E8%99%8E%E6%A6%9C")
}
