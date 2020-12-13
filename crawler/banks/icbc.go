// Package banks 中国工商银行
package banks

import (
	"exchangerate/crawler"
	dto "exchangerate/dto"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var icbcCrawler *crawler.Crawler
var icbcExchangeRateMap map[int]*dto.ExchangeRateInfo

// IndustrialAndCommercialBankOfChina 爬取数据
func IndustrialAndCommercialBankOfChina() {

	if icbcCrawler == nil {

		c, err := crawler.NewInstance()

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		icbcCrawler = c
	}

	doc := icbcCrawler.GetPageDoc("http://www.icbc.com.cn/ICBCDynamicSite/Optimize/Quotation/QuotationListIframe.aspx")

	if doc == nil {
		return
	}

	if icbcExchangeRateMap == nil {

		erMap, err := icbcCrawler.InitExchangeRateInfoMap("b_china_icbc")

		if err != nil {
			fmt.Println("Init exchange rate info map failed.")
			return
		}

		icbcExchangeRateMap = erMap
	}

	doc.Find("#form1 > div > table > tbody").Each(func(i int, table *goquery.Selection) {

		table.Find("tr").Each(func(j int, tr *goquery.Selection) {

			rateInfo := &dto.ExchangeRateInfo{}

			tr.Find("td").Each(func(k int, td *goquery.Selection) {

				value := td.Text()
				value = strings.ReplaceAll(value, "\n", "")
				value = strings.Trim(value, " ")

				if value == "--" {
					value = "0"
				}

				switch k {
				case 0:
					rateInfo.CurrencyName = value
					currencyID := crawler.CurrencyInfoMap[rateInfo.CurrencyName]
					rateInfo.CurrencyID = currencyID
				case 1:
					rateInfo.BuyingRate, _ = strconv.ParseFloat(value, 64)
				case 2:
					rateInfo.CashBuyingRate, _ = strconv.ParseFloat(value, 64)
				case 3:
					rateInfo.SellingRate, _ = strconv.ParseFloat(value, 64)
				case 4:
					rateInfo.CashSellingRate, _ = strconv.ParseFloat(value, 64)
				case 5:
					rateInfo.ReleaseTime = icbcCrawler.ConvertReleaseTime(value, 3)
				}
			})

			fmt.Println(rateInfo)

			if len(icbcExchangeRateMap) == 0 {

				icbcCrawler.SaveEachBankExchangeRate(rateInfo, "b_china_icbc")
				icbcExchangeRateMap[rateInfo.CurrencyID] = rateInfo

			} else {

				if rateInfo.CurrencyID != 0 {

					pre := icbcExchangeRateMap[rateInfo.CurrencyID]

					if pre != nil {

						if isChanged(pre, rateInfo) {

							icbcCrawler.SaveEachBankExchangeRate(rateInfo, "b_china_icbc")
							icbcExchangeRateMap[rateInfo.CurrencyID] = rateInfo

							fmt.Println("China Merchants Bank Data changed.")

						} else {

							// 回收内存
							rateInfo = nil
						}

					} else {

						icbcCrawler.SaveEachBankExchangeRate(rateInfo, "b_china_icbc")
						icbcExchangeRateMap[rateInfo.CurrencyID] = rateInfo
					}
				}
			}
		})
	})
}
