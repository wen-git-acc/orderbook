package tarantool_pkg

import "sort"

type MatchingEngineInterface interface {
	CalculateAccountEquity(walletBalance float64, positions []*PositionStruct) float64
}

// Pass in ask orderbook
func (c *TarantoolClient) MatchingEngineForLongOrder(order *OrderStruct, orderBook []*OrderStruct) bool {
	market := order.Market

	if len(orderBook) == 0 {
		return false
	}

	sortedOrderBook := c.sortOrderBook(orderBook)
	if len(sortedOrderBook) > 0 {
		if sortedOrderBook[0][0].Price > order.Price {
			isMarginSufficient := c.checkAccountMargin(&ExecutionDetailsStruct{
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
			c.InsertNewOrder(order)
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

			userExecutionDetails := &ExecutionDetailsStruct{
				UserId:         order.UserId,
				ExecutionPrice: executionPrice,
				Market:         market,
				Side:           "1",
			}

			makerExecutionDetails := &ExecutionDetailsStruct{
				UserId:         makerOrder.UserId,
				ExecutionPrice: executionPrice,
				Market:         market,
				Side:           "-1",
			}

			userMatchPositionDetails := &PositionStruct{
				UserID:   order.UserId,
				Market:   market,
				AvgPrice: executionPrice,
				Side:     "1",
			}

			makerMatchPositionDetails := &PositionStruct{
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
				c.InsertMatchedPosition(userMatchPositionDetails)
				makerMatchPositionDetails.PositionSize = executionPositionSize
				c.InsertMatchedPosition(makerMatchPositionDetails)

				// Delete Maker Order from orderbook
				c.DeleteOrderByPrimaryKey(makerOrder.UserId, makerOrder.Price, makerOrder.Side, makerOrder.Market)

				// Update Market Price
				c.UpdateMarketPrice(market, executionPrice)
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
				c.InsertMatchedPosition(userMatchPositionDetails)
				makerMatchPositionDetails.PositionSize = executionPositionSize
				c.InsertMatchedPosition(makerMatchPositionDetails)

				// Update Maker Order from orderbook
				if (makerOrder.PositionSize) == 0 {
					c.DeleteOrderByPrimaryKey(makerOrder.UserId, makerOrder.Price, makerOrder.Side, makerOrder.Market)
				} else {
					c.UpdateOrderByPrimaryKey(makerOrder.UserId, makerOrder.Price, makerOrder.Side, makerOrder.Market, makerOrder.PositionSize)
				}
				// Update Market Price
				c.UpdateMarketPrice(market, executionPrice)
				c.logger.Info("matchingEngineForShortOrder", "executionPositionSize", executionPositionSize, "executionPrice", executionPrice, "orderUser", order.UserId, "maker", makerOrder.UserId)

			}

		}
	}
	return true
}

// Pass in bid orderbook
func (c *TarantoolClient) MatchingEngineForShortOrder(order *OrderStruct, orderBook []*OrderStruct) bool {
	market := order.Market

	if len(orderBook) == 0 {
		return false
	}

	sortedOrderBook := c.sortOrderBook(orderBook)
	if len(sortedOrderBook) > 0 {
		if sortedOrderBook[len(sortedOrderBook)-1][0].Price < order.Price {
			isMarginSufficient := c.checkAccountMargin(&ExecutionDetailsStruct{
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
			c.InsertNewOrder(order)
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

			userExecutionDetails := &ExecutionDetailsStruct{
				UserId:         order.UserId,
				ExecutionPrice: executionPrice,
				Market:         market,
				Side:           "-1",
			}

			makerExecutionDetails := &ExecutionDetailsStruct{
				UserId:         makerOrder.UserId,
				ExecutionPrice: executionPrice,
				Market:         market,
				Side:           "1",
			}

			userMatchPositionDetails := &PositionStruct{
				UserID:   order.UserId,
				Market:   market,
				AvgPrice: executionPrice,
				Side:     "-1",
			}

			makerMatchPositionDetails := &PositionStruct{
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
				c.InsertMatchedPosition(userMatchPositionDetails)
				makerMatchPositionDetails.PositionSize = executionPositionSize
				c.InsertMatchedPosition(makerMatchPositionDetails)

				// Delete Maker Order from orderbook
				c.DeleteOrderByPrimaryKey(makerOrder.UserId, makerOrder.Price, makerOrder.Side, makerOrder.Market)

				// Update Market Price
				c.UpdateMarketPrice(market, executionPrice)
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
				c.InsertMatchedPosition(userMatchPositionDetails)
				makerMatchPositionDetails.PositionSize = executionPositionSize
				c.InsertMatchedPosition(makerMatchPositionDetails)

				// Update Maker Order from orderbook
				if (makerOrder.PositionSize) == 0 {
					c.DeleteOrderByPrimaryKey(makerOrder.UserId, makerOrder.Price, makerOrder.Side, makerOrder.Market)
				} else {
					c.UpdateOrderByPrimaryKey(makerOrder.UserId, makerOrder.Price, makerOrder.Side, makerOrder.Market, makerOrder.PositionSize)
				}

				// Update Market Price
				c.UpdateMarketPrice(market, executionPrice)
				c.logger.Info("matchingEngineForShortOrder", "executionPositionSize", executionPositionSize, "executionPrice", executionPrice, "orderUser", order.UserId, "maker", makerOrder.UserId)
			}

		}
	}
	return true
}

func (c *TarantoolClient) checkAccountMargin(executionDetail *ExecutionDetailsStruct) bool {

	userWalletBalance := c.GetUserWalletBalance(executionDetail.UserId)

	userPositions, err := c.GetUserPositions(executionDetail.UserId)

	if err != nil {
		c.logger.Error("failed to get user positions", err)
		return false
	}

	userPositions = append(userPositions, &PositionStruct{
		Market:       executionDetail.Market,
		PositionSize: executionDetail.ExecutionPositionSize,
		AvgPrice:     executionDetail.ExecutionPrice,
		Side:         executionDetail.Side,
	})

	accountEquity := c.CalculateAccountEquity(userWalletBalance, userPositions)
	totalAccountNotional := c.calculateTotalAccountNotional(userPositions)

	accountMargin := c.calculateAccountMargin(accountEquity, totalAccountNotional)

	return accountMargin >= 0.1
}

func (c *TarantoolClient) sortOrderBook(orderbooks []*OrderStruct) [][]*OrderStruct {
	priceMap := map[float64][]*OrderStruct{}
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
	var sortedOrders [][]*OrderStruct
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

func (c *TarantoolClient) GetBidOrderBook(currentOrder *OrderStruct) []*OrderStruct {
	return c.GetOrdersByMarketAndSide(currentOrder.Market, "1")
}

func (c *TarantoolClient) GetAskOrderBook(currentOrder *OrderStruct) []*OrderStruct {
	return c.GetOrdersByMarketAndSide(currentOrder.Market, "-1")
}

func (c *TarantoolClient) deleteOrderFromOrderBook(makerOrder *OrderStruct) {
	userOrderBook := c.GetAllOrders()
	for _, order := range userOrderBook {
		if order.UserId == makerOrder.UserId {
			c.DeleteOrderByPrimaryKey(order.UserId, order.Price, order.Side, order.Market)
		}
	}
}

func (c *TarantoolClient) calculateAccountMargin(accountEquity float64, totalAccountNotional float64) float64 {
	if totalAccountNotional == 0 {
		if accountEquity > 0 {
			return float64(^uint(0) >> 1) // Return the maximum float64 value (represents high margin)
		}
		return 0 //if both zero
	}
	return accountEquity / totalAccountNotional
}

func (c *TarantoolClient) CalculateAccountEquity(walletBalance float64, positions []*PositionStruct) float64 {
	equity := walletBalance
	for _, position := range positions {
		marketPrice := c.GetMarketPriceByMarket(position.Market)
		if position.Side == "1" {
			equity += position.PositionSize * (marketPrice - position.AvgPrice)
		} else {
			equity += position.PositionSize * (position.AvgPrice - marketPrice)
		}
	}
	return equity
}

func (c *TarantoolClient) calculateTotalAccountNotional(positions []*PositionStruct) float64 {
	totalNotional := 0.0
	for _, position := range positions {
		totalNotional += position.PositionSize * c.GetMarketPriceByMarket(position.Market)
	}
	return totalNotional
}
