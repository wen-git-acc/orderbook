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
	exampleUpsert(conn)
	exampeCallStoredProcedure(conn)

}

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
}
