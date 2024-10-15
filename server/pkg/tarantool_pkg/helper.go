package tarantool_pkg

import (
	"fmt"
	"strconv"
)

func (c *TarantoolClient) convertToFloat64(data interface{}) float64 {
	if balance, ok := data.(float64); ok {
		return balance
	}
	strnum := fmt.Sprintf("%v", data)
	floatNum, err := strconv.ParseFloat(strnum, 64)

	if err != nil {
		return 0
	}
	return floatNum
}

func (c *TarantoolClient) convertToInt(data interface{}) int {
	if intValue, ok := data.(int); ok {
		return intValue
	}
	strnum := fmt.Sprintf("%v", data)
	intNum, err := strconv.Atoi(strnum)

	if err != nil {
		return 0
	}
	return intNum
}

func (c *TarantoolClient) calculateAccountMargin(accountEquity float64, totalAccountNotional float64) float64 {
	return accountEquity / totalAccountNotional
}

func (c *TarantoolClient) calculateAccountEquity(walletBalance float64, positions []*PositionStruct) float64 {
	equity := walletBalance
	for _, position := range positions {
		marketPrice := c.GetMarketPriceByMarket(position.Market)
		if position.Side == "1" {
			equity += position.PositionSize * (marketPrice - position.AvgPrice)
		} else {
			equity += position.PositionSize * (position.AvgPrice - marketPrice)
		}
	}
	return equity
}

func (c *TarantoolClient) calculateTotalAccountNotional(positions []*PositionStruct) float64 {
	totalNotional := 0.0
	for _, position := range positions {
		totalNotional += position.PositionSize * c.GetMarketPriceByMarket(position.Market)
	}
	return totalNotional
}
