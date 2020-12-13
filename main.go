package main

import (
	crawler "exchangerate/crawler"
	"exchangerate/crawler/banks"
	"fmt"
	"time"
)

func init() {

	crawler.InitializeInfo()
}

func main() {

	timer := time.Tick(120 * 1e9)

	for {
		select {
		case <-timer:
			fmt.Println("Get exchange rate")
			go banks.BankOfChina()
			go banks.ChinaMerchantsBank()
		}
	}
}
