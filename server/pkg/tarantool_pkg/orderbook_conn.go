package tarantool_pkg

import (
	"fmt"
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
	createdTime  int64
}

const (
	orderBookSpace       = "order_book"
	insertOrderData      = "insert_order_data"
	getOrderByPrimaryKey = "get_order_by_primary_key"
	marketSideIndex      = "market_side_index"
)

type TarantoolOrderBookConnInterface interface {
	GetAllOrders() []*OrderStruct
	InsertNewOrder(order *OrderStruct) error
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

	// primaryKey := fmt.Sprintf("%s:%.2f:%d:%s", userId, price, side, market)
	// sideStr := fmt.Sprintf("%d", side)
	// Update user wallet balance

	_, err := conn.Do(
		tarantool.NewCallRequest(insertOrderData).Args([]interface{}{primaryKey, order.Price, order.Market, order.Side, order.UserId, order.PositionSize, timestamp}), // Ensure this matches the space format
	).Get()

	if err != nil {
		fmt.Println("Got an error:", err)
	}
	return err
}

func (c *TarantoolClient) GetOrderByPrimaryKey(userId string, price float64, side int, market string) error {
	primaryKey := fmt.Sprintf("%s:%.2f:%d:%s", userId, price, side, market)
	conn := c.conn

	// Update user wallet balance
	result, err := conn.Do(
		tarantool.NewCallRequest(getOrderByPrimaryKey).Args([]interface{}{primaryKey}), // Ensure this matches the space format
	).Get()

	if err != nil {
		fmt.Println("Got an error:", err)
	}
	fmt.Println("result", result)

	return err
}

func (c *TarantoolClient) GetAllOrders() []*OrderStruct {
	conn := c.conn

	result, err := conn.Do(
		tarantool.NewSelectRequest(orderBookSpace).
			Iterator(tarantool.IterAll).
			Key([]interface{}{}), // Ensure this matches the space format
	).Get()

	if err != nil {
		fmt.Println("Got an error:", err)
	}

	// orderList := []*OrderStruct{}
	// for _, item := range result {
	// 	data, ok := item.([]interface{})
	// 	if !ok {
	// 		fmt.Println("Unexpected data format")
	// 		continue
	// 	}
	// 	fmt.Println("data", data)
	// 	orderList = append(orderList, c.transformToOrderStruct(data))
	// }
	orderList := c.getTransformOrderList(result)
	fmt.Println("result", result)

	return orderList
}
func (c *TarantoolClient) getTransformOrderList(data []interface{}) []*OrderStruct {
	orderList := []*OrderStruct{}
	for _, item := range data {
		data, ok := item.([]interface{})
		if !ok {
			fmt.Println("Unexpected data format for order list")
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
		createdTime:  int64(c.convertToInt(data[6])),
	}
}

// with index request
func (c *TarantoolClient) GetOrdersByMarketAndSide(market string, side int) error {
	conn := c.conn

	sideStr := fmt.Sprintf("%d", side)
	result, err := conn.Do(
		tarantool.NewSelectRequest(orderBookSpace).
			Index(marketSideIndex).
			Key([]interface{}{market, sideStr}).
			Iterator(tarantool.IterEq),
	).Get()

	if err != nil {
		fmt.Println("Got an error:", err)
	}
	fmt.Println("result", result)

	return err
}
