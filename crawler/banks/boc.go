// Package banks 中国银行
package banks

import (
	"exchangerate/crawler"
	dto "exchangerate/dto"
	"fmt"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var bocCrawler *crawler.Crawler
var bocExchangeRateMap map[int]*dto.ExchangeRateInfo

// BankOfChina 爬取数据
func BankOfChina() {

	if bocCrawler == nil {

		c, err := crawler.NewInstance()

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		bocCrawler = c
	}

	doc := bocCrawler.GetPageDoc("https://www.boc.cn/sourcedb/whpj/index.html")

	if doc == nil {
		return
	}

	if bocExchangeRateMap == nil {

		erMap, err := bocCrawler.InitExchangeRateInfoMap("b_china_boc")

		if err != nil {
			fmt.Printf("Init exchange rate info map failed. %v\n", err.Error())
			return
		}

		bocExchangeRateMap = erMap
	}

	doc.Find("body > div > div.BOC_main > div.publish > div:nth-child(3) > table > tbody").Each(func(i int, table *goquery.Selection) {

		table.Find("tr").Each(func(j int, tr *goquery.Selection) {

			rateInfo := &dto.ExchangeRateInfo{}

			tr.Find("td").Each(func(k int, td *goquery.Selection) {

				value := td.Text()

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
					rateInfo.MiddleRate, _ = strconv.ParseFloat(value, 64)
				case 6:
					rateInfo.ReleaseTime = bocCrawler.ConvertReleaseTime(value, 1) // 中国银行的银行ID为1
				}
			})

			if len(bocExchangeRateMap) == 0 {

				bocCrawler.SaveEachBankExchangeRate(rateInfo, "b_china_boc")
				bocExchangeRateMap[rateInfo.CurrencyID] = rateInfo

			} else {

				if rateInfo.CurrencyID != 0 {

					pre := bocExchangeRateMap[rateInfo.CurrencyID]

					if pre != nil {

						if isChanged(pre, rateInfo) {

							bocCrawler.SaveEachBankExchangeRate(rateInfo, "b_china_boc")
							bocExchangeRateMap[rateInfo.CurrencyID] = rateInfo

							fmt.Println("Bank of China Data changed.")
						} else {

							// 回收内存
							rateInfo = nil
						}

					} else {

						bocCrawler.SaveEachBankExchangeRate(rateInfo, "b_china_boc")
						bocExchangeRateMap[rateInfo.CurrencyID] = rateInfo
					}
				}
			}
		})
	})
}

func isChanged(pre *dto.ExchangeRateInfo, now *dto.ExchangeRateInfo) bool {

	if pre.BuyingRate != now.BuyingRate {
		return true
	}

	if pre.CashBuyingRate != now.CashBuyingRate {
		return true
	}

	if pre.CashBuyingRate != now.CashBuyingRate {
		return true
	}

	if pre.SellingRate != now.SellingRate {
		return true
	}

	if pre.CashSellingRate != now.CashSellingRate {
		return true
	}

	if pre.MiddleRate != now.MiddleRate {
		return true
	}

	if pre.Benchmark != now.Benchmark {
		return true
	}

	if pre.CentralParityRate != now.CentralParityRate {
		return true
	}

	if pre.ReferenceRate != now.ReferenceRate {
		return true
	}

	return false
}
