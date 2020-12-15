// Package banks 中国建设银行
package banks

import (
	"exchangerate/crawler"
	dto "exchangerate/dto"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var ccbCrawler *crawler.Crawler
var ccbExchangeRateMap map[int]*dto.ExchangeRateInfo

// ChinaConstructionBank 爬取数据
func ChinaConstructionBank() {

	if ccbCrawler == nil {

		c, err := crawler.NewInstance()

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		ccbCrawler = c
	}

	doc := ccbCrawler.GetPageDoc("http://forex.ccb.com/cn/forex/exchange-quotations.html")

	if doc == nil {
		return
	}

	if ccbExchangeRateMap == nil {

		erMap, err := ccbCrawler.InitExchangeRateInfoMap("b_china_ccb")

		if err != nil {
			fmt.Println("Init exchange rate info map failed.")
			return
		}

		ccbExchangeRateMap = erMap
	}

	selection := doc.Find("#jshckpj")
	nodes := selection.Children()

	fmt.Println(len(nodes.Nodes))

	doc.Find("#jshckpj").Each(func(i int, div *goquery.Selection) {

		div.Find("ul").Each(func(j int, ul *goquery.Selection) {

			rateInfo := &dto.ExchangeRateInfo{}

			ul.Find("li").Each(func(k int, li *goquery.Selection) {

				value := li.Text()
				value = strings.ReplaceAll(value, "\n", "")
				value = strings.Trim(value, " ")

				switch k {
				case 0:
					rateInfo.CurrencyName = value
					currencyID := crawler.CurrencyInfoMap[rateInfo.CurrencyName]
					rateInfo.CurrencyID = currencyID
				case 1:
					rateInfo.BuyingRate, _ = strconv.ParseFloat(value, 64)
				case 2:
					rateInfo.SellingRate, _ = strconv.ParseFloat(value, 64)
				case 3:
					rateInfo.CashBuyingRate, _ = strconv.ParseFloat(value, 64)
				case 4:
					rateInfo.CashSellingRate, _ = strconv.ParseFloat(value, 64)
				case 5:
					rateInfo.ReleaseTime = ccbCrawler.ConvertReleaseTime(value, 4)
				}
			})

			fmt.Println(rateInfo)

			if len(ccbExchangeRateMap) == 0 {

				ccbCrawler.SaveEachBankExchangeRate(rateInfo, "b_china_ccb")
				ccbExchangeRateMap[rateInfo.CurrencyID] = rateInfo

			} else {

				if rateInfo.CurrencyID != 0 {

					pre := ccbExchangeRateMap[rateInfo.CurrencyID]

					if pre != nil {

						if isChanged(pre, rateInfo) {

							ccbCrawler.SaveEachBankExchangeRate(rateInfo, "b_china_ccb")
							ccbExchangeRateMap[rateInfo.CurrencyID] = rateInfo

							fmt.Println("China Merchants Bank Data changed.")

						} else {

							// 回收内存
							rateInfo = nil
						}

					} else {

						ccbCrawler.SaveEachBankExchangeRate(rateInfo, "b_china_ccb")
						ccbExchangeRateMap[rateInfo.CurrencyID] = rateInfo
					}
				}
			}
		})
	})
}
