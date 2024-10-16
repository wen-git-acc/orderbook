package tarantool_pkg

import (
	"github.com/tarantool/go-tarantool/v2"
)

type TarantoolUserConnInterface interface {
	IsUserRegistered(userID string) bool
	GetUserWalletBalance(user_id string) float64
	UpdateUserWalletBalance(user_id string, balance float64) error
	CreateUserWalletBalance(user_id string, balance float64) error
}

const (
	getUserWalletBalance    = "get_user_wallet_balance"
	updateUserWalletBalance = "update_user_wallet_balance"
	createUserWalltBalance  = "create_user_wallet_balance"
)

func (c *TarantoolClient) IsUserRegistered(userID string) bool {
	// Check if user is registered
	conn := c.conn
	result, err := conn.Do(
		tarantool.NewCallRequest(getUserWalletBalance).Args([]interface{}{userID}), // Ensure this matches the space format
	).Get()

	if err != nil {
		c.logger.Error("Got an error:", err)
	}

	var data interface{}
	if len(result) > 0 {
		data = result[0]
	}

	if data != nil {
		return true
	}

	return false
}

func (c *TarantoolClient) GetUserWalletBalance(userId string) float64 {
	// Get user wallet balance
	conn := c.conn
	result, err := conn.Do(
		tarantool.NewCallRequest(getUserWalletBalance).Args([]interface{}{userId}), // Ensure this matches the space format
	).Get()

	if err != nil {
		c.logger.Error("Got an error:", err)
	}

	var data interface{}
	if len(result) > 0 {
		data = result[0]
	}
	if data != nil {
		return c.convertToFloat64(data)
	}

	return 0
}

func (c *TarantoolClient) UpdateUserWalletBalance(userId string, balance float64) error {
	// Update user wallet balance
	conn := c.conn
	_, err := conn.Do(
		tarantool.NewCallRequest(updateUserWalletBalance).Args([]interface{}{userId, balance}), // Ensure this matches the space format
	).Get()

	if err != nil {
		c.logger.Error("Got an error:", err)
	}

	return err
}

func (c *TarantoolClient) CreateUserWalletBalance(userId string, balance float64) error {
	// Update user wallet balance
	conn := c.conn
	_, err := conn.Do(
		tarantool.NewCallRequest(createUserWalltBalance).Args([]interface{}{userId, balance}), // Ensure this matches the space format
	).Get()

	if err != nil {
		c.logger.Error("Got an error:", err)
	}

	return err
}
