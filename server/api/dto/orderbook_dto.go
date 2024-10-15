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

type GetAllPositionsResposne struct {
	Positions []*tarantool_pkg.PositionStruct `json:"positions"`
}

type GetMarketPriceResponse struct {
	MarketPrice float64 `json:"market_price"`
}

type DeleteOrderRequest struct {
	UserId string  `json:"user_id"`
	Price  float64 `json:"price"`
	Side   string  `json:"side"`
	Market string  `json:"market"`
}

type DeleteOrderResponse struct {
	IsSuccess bool `json:"is_success"`
}

type GetUserPositionResponse struct {
	Positions []*tarantool_pkg.PositionStruct `json:"positions"`
}
