package dto

import "github.com/wen-git-acc/orderbook/pkg/tarantool_pkg"

// help me generate user request the with body user id and wallet balance
type UserDepositRequest struct {
	UserID        string  `json:"user_id"`
	DepositAmount float64 `json:"deposit_amount"`
}

type UserDepositResponse struct {
	UserID       string  `json:"user_id"`
	WalletAmount float64 `json:"wallet_amount"`
}

type InsertOrderRequest struct {
	UserId       string  `json:"user_id"`
	Price        float64 `json:"price"`
	Market       string  `json:"market"`
	Side         string  `json:"side"`
	PositionSize float64 `json:"position_size"`
}

type InsertOrderResponse struct {
	IsSuccess bool `json:"is_success"`
}

type GetAllOrderResponse struct {
	Orders []*tarantool_pkg.OrderStruct `json:"orders"`
}
