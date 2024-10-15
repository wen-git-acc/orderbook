package tarantool_pkg

import (
	"fmt"

	"github.com/tarantool/go-tarantool/v2"
)

type TarantoolUserConnInterface interface {
	IsUserRegistered(userID string) bool
	GetUserWalletBalance(user_id string) float64
	UpdateUserWalletBalance(user_id string, balance float64) error
	CreateUserWalletBalance(user_id string, balance float64) error
}

func (c *TarantoolClient) IsUserRegistered(userID string) bool {
	// Check if user is registered
	conn := c.conn
	result, err := conn.Do(
		tarantool.NewCallRequest("get_user_wallet_balance").Args([]interface{}{userID}), // Ensure this matches the space format
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

func (c *TarantoolClient) GetUserWalletBalance(user_id string) float64 {
	// Get user wallet balance
	conn := c.conn
	result, err := conn.Do(
		tarantool.NewCallRequest("get_user_wallet_balance").Args([]interface{}{user_id}), // Ensure this matches the space format
	).Get()

	if err != nil {
		fmt.Println("Got an error:", err)
	}

	data := result[0]
	if data != nil {
		return c.convertToFloat64(data)
	}

	return 0
}

func (c *TarantoolClient) UpdateUserWalletBalance(user_id string, balance float64) error {
	// Update user wallet balance
	conn := c.conn
	_, err := conn.Do(
		tarantool.NewCallRequest("update_user_wallet_balance").Args([]interface{}{user_id, balance}), // Ensure this matches the space format
	).Get()

	if err != nil {
		fmt.Println("Got an error:", err)
	}

	return err
}

func (c *TarantoolClient) CreateUserWalletBalance(user_id string, balance float64) error {
	// Update user wallet balance
	conn := c.conn
	_, err := conn.Do(
		tarantool.NewCallRequest("create_user_wallet_balance").Args([]interface{}{user_id, balance}), // Ensure this matches the space format
	).Get()

	if err != nil {
		fmt.Println("Got an error:", err)
	}

	return err
}
