package matching_engine

import (
	"sort"

	"github.com/wen-git-acc/orderbook/pkg/logger"
	"github.com/wen-git-acc/orderbook/pkg/tarantool_pkg"
	"github.com/wen-git-acc/orderbook/pkg/utils"
)

type MatchingEngineInterface interface {
	// CalculateAccountEquity(walletBalance float64, positions []*tarantool_pkg.PositionStruct) float64
	MatchingEngineForLongOrder(order *tarantool_pkg.OrderStruct, orderBook []*tarantool_pkg.OrderStruct) bool
	MatchingEngineForShortOrder(order *tarantool_pkg.OrderStruct, orderBook []*tarantool_pkg.OrderStruct) bool
	OrderMatcher(order *tarantool_pkg.OrderStruct) bool
}

type MatchingEngineOptions struct {
	Logger    logger.LoggerClientInterface
	Utils     utils.UtilsClientInterface
	Tarantool tarantool_pkg.TarantoolClientInterface
}

type MatchingEngine struct {
	//Placeholder for logger client
	logger    logger.LoggerClientInterface
	utils     utils.UtilsClientInterface
	tarantool tarantool_pkg.TarantoolClientInterface
}

func NewMatchingEngine(opt *MatchingEngineOptions) MatchingEngineInterface {
	return &MatchingEngine{
		logger:    opt.Logger,
		utils:     opt.Utils,
		tarantool: opt.Tarantool,
	}
}

func (c *MatchingEngine) OrderMatcher(order *tarantool_pkg.OrderStruct) bool {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("OrderMatcher panicked: %v", r)
		}
	}()

	isComplete := false
	if order.Side == "1" {
		isComplete = c.MatchingEngineForLongOrder(order, c.tarantool.GetAskOrderBook(order))
	} else {
		isComplete = c.MatchingEngineForShortOrder(order, c.tarantool.GetBidOrderBook(order))
	}

	return isComplete
}

// Pass in ask orderbook
func (c *MatchingEngine) MatchingEngineForLongOrder(order *tarantool_pkg.OrderStruct, orderBook []*tarantool_pkg.OrderStruct) bool {
	market := order.Market

	if len(orderBook) == 0 {
		c.tarantool.InsertNewOrder(order)
		return true
	}

	sortedOrderBook := c.sortOrderBook(orderBook)
	if len(sortedOrderBook) > 0 {
		if sortedOrderBook[0][0].Price > order.Price {
			isMarginSufficient := c.checkAccountMargin(&tarantool_pkg.ExecutionDetailsStruct{
				UserId:                order.UserId,
				Market:                market,
				Side:                  order.Side,
				ExecutionPositionSize: order.PositionSize,
				ExecutionPrice:        order.Price,
			})

			if !isMarginSufficient {
				c.logger.Info("user margin less than 10% when attempting to write orderbook", "orderUser", order.UserId)
				c.deleteOrderFromOrderBook(order)
				return false
			}
			c.tarantool.InsertNewOrder(order)
			return true
		}
	}

	for _, makerOrders := range sortedOrderBook {
		if len(makerOrders) == 0 {
			continue
		}

		for _, makerOrder := range makerOrders {
			if order.PositionSize == 0 {
				//End of order matching
				c.logger.Info("order matching end, fully matched", "orderUser", order.UserId, "side", order.Side)
				return true
			}

			if order.UserId == makerOrder.UserId {
				continue
			}

			executionPrice := makerOrder.Price

			if executionPrice > order.Price {
				//End of order matching
				c.logger.Info("order matching end, partially matched", "orderUser", order.UserId, "side", order.Side)
				return true
			}

			userExecutionDetails := &tarantool_pkg.ExecutionDetailsStruct{
				UserId:         order.UserId,
				ExecutionPrice: executionPrice,
				Market:         market,
				Side:           "1",
			}

			makerExecutionDetails := &tarantool_pkg.ExecutionDetailsStruct{
				UserId:         makerOrder.UserId,
				ExecutionPrice: executionPrice,
				Market:         market,
				Side:           "-1",
			}

			userMatchPositionDetails := &tarantool_pkg.PositionStruct{
				UserID:   order.UserId,
				Market:   market,
				AvgPrice: executionPrice,
				Side:     "1",
			}

			makerMatchPositionDetails := &tarantool_pkg.PositionStruct{
				UserID:   makerOrder.UserId,
				Market:   market,
				AvgPrice: executionPrice,
				Side:     "-1",
			}

			if order.PositionSize > makerOrder.PositionSize {
				executionPositionSize := makerOrder.PositionSize

				userExecutionDetails.ExecutionPositionSize = executionPositionSize
				checkUserMargin := c.checkAccountMargin(userExecutionDetails)
				if !checkUserMargin {
					c.logger.Info("user margin less than 10%", "orderUser", order.UserId)
					c.deleteOrderFromOrderBook(order)
					return false
				}

				makerExecutionDetails.ExecutionPositionSize = executionPositionSize
				checkMakerMargin := c.checkAccountMargin(makerExecutionDetails)
				if !checkMakerMargin {
					c.logger.Info("user margin less than 10%", "maker", makerOrder.UserId)
					c.deleteOrderFromOrderBook(makerOrder)
					continue
				}

				order.PositionSize = order.PositionSize - executionPositionSize

				// Insert Match Position for user and maker
				userMatchPositionDetails.PositionSize = executionPositionSize
				c.tarantool.InsertMatchedPosition(userMatchPositionDetails)
				makerMatchPositionDetails.PositionSize = executionPositionSize
				c.tarantool.InsertMatchedPosition(makerMatchPositionDetails)

				// Delete Maker Order from orderbook
				c.tarantool.DeleteOrderByPrimaryKey(makerOrder.UserId, makerOrder.Price, makerOrder.Side, makerOrder.Market)

				// Update Market Price
				c.tarantool.UpdateMarketPrice(market, executionPrice)
				c.logger.Info("matchingEngineForShortOrder", "executionPositionSize", executionPositionSize, "executionPrice", executionPrice, "orderUser", order.UserId, "maker", makerOrder.UserId)

			} else {
				executionPositionSize := order.PositionSize

				userExecutionDetails.ExecutionPositionSize = executionPositionSize
				checkUserMargin := c.checkAccountMargin(userExecutionDetails)
				if !checkUserMargin {
					c.logger.Info("user margin less than 10%", "orderUser", order.UserId)
					c.deleteOrderFromOrderBook(order)
					return false
				}

				makerExecutionDetails.ExecutionPositionSize = executionPositionSize
				checkMakerMargin := c.checkAccountMargin(makerExecutionDetails)
				if !checkMakerMargin {
					c.logger.Info("user margin less than 10%", "maker", makerOrder.UserId)
					c.deleteOrderFromOrderBook(makerOrder)
					continue
				}

				makerOrder.PositionSize = makerOrder.PositionSize - executionPositionSize
				order.PositionSize = 0

				// Insert Match Position for user and maker
				userMatchPositionDetails.PositionSize = executionPositionSize
				c.tarantool.InsertMatchedPosition(userMatchPositionDetails)
				makerMatchPositionDetails.PositionSize = executionPositionSize
				c.tarantool.InsertMatchedPosition(makerMatchPositionDetails)

				// Update Maker Order from orderbook
				if (makerOrder.PositionSize) == 0 {
					c.tarantool.DeleteOrderByPrimaryKey(makerOrder.UserId, makerOrder.Price, makerOrder.Side, makerOrder.Market)
				} else {
					c.tarantool.UpdateOrderByPrimaryKey(makerOrder.UserId, makerOrder.Price, makerOrder.Side, makerOrder.Market, makerOrder.PositionSize)
				}
				// Update Market Price
				c.tarantool.UpdateMarketPrice(market, executionPrice)
				c.logger.Info("matchingEngineForShortOrder", "executionPositionSize", executionPositionSize, "executionPrice", executionPrice, "orderUser", order.UserId, "maker", makerOrder.UserId)

			}

		}
	}
	return true
}

// Pass in bid orderbook
func (c *MatchingEngine) MatchingEngineForShortOrder(order *tarantool_pkg.OrderStruct, orderBook []*tarantool_pkg.OrderStruct) bool {
	market := order.Market

	if len(orderBook) == 0 {
		return false
	}

	sortedOrderBook := c.sortOrderBook(orderBook)
	if len(sortedOrderBook) > 0 {
		if sortedOrderBook[len(sortedOrderBook)-1][0].Price < order.Price {
			isMarginSufficient := c.checkAccountMargin(&tarantool_pkg.ExecutionDetailsStruct{
				UserId:                order.UserId,
				Market:                market,
				Side:                  order.Side,
				ExecutionPositionSize: order.PositionSize,
				ExecutionPrice:        order.Price,
			})

			if !isMarginSufficient {
				c.logger.Info("user margin less than 10% when attempting to write orderbook", "orderUser", order.UserId)
				c.deleteOrderFromOrderBook(order)
				return false
			}
			c.tarantool.InsertNewOrder(order)
			return true
		}
	}

	//Iterate from last element highest price
	for i := len(sortedOrderBook) - 1; i >= 0; i-- {
		makerOrders := sortedOrderBook[i]
		if len(makerOrders) == 0 {
			continue
		}

		for _, makerOrder := range makerOrders {

			if order.PositionSize == 0 {
				//End of order matching
				c.logger.Info("order matching end, fully matched", "orderUser", order.UserId, "side", order.Side)
				return true
			}

			if order.UserId == makerOrder.UserId {
				continue
			}

			executionPrice := makerOrder.Price

			if executionPrice < order.Price {
				//End of order matching
				c.logger.Info("order matching end, partially matched", "orderUser", order.UserId, "side", order.Side)
				return true
			}

			userExecutionDetails := &tarantool_pkg.ExecutionDetailsStruct{
				UserId:         order.UserId,
				ExecutionPrice: executionPrice,
				Market:         market,
				Side:           "-1",
			}

			makerExecutionDetails := &tarantool_pkg.ExecutionDetailsStruct{
				UserId:         makerOrder.UserId,
				ExecutionPrice: executionPrice,
				Market:         market,
				Side:           "1",
			}

			userMatchPositionDetails := &tarantool_pkg.PositionStruct{
				UserID:   order.UserId,
				Market:   market,
				AvgPrice: executionPrice,
				Side:     "-1",
			}

			makerMatchPositionDetails := &tarantool_pkg.PositionStruct{
				UserID:   makerOrder.UserId,
				Market:   market,
				AvgPrice: executionPrice,
				Side:     "1",
			}

			if order.PositionSize > makerOrder.PositionSize {
				executionPositionSize := makerOrder.PositionSize

				userExecutionDetails.ExecutionPositionSize = executionPositionSize
				checkUserMargin := c.checkAccountMargin(userExecutionDetails)
				if !checkUserMargin {
					c.logger.Info("user margin less than 10%", "orderUser", order.UserId)
					c.deleteOrderFromOrderBook(order)
					return false
				}

				makerExecutionDetails.ExecutionPositionSize = executionPositionSize
				checkMakerMargin := c.checkAccountMargin(makerExecutionDetails)
				if !checkMakerMargin {
					c.logger.Info("user margin less than 10%", "maker", makerOrder.UserId)
					c.deleteOrderFromOrderBook(makerOrder)
					continue
				}

				order.PositionSize = order.PositionSize - executionPositionSize

				// Insert Match Position for user and maker
				userMatchPositionDetails.PositionSize = executionPositionSize
				c.tarantool.InsertMatchedPosition(userMatchPositionDetails)
				makerMatchPositionDetails.PositionSize = executionPositionSize
				c.tarantool.InsertMatchedPosition(makerMatchPositionDetails)

				// Delete Maker Order from orderbook
				c.tarantool.DeleteOrderByPrimaryKey(makerOrder.UserId, makerOrder.Price, makerOrder.Side, makerOrder.Market)

				// Update Market Price
				c.tarantool.UpdateMarketPrice(market, executionPrice)
				c.logger.Info("matchingEngineForShortOrder", "executionPositionSize", executionPositionSize, "executionPrice", executionPrice, "orderUser", order.UserId, "maker", makerOrder.UserId)

			} else {
				executionPositionSize := order.PositionSize

				userExecutionDetails.ExecutionPositionSize = executionPositionSize
				checkUserMargin := c.checkAccountMargin(userExecutionDetails)
				if !checkUserMargin {
					c.logger.Info("user margin less than 10%", "orderUser", order.UserId)
					c.deleteOrderFromOrderBook(order)
					return false
				}

				makerExecutionDetails.ExecutionPositionSize = executionPositionSize
				checkMakerMargin := c.checkAccountMargin(makerExecutionDetails)
				if !checkMakerMargin {
					c.logger.Info("user margin less than 10%", "maker", makerOrder.UserId)
					c.deleteOrderFromOrderBook(makerOrder)
					continue
				}

				makerOrder.PositionSize = makerOrder.PositionSize - executionPositionSize
				order.PositionSize = 0

				// Insert Match Position for user and maker
				userMatchPositionDetails.PositionSize = executionPositionSize
				c.tarantool.InsertMatchedPosition(userMatchPositionDetails)
				makerMatchPositionDetails.PositionSize = executionPositionSize
				c.tarantool.InsertMatchedPosition(makerMatchPositionDetails)

				// Update Maker Order from orderbook
				if (makerOrder.PositionSize) == 0 {
					c.tarantool.DeleteOrderByPrimaryKey(makerOrder.UserId, makerOrder.Price, makerOrder.Side, makerOrder.Market)
				} else {
					c.tarantool.UpdateOrderByPrimaryKey(makerOrder.UserId, makerOrder.Price, makerOrder.Side, makerOrder.Market, makerOrder.PositionSize)
				}

				// Update Market Price
				c.tarantool.UpdateMarketPrice(market, executionPrice)
				c.logger.Info("matchingEngineForShortOrder", "executionPositionSize", executionPositionSize, "executionPrice", executionPrice, "orderUser", order.UserId, "maker", makerOrder.UserId)
			}

		}
	}
	return true
}

func (c *MatchingEngine) checkAccountMargin(executionDetail *tarantool_pkg.ExecutionDetailsStruct) bool {

	userWalletBalance := c.tarantool.GetUserWalletBalance(executionDetail.UserId)

	userPositions, err := c.tarantool.GetUserPositions(executionDetail.UserId)

	if err != nil {
		c.logger.Error("failed to get user positions", err)
		return false
	}

	userPositions = append(userPositions, &tarantool_pkg.PositionStruct{
		Market:       executionDetail.Market,
		PositionSize: executionDetail.ExecutionPositionSize,
		AvgPrice:     executionDetail.ExecutionPrice,
		Side:         executionDetail.Side,
	})

	accountEquity := c.tarantool.CalculateAccountEquity(userWalletBalance, userPositions)
	totalAccountNotional := c.tarantool.CalculateTotalAccountNotional(userPositions)

	accountMargin := c.tarantool.CalculateAccountMargin(accountEquity, totalAccountNotional)

	return accountMargin >= 0.1
}

func (c *MatchingEngine) sortOrderBook(orderbooks []*tarantool_pkg.OrderStruct) [][]*tarantool_pkg.OrderStruct {
	priceMap := map[float64][]*tarantool_pkg.OrderStruct{}
	for _, order := range orderbooks {
		priceMap[order.Price] = append(priceMap[order.Price], order)
	}

	// Extract the keys and sort them
	var prices []float64
	for price := range priceMap {
		prices = append(prices, price)
	}

	sort.Float64s(prices)

	// Create a sorted list of orders
	var sortedOrders [][]*tarantool_pkg.OrderStruct
	for _, price := range prices {
		orders := priceMap[price]
		// Sort orders by createdTime
		sort.SliceStable(orders, func(i, j int) bool {
			return orders[i].CreatedTime < orders[j].CreatedTime
		})
		sortedOrders = append(sortedOrders, orders)
	}
	return sortedOrders
}

func (c *MatchingEngine) deleteOrderFromOrderBook(makerOrder *tarantool_pkg.OrderStruct) {
	userOrderBook := c.tarantool.GetAllOrders()
	for _, order := range userOrderBook {
		if order.UserId == makerOrder.UserId {
			c.tarantool.DeleteOrderByPrimaryKey(order.UserId, order.Price, order.Side, order.Market)
		}
	}
}

// func (c *MatchingEngine) calculateAccountMargin(accountEquity float64, totalAccountNotional float64) float64 {
// 	if totalAccountNotional == 0 {
// 		if accountEquity > 0 {
// 			return float64(^uint(0) >> 1) // Return the maximum float64 value (represents high margin)
// 		}
// 		return 0 //if both zero
// 	}
// 	return accountEquity / totalAccountNotional
// }

// func (c *MatchingEngine) CalculateAccountEquity(walletBalance float64, positions []*tarantool_pkg.PositionStruct) float64 {
// 	equity := walletBalance
// 	for _, position := range positions {
// 		marketPrice := c.tarantool.GetMarketPriceByMarket(position.Market)
// 		if position.Side == "1" {
// 			equity += position.PositionSize * (marketPrice - position.AvgPrice)
// 		} else {
// 			equity += position.PositionSize * (position.AvgPrice - marketPrice)
// 		}
// 	}
// 	return equity
// }

// func (c *MatchingEngine) calculateTotalAccountNotional(positions []*tarantool_pkg.PositionStruct) float64 {
// 	totalNotional := 0.0
// 	for _, position := range positions {
// 		totalNotional += position.PositionSize * c.tarantool.GetMarketPriceByMarket(position.Market)
// 	}
// 	return totalNotional
// }
