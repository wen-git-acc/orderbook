package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tarantool/go-tarantool/v2"
	_ "github.com/tarantool/go-tarantool/v2/datetime"
	//	_ "github.com/tarantool/go-tarantool/v2/decimal"
	//	_ "github.com/tarantool/go-tarantool/v2/uuid"
	//
	// )
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	dialer := tarantool.NetDialer{
		Address:  "127.0.0.1:3301",
		User:     "sampleuser",
		Password: "123456",
	}
	opts := tarantool.Opts{
		Timeout: time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		fmt.Println("Connection refused:", err)
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) // Listen for interrupt (Ctrl + C) and termination signals

	go func() {
		<-sigs
		fmt.Println("Got interrupt signal. Close connection")
		conn.CloseGraceful()
		cancel()
	}()
	defer func() {
		conn.CloseGraceful()
		fmt.Println("Connection is closed")
	}()

	//Test Ground

	//Market Price
	// exampleUpsert(conn)
	// exampeCallStoredProcedure(conn)

	//Users
	// isRegisted := isUserRegistered(conn, "123")
	// fmt.Println("isRegisted", isRegisted)

	// if !isRegisted {
	// 	createUserWalletBalance(conn, "123", 1000.53)
	// }

	// balance := getUserWalletBalance(conn, "123")
	// fmt.Println("balance", balance)

	// updateUserWalletBalance(conn, "123", 2023123.31)
	// balance = getUserWalletBalance(conn, "123")
	// fmt.Printf("balance %.2f", balance)

	//Orderbooks
	// insertNewOrder(conn, 90.2, "BTC", 1, "userC", 2.0)
	// time.Sleep(1 * time.Second)
	// insertNewOrder(conn, 123.12, "BTC", 1, "userA", 10.0)
	// time.Sleep(1 * time.Second)

	// insertNewOrder(conn, 123.12, "BTC", 1, "userC", 3.0)
	// time.Sleep(1 * time.Second)

	// insertNewOrder(conn, 300.2, "BTC", 1, "userD", 2.0)
	// time.Sleep(1 * time.Second)

	// insertNewOrder(conn, 90.2, "BTC", -1, "userC", 2.0)
	// time.Sleep(1 * time.Second)

	// insertNewOrder(conn, 123.12, "BTC", -1, "userA", 10.0)
	// time.Sleep(1 * time.Second)

	// insertNewOrder(conn, 123.12, "BTC", -1, "userC", 3.0)
	// time.Sleep(1 * time.Second)

	// insertNewOrder(conn, 300.2, "BTC", -1, "userD", 2.0)
	// time.Sleep(1 * time.Second)

	// insertNewOrder(conn, 90.2, "ETH", 1, "userC", 2.0)
	// time.Sleep(1 * time.Second)

	// insertNewOrder(conn, 123.12, "ETH", 1, "userA", 10.0)
	// time.Sleep(1 * time.Second)

	// insertNewOrder(conn, 123.12, "ETH", 1, "userC", 3.0)
	// time.Sleep(1 * time.Second)

	// insertNewOrder(conn, 300.2, "ETH", 1, "userD", 2.0)
	// time.Sleep(1 * time.Second)

	// insertNewOrder(conn, 90.2, "ETH", -1, "userC", 2.0)
	// time.Sleep(1 * time.Second)

	// insertNewOrder(conn, 123.12, "ETH", -1, "userA", 10.0)
	// time.Sleep(1 * time.Second)

	// insertNewOrder(conn, 123.12, "ETH", -1, "userC", 3.0)
	// time.Sleep(1 * time.Second)

	// insertNewOrder(conn, 300.2, "ETH", -1, "userD", 2.0)
	// time.Sleep(1 * time.Second)

	// getAllOrders(conn)
	getOrdersByMarketAndSide(conn, "ETH", -1)
	// getOrder(conn, "123", 123.12)

	// getOrder(conn, "123", 123.12)
	// getAskOrderBooksByPriceSelect2(conn, "BTC", 1, 150.00)

	// getAskOrderBooksByPrice(conn, "BTC", 1, 123.00)
}

// Market price related
func exampleUpsert(conn *tarantool.Connection) {
	// Upsert a market price
	result, err := conn.Do(
		tarantool.NewUpsertRequest("market_price").
			Tuple([]interface{}{"BTC", 1233.12}).                     // Ensure this matches the space format
			Operations(tarantool.NewOperations().Assign(1, 1233.16)), // Update price
	).Get()

	if err != nil {
		fmt.Println("hsds")
		fmt.Println("Got an error:", err)
	}
	fmt.Println("Stored  result:", result)

}

func exampeCallStoredProcedure(conn *tarantool.Connection) {
	s, err := conn.Do(
		tarantool.NewCallRequest("get_market_price").Args([]interface{}{"BTC"}),
	).Get()
	if err != nil {
		fmt.Println("why")
		fmt.Println("Got an error:", err)
	}
	fmt.Println("Stored procedure result:", s)

	if len(s) > 0 {
		fmt.Println("Stored procedure result:", s[0])
	}
}

// Users
func isUserRegistered(conn *tarantool.Connection, user_id string) bool {
	// Check if user is registered
	result, err := conn.Do(
		tarantool.NewCallRequest("get_user_wallet_balance").Args([]interface{}{user_id}), // Ensure this matches the space format
	).Get()

	data := result[0]
	if err != nil {
		fmt.Println("Got an error:", err)
	}
	if data != nil {
		return true
	}

	return false
}

func getUserWalletBalance(conn *tarantool.Connection, user_id string) float64 {
	// Get user wallet balance
	result, err := conn.Do(
		tarantool.NewCallRequest("get_user_wallet_balance").Args([]interface{}{user_id}), // Ensure this matches the space format
	).Get()

	if err != nil {
		fmt.Println("Got an error:", err)
	}

	if len(result) > 0 {
		return result[0].(float64)
	}

	return 0
}

func updateUserWalletBalance(conn *tarantool.Connection, user_id string, balance float64) error {
	// Update user wallet balance
	_, err := conn.Do(
		tarantool.NewCallRequest("update_user_wallet_balance").Args([]interface{}{user_id, balance}), // Ensure this matches the space format
	).Get()

	if err != nil {
		fmt.Println("Got an error:", err)
	}

	return err
}

func createUserWalletBalance(conn *tarantool.Connection, user_id string, balance float64) error {
	// Update user wallet balance
	_, err := conn.Do(
		tarantool.NewCallRequest("create_user_wallet_balance").Args([]interface{}{user_id, balance}), // Ensure this matches the space format
	).Get()

	if err != nil {
		fmt.Println("Got an error:", err)
	}

	return err
}

//Orderbook

func insertNewOrder(conn *tarantool.Connection, price float64, market string, side int, userId string, positionSize float64) error {
	timestamp := time.Now().Unix()
	fmt.Println("timestamp", timestamp)
	primaryKey := fmt.Sprintf("%s:%.2f:%d:%s", userId, price, side, market)
	sideStr := fmt.Sprintf("%d", side)
	// Update user wallet balance
	_, err := conn.Do(
		tarantool.NewCallRequest("insert_order_data").Args([]interface{}{primaryKey, price, market, sideStr, userId, positionSize, timestamp}), // Ensure this matches the space format
	).Get()

	if err != nil {
		fmt.Println("Got an error:", err)
	}
	return err
}

// func getOrder(conn *tarantool.Connection, userId string, price float64) error {
// 	// Update user wallet balance
// 	result, err := conn.Do(
// 		tarantool.NewCallRequest("get_order_by_price_and_user_id").Args([]interface{}{userId, price}), // Ensure this matches the space format
// 	).Get()

// 	if err != nil {
// 		fmt.Println("Got an error:", err)
// 	}
// 	fmt.Println("result", result)

// 	return err
// }

func getOrderByPrimaryKey(conn *tarantool.Connection, userId string, price float64, side int, market string) error {
	primaryKey := fmt.Sprintf("%s:%.2f:%d:%s", userId, price, side, market)

	// Update user wallet balance
	result, err := conn.Do(
		tarantool.NewCallRequest("get_order_by_primary_key").Args([]interface{}{primaryKey}), // Ensure this matches the space format
	).Get()

	if err != nil {
		fmt.Println("Got an error:", err)
	}
	fmt.Println("result", result)

	return err
}

// Get sell roder book by select request using index
// Side 1 means get ask price
// Side -1 means get bid price
func getAllOrders(conn *tarantool.Connection) error {
	result, err := conn.Do(
		tarantool.NewSelectRequest("order_book").
			Iterator(tarantool.IterAll).
			Key([]interface{}{}), // Ensure this matches the space format
	).Get()

	if err != nil {
		fmt.Println("Got an error:", err)
	}
	fmt.Println("result", result)

	return err
}

// with index request
func getOrdersByMarketAndSide(conn *tarantool.Connection, market string, side int) error {
	sideStr := fmt.Sprintf("%d", side)
	result, err := conn.Do(
		tarantool.NewSelectRequest("order_book").
			Index("market_side_index").
			Key([]interface{}{market, sideStr}).
			Iterator(tarantool.IterEq),
	).Get()

	if err != nil {
		fmt.Println("Got an error:", err)
	}
	fmt.Println("result", result)

	return err
}
