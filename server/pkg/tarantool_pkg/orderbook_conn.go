package tarantool_pkg

import (
	"fmt"
	"sort"
	"time"

	"github.com/tarantool/go-tarantool/v2"
)

type OrderStruct struct {
	PrimaryKey   string
	Price        float64
	Market       string
	Side         string
	UserId       string
	PositionSize float64
	CreatedTime  int64
}

type ExecutionDetailsStruct struct {
	UserId                string
	Market                string
	Side                  string
	ExecutionPositionSize float64
	ExecutionPrice        float64
}

type SimplifiedOrderBook struct {
	AskOrderBook [][]float64 `json:"ask"`
	BidOrderBook [][]float64 `json:"bid"`
}

const (
	orderBookSpace       = "order_book"
	insertOrderData      = "insert_order_data"
	getOrderByPrimaryKey = "get_order_by_primary_key"
	marketSideIndex      = "market_side_index"
)

type TarantoolOrderBookConnInterface interface {
	InsertNewOrder(order *OrderStruct) error
	GetOrderBook(market string) *SimplifiedOrderBook
	DeleteOrderByPrimaryKey(userId string, price float64, side string, market string) error
	GetOrderByPrimaryKey(userId string, price float64, side string, market string) *OrderStruct
	UpdateOrderByPrimaryKey(userId string, price float64, side string, market string, positionSize float64) error
	GetAllOrders() []*OrderStruct
	GetOrdersByMarketAndSide(market string, side string) []*OrderStruct
	GetBidOrderBook(currentOrder *OrderStruct) []*OrderStruct
	GetAskOrderBook(currentOrder *OrderStruct) []*OrderStruct
}

func (c *TarantoolClient) GetPrimaryKeyForOrder(order *OrderStruct) string {
	primaryKey := fmt.Sprintf("%s:%.2f:%s:%s", order.UserId, order.Price, order.Side, order.Market)
	order.PrimaryKey = primaryKey
	return primaryKey
}

func (c *TarantoolClient) InsertNewOrder(order *OrderStruct) error {
	conn := c.conn
	timestamp := time.Now().Unix()
	primaryKey := c.GetPrimaryKeyForOrder(order)

	_, err := conn.Do(
		tarantool.NewCallRequest(insertOrderData).Args([]interface{}{primaryKey, order.Price, order.Market, order.Side, order.UserId, order.PositionSize, timestamp}),
	).Get()

	if err != nil {
		c.logger.Error("Got an error:", err)
	}
	return err
}

func (c *TarantoolClient) GetOrderByPrimaryKey(userId string, price float64, side string, market string) *OrderStruct {
	primaryKey := fmt.Sprintf("%s:%.2f:%s:%s", userId, price, side, market)
	conn := c.conn

	// Update user wallet balance
	result, err := conn.Do(
		tarantool.NewCallRequest(getOrderByPrimaryKey).Args([]interface{}{primaryKey}),
	).Get()

	if err != nil {
		c.logger.Error("Got an error:", err)
	}

	orderList := c.getTransformOrderList(result)
	if len(orderList) > 0 {
		return orderList[0]
	}

	return nil
}

func (c *TarantoolClient) UpdateOrderByPrimaryKey(userId string, price float64, side string, market string, positionSize float64) error {
	primaryKey := fmt.Sprintf("%s:%.2f:%s:%s", userId, price, side, market)
	conn := c.conn

	// Update user wallet balance
	_, err := conn.Do(
		tarantool.NewUpdateRequest(orderBookSpace).
			Key([]interface{}{primaryKey}).
			Operations(tarantool.NewOperations().Assign(5, positionSize)),
	).Get()

	if err != nil {
		c.logger.Error("Got an error:", err)
	}

	return err
}

func (c *TarantoolClient) DeleteOrderByPrimaryKey(userId string, price float64, side string, market string) error {
	primaryKey := fmt.Sprintf("%s:%.2f:%s:%s", userId, price, side, market)
	conn := c.conn

	// Update user wallet balance
	_, err := conn.Do(
		tarantool.NewDeleteRequest(orderBookSpace).
			Key([]interface{}{primaryKey}),
	).Get()

	if err != nil {
		c.logger.Error("Got an error:", err)
	}

	return err
}

func (c *TarantoolClient) GetAllOrders() []*OrderStruct {
	conn := c.conn

	result, err := conn.Do(
		tarantool.NewSelectRequest(orderBookSpace).
			Iterator(tarantool.IterAll).
			Key([]interface{}{}),
	).Get()

	if err != nil {
		c.logger.Error("Got an error:", err)
	}
	orderList := c.getTransformOrderList(result)

	return orderList
}

func (c *TarantoolClient) getTransformOrderList(data []interface{}) []*OrderStruct {
	orderList := []*OrderStruct{}
	for _, item := range data {
		if item == nil {
			continue
		}
		data, ok := item.([]interface{})
		if !ok {
			c.logger.Info("Unexpected data format for order list")
			continue
		}
		orderList = append(orderList, c.transformToOrderStruct(data))
	}
	return orderList
}

func (c *TarantoolClient) transformToOrderStruct(data []interface{}) *OrderStruct {
	return &OrderStruct{
		PrimaryKey:   data[0].(string),
		Price:        c.convertToFloat64(data[1]),
		Market:       data[2].(string),
		Side:         data[3].(string),
		UserId:       data[4].(string),
		PositionSize: c.convertToFloat64(data[5]),
		CreatedTime:  int64(c.convertToInt(data[6])),
	}
}

func (c *TarantoolClient) GetOrdersByMarketAndSide(market string, side string) []*OrderStruct {
	conn := c.conn

	result, err := conn.Do(
		tarantool.NewSelectRequest(orderBookSpace).
			Index(marketSideIndex).
			Key([]interface{}{market, side}).
			Iterator(tarantool.IterEq),
	).Get()

	if err != nil {
		c.logger.Error("Got an error:", err)
	}

	orderList := c.getTransformOrderList(result)

	return orderList
}

func (c *TarantoolClient) GetOrderBook(market string) *SimplifiedOrderBook {
	askOrderBook := c.GetAskOrderBook(&OrderStruct{Market: market})
	bidOrderBook := c.GetBidOrderBook(&OrderStruct{Market: market})

	return &SimplifiedOrderBook{
		AskOrderBook: c.tranformOrderBookToPriceAndSize(askOrderBook),
		BidOrderBook: c.tranformOrderBookToPriceAndSize(bidOrderBook),
	}
}

func (c *TarantoolClient) tranformOrderBookToPriceAndSize(orderBook []*OrderStruct) [][]float64 {
	orderBookList := [][]float64{}
	for _, order := range orderBook {
		orderBookList = append(orderBookList, []float64{order.Price, order.PositionSize})
	}

	// Sort the orderBookList by price (first element of each slice)
	sort.Slice(orderBookList, func(i, j int) bool {
		return orderBookList[i][0] < orderBookList[j][0]
	})

	return orderBookList
}

func (c *TarantoolClient) GetBidOrderBook(currentOrder *OrderStruct) []*OrderStruct {
	return c.GetOrdersByMarketAndSide(currentOrder.Market, "1")
}

func (c *TarantoolClient) GetAskOrderBook(currentOrder *OrderStruct) []*OrderStruct {
	return c.GetOrdersByMarketAndSide(currentOrder.Market, "-1")
}
