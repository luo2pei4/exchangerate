package db

import (
	"database/sql"
	dto "exchangerate/dto"
	"fmt"
	"strings"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Connection 数据库连接
type Connection struct {
	conn *sql.DB
}

// Close 关闭连接
func (c *Connection) Close() {
	c.Close()
}

// NewConnection 获取数据库连接
func NewConnection() (connection *Connection, err error) {

	c, err := sql.Open("mysql", "dbo:caecaodb@tcp(192.168.3.168:3306)/exrate?charset=utf8&parseTime=true")

	if err != nil {
		return nil, err
	}

	fmt.Println("db connect successful.")

	connection = &Connection{
		conn: c,
	}

	return
}

// GetCurrencyMap 获取货币信息Map
func GetCurrencyMap(c *Connection) map[string]int {

	sql := "select id, names from currency"
	rows, err := c.conn.Query(sql)

	if err != nil {
		fmt.Println("Get currency info failed.")
		return nil
	}

	var currencyInfoSlice []dto.CurrencyInfo

	for rows.Next() {

		currencyInfo := dto.CurrencyInfo{}

		if err := rows.Scan(&currencyInfo.ID, &currencyInfo.Names); err != nil {
			fmt.Println("Rows scan failed.")
			return nil
		}

		currencyInfoSlice = append(currencyInfoSlice, currencyInfo)
	}

	currencyInfoMap := make(map[string]int)

	for _, ci := range currencyInfoSlice {

		names := ci.Names
		nameList := strings.Split(names, ",")

		for _, name := range nameList {

			currencyInfoMap[name] = ci.ID
		}
	}

	fmt.Println("Get currency map successful.")

	return currencyInfoMap
}

// GetBankTimeZoneMap 获取银行所在时区Map
func GetBankTimeZoneMap(c *Connection) map[int]string {

	sql := "select id, time_zone from banks"
	rows, err := c.conn.Query(sql)

	if err != nil {
		fmt.Println("Get banks info failed.")
		return nil
	}

	var banksInfoSlice []dto.BanksInfo

	for rows.Next() {

		banksInfo := dto.BanksInfo{}

		if err := rows.Scan(&banksInfo.ID, &banksInfo.TimeZone); err != nil {
			fmt.Println("Rows scan failed.")
			return nil
		}

		banksInfoSlice = append(banksInfoSlice, banksInfo)
	}

	bankTimeZoneMap := make(map[int]string)

	for _, bi := range banksInfoSlice {

		bankTimeZoneMap[bi.ID] = bi.TimeZone
	}

	fmt.Println("Get bank time zone map successful.")

	return bankTimeZoneMap
}

// InsertExchangeRateData 插入汇率信息表
func (c *Connection) InsertExchangeRateData(table string, info *dto.ExchangeRateInfo) (lastInsertID, rowsAffected int64, err error) {

	sql := fmt.Sprintf("INSERT INTO %v(currencyid, buying_rate, cash_buying_rate, selling_rate, cash_selling_rate, middle_rate, benchmark, central_parity_rate, reference_rate, release_time) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", table)
	ins, err := c.conn.Prepare(sql)

	if err != nil {
		return 0, 0, err
	}

	result, err := ins.Exec(info.CurrencyID, info.BuyingRate, info.CashBuyingRate, info.SellingRate, info.CashSellingRate, info.MiddleRate, info.Benchmark, info.CentralParityRate, info.ReferenceRate, info.ReleaseTime)

	if err != nil {
		fmt.Println(err.Error())
	}

	lastInsertID, _ = result.LastInsertId()
	rowsAffected, _ = result.RowsAffected()

	return
}

// SelectExchangeRateInfoMap 查询初始化时获取各个银行最后一次的汇率数据
func (c *Connection) SelectExchangeRateInfoMap(table string) (exchangeRateInfoMap map[int]*dto.ExchangeRateInfo, err error) {

	sql := `
		select 
			id, 
			currencyid,
			buying_rate,
			cash_buying_rate,
			selling_rate,
			cash_selling_rate,
			middle_rate,
			benchmark,
			central_parity_rate,
			reference_rate,
			release_time,
			createtime
		from (
			select 
				id, 
				currencyid,
				buying_rate,
				cash_buying_rate,
				selling_rate,
				cash_selling_rate,
				middle_rate,
				benchmark,
				central_parity_rate,
				reference_rate,
				release_time,
				createtime, 
				rank()over(partition by currencyid order by release_time desc) m 
			from 
				%v
		) temp
		where m = 1`

	sql = fmt.Sprintf(sql, table)
	rows, err := c.conn.Query(sql)

	if err != nil {
		fmt.Println("Get exchange rate Info failed.")
		return nil, err
	}

	exchangeRateInfoMap = make(map[int]*dto.ExchangeRateInfo)

	for rows.Next() {

		info := dto.ExchangeRateInfo{}
		err = rows.Scan(
			&info.ID,
			&info.CurrencyID,
			&info.BuyingRate,
			&info.CashBuyingRate,
			&info.SellingRate,
			&info.CashSellingRate,
			&info.MiddleRate,
			&info.Benchmark,
			&info.CentralParityRate,
			&info.ReferenceRate,
			&info.ReleaseTime,
			&info.CreateTime,
		)

		exchangeRateInfoMap[info.CurrencyID] = &info
	}

	return
}
