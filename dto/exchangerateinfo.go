package dto

import (
	"time"
)

// ExchangeRateInfo 汇率信息
type ExchangeRateInfo struct {
	ID                int       // 数据记录ID
	BankID            int       // 银行代码
	CurrencyID        int       // 货币代码
	CurrencyName      string    // 货币名称
	BuyingRate        float64   // 现汇买入价
	CashBuyingRate    float64   // 现钞买入价
	SellingRate       float64   // 现汇卖出价
	CashSellingRate   float64   // 现钞卖出价
	MiddleRate        float64   // 折算价
	Benchmark         float64   // 基准价
	CentralParityRate float64   // 汇率中间价
	ReferenceRate     float64   // 参考汇率
	ReleaseTime       time.Time // 发布时间
	CreateTime        time.Time // 创建时间
}

// CurrencyInfo 货币信息
type CurrencyInfo struct {
	ID             int       // ID
	Names          string    // 名称
	CreateTime     time.Time // 创建时间
	LastUpdateTime time.Time // 最后更新时间
}

// BanksInfo 银行信息
type BanksInfo struct {
	ID         int       // ID
	BankName   string    // 银行名称
	BankNameEn string    // 银行英文名称
	TimeZone   string    // 所在时区
	CreateTime time.Time // 创建时间
}
