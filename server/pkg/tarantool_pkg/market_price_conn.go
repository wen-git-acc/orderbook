package tarantool_pkg

import (
	"github.com/tarantool/go-tarantool/v2"
)

const (
	marketPriceSpace = "market_price"
	getMarketPrice   = "get_market_price"
)

type TarantoolMarketPriceConnInterface interface {
	GetMarketPriceByMarket(market string) float64
	UpdateMarketPrice(market string, marketPrice float64)
}

func (c *TarantoolClient) UpdateMarketPrice(market string, marketPrice float64) {
	conn := c.conn

	// Upsert a market price
	_, err := conn.Do(
		tarantool.NewUpsertRequest(marketPriceSpace).
			Tuple([]interface{}{market, marketPrice}).
			Operations(tarantool.NewOperations().Assign(1, marketPrice)),
	).Get()

	if err != nil {
		c.logger.Error("failed to upsert market price", err)
	}
}

func (c *TarantoolClient) GetMarketPriceByMarket(market string) float64 {
	conn := c.conn
	result, err := conn.Do(
		tarantool.NewCallRequest(getMarketPrice).Args([]interface{}{market}),
	).Get()
	if err != nil {
		c.logger.Error("failed to upsert market price", err)
	}

	if len(result) > 0 {
		return c.convertToFloat64(result[0])
	}

	return 0
}
