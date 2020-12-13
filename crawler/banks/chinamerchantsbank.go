package banks

import (
	"exchangerate/crawler"
	dto "exchangerate/dto"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var chinaMerchantsBankCrawler *crawler.Crawler
var chinaMerchantsBankExchangeRateMap map[int]*dto.ExchangeRateInfo

// ChinaMerchantsBank 爬取数据
func ChinaMerchantsBank() {

	if chinaMerchantsBankCrawler == nil {

		c, err := crawler.NewInstance()

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		chinaMerchantsBankCrawler = c
	}

	doc := chinaMerchantsBankCrawler.GetPageDoc("http://fx.cmbchina.com/hq/")

	if doc == nil {
		return
	}

	if chinaMerchantsBankExchangeRateMap == nil {

		erMap, err := chinaMerchantsBankCrawler.InitExchangeRateInfoMap("b_china_cmb")

		if err != nil {
			fmt.Println("Init exchange rate info map failed.")
			return
		}

		chinaMerchantsBankExchangeRateMap = erMap
	}

	// 因为招商银行的发布时间没有带日期，所以这个地方需要按时区来获取日期
	date := chinaMerchantsBankCrawler.GetDateByTimezone("Asia/Shanghai")

	doc.Find("#realRateInfo > table > tbody").Each(func(i int, table *goquery.Selection) {

		table.Find("tr").Each(func(j int, tr *goquery.Selection) {

			rateInfo := &dto.ExchangeRateInfo{}

			tr.Find("td").Each(func(k int, td *goquery.Selection) {

				value := td.Text()
				value = strings.ReplaceAll(value, "\n", "")
				value = strings.Trim(value, " ")

				switch k {
				case 0:
					rateInfo.CurrencyName = value
					currencyID := crawler.CurrencyInfoMap[rateInfo.CurrencyName]
					rateInfo.CurrencyID = currencyID
				case 3:
					rateInfo.SellingRate, _ = strconv.ParseFloat(value, 64)
				case 4:
					rateInfo.CashSellingRate, _ = strconv.ParseFloat(value, 64)
				case 5:
					rateInfo.BuyingRate, _ = strconv.ParseFloat(value, 64)
				case 6:
					rateInfo.CashBuyingRate, _ = strconv.ParseFloat(value, 64)
				case 7:
					value = date + "T" + value
					rateInfo.ReleaseTime = chinaMerchantsBankCrawler.ConvertReleaseTime(value, 2) // 招商银行的银行ID为2
				}
			})

			if len(chinaMerchantsBankExchangeRateMap) == 0 {

				chinaMerchantsBankCrawler.SaveEachBankExchangeRate(rateInfo, "b_china_cmb")
				chinaMerchantsBankExchangeRateMap[rateInfo.CurrencyID] = rateInfo

			} else {

				if rateInfo.CurrencyID != 0 {

					pre := chinaMerchantsBankExchangeRateMap[rateInfo.CurrencyID]

					if pre != nil {

						if isChanged(pre, rateInfo) {

							chinaMerchantsBankCrawler.SaveEachBankExchangeRate(rateInfo, "b_china_cmb")
							chinaMerchantsBankExchangeRateMap[rateInfo.CurrencyID] = rateInfo

							fmt.Println("China Merchants Bank Data changed.")

						} else {

							// 回收内存
							rateInfo = nil
						}

					} else {

						chinaMerchantsBankCrawler.SaveEachBankExchangeRate(rateInfo, "b_china_cmb")
						chinaMerchantsBankExchangeRateMap[rateInfo.CurrencyID] = rateInfo
					}
				}
			}
		})
	})
}
