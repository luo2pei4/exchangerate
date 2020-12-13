package crawler

import (
	db "exchangerate/db"
	dto "exchangerate/dto"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// CurrencyInfoMap 货币基础数据信息
var CurrencyInfoMap map[string]int

// BankTimeZoneMap 银行所在时区Map
var BankTimeZoneMap map[int]string

// Crawler 爬虫
type Crawler struct {
	conn *db.Connection // 数据库连接
}

// NewInstance 获取爬虫实例
func NewInstance() (crawler *Crawler, err error) {

	conn, err := db.NewConnection()

	if err != nil {
		return nil, err
	}

	return &Crawler{
		conn: conn,
	}, nil
}

// InitializeInfo 初始化相关数据
func InitializeInfo() error {

	conn, err := db.NewConnection()

	if err != nil {
		return err
	}

	CurrencyInfoMap = db.GetCurrencyMap(conn)
	BankTimeZoneMap = db.GetBankTimeZoneMap(conn)

	// conn.Close()

	return nil
}

// ConvertReleaseTime 转换发布时间
func (crawler *Crawler) ConvertReleaseTime(value string, bankID int) time.Time {

	if value != "" {

		value = strings.ReplaceAll(value, ".", "-")
		value = strings.ReplaceAll(value, "年", "-")
		value = strings.ReplaceAll(value, "月", "-")
		value = strings.ReplaceAll(value, "日", "")
		value = strings.ReplaceAll(value, " ", "T")
		value = value + BankTimeZoneMap[bankID]
		releaseTime, err := time.Parse(time.RFC3339, value)

		if err == nil {
			return releaseTime
		}
	}

	return time.Now()
}

// GetDateByTimezone 获取所在时区的日期
func (crawler *Crawler) GetDateByTimezone(timezone string) string {

	if timezone == "" {
		return ""
	}

	loc, err := time.LoadLocation(timezone)

	if err != nil {
		fmt.Println("Get date failed.")
		return ""
	}

	now := time.Now().In(loc)
	timestamp := now.Format(time.RFC3339)
	date := timestamp[:10]

	return date
}

// GetPageDoc 获取网页的Document
func (crawler *Crawler) GetPageDoc(pageURL string) *goquery.Document {

	response, err := http.Get(pageURL)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	statusCode := response.StatusCode

	if statusCode != 200 {
		fmt.Println(response.Status)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return doc
}

// SaveEachBankExchangeRate 保存各个银行发布的汇率数据
func (crawler *Crawler) SaveEachBankExchangeRate(rateInfo *dto.ExchangeRateInfo, table string) {

	if rateInfo == nil {
		fmt.Println("rate info is nil.")
		return
	}

	if rateInfo.CurrencyID != 0 {

		_, _, err := crawler.conn.InsertExchangeRateData(table, rateInfo)

		if err != nil {
			fmt.Println(err.Error())
		}
	}

}

// InitExchangeRateInfoMap 初始化各个银行最后一次的汇率数据Map
func (crawler *Crawler) InitExchangeRateInfoMap(table string) (exchangeRateInfoMap map[int]*dto.ExchangeRateInfo, err error) {
	return crawler.conn.SelectExchangeRateInfoMap(table)
}
