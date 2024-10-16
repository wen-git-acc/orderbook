package tarantool_pkg

import (
	"fmt"
	"strconv"
)

type HelperInterface interface {
	CalculateAccountEquity(walletBalance float64, positions []*PositionStruct) float64
	CalculateAccountMargin(accountEquity float64, totalAccountNotional float64) float64
	CalculateTotalAccountNotional(positions []*PositionStruct) float64
}

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

func (c *TarantoolClient) CalculateAccountMargin(accountEquity float64, totalAccountNotional float64) float64 {
	if totalAccountNotional == 0 {
		if accountEquity > 0 {
			return float64(^uint(0) >> 1) // Return the maximum float64 value (represents high margin)
		}
		return 0 //if both zero
	}
	return accountEquity / totalAccountNotional
}

func (c *TarantoolClient) CalculateAccountEquity(walletBalance float64, positions []*PositionStruct) float64 {
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

func (c *TarantoolClient) CalculateTotalAccountNotional(positions []*PositionStruct) float64 {
	totalNotional := 0.0
	for _, position := range positions {
		totalNotional += position.PositionSize * c.GetMarketPriceByMarket(position.Market)
	}
	return totalNotional
}
