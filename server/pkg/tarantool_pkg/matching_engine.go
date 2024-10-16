package tarantool_pkg

import "sort"

// need to put matching, and get bid and get ask orderbook to expose
// Pass in ask orderbook
func (c *TarantoolClient) matchingEngineForLongOrder(order *OrderStruct, orderBook []*OrderStruct) bool {
	market := order.Market

	if len(orderBook) == 0 {
		return false
	}

	sortedOrderBook := c.SortOrderBook(orderBook)
	if len(sortedOrderBook) > 0 {
		if sortedOrderBook[0][0].Price > order.Price {
			c.InsertNewOrder(order)
			c.updateUserWalletAmountWithDeduction(order.UserId, order.PositionSize*order.Price)
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
				checkUserMargin := c.checkAccountMargin(userExecutionDetails, false)
				if !checkUserMargin {
					c.logger.Info("user margin less than 10%", "orderUser", order.UserId)
					return false
				}

				makerExecutionDetails.ExecutionPositionSize = executionPositionSize
				checkMakerMargin := c.checkAccountMargin(makerExecutionDetails, true)
				if !checkMakerMargin {
					c.logger.Info("user margin less than 10%", "maker", makerOrder.UserId)
					c.deleteMakerOrderFromOrderBook(makerOrder)
					continue
				}

				order.PositionSize = order.PositionSize - executionPositionSize

				// Insert Match Position for user and maker
				userMatchPositionDetails.PositionSize = executionPositionSize
				c.InsertMatchedPosition(userMatchPositionDetails)
				c.updateUserWalletAmountWithDeduction(order.UserId, executionPositionSize*executionPrice)
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
				checkUserMargin := c.checkAccountMargin(userExecutionDetails, false)
				if !checkUserMargin {
					c.logger.Info("user margin less than 10%", "orderUser", order.UserId)
					return false
				}

				makerExecutionDetails.ExecutionPositionSize = executionPositionSize
				checkMakerMargin := c.checkAccountMargin(makerExecutionDetails, true)
				if !checkMakerMargin {
					c.logger.Info("user margin less than 10%", "maker", makerOrder.UserId)
					c.deleteMakerOrderFromOrderBook(makerOrder)
					continue
				}

				makerOrder.PositionSize = makerOrder.PositionSize - executionPositionSize
				order.PositionSize = 0

				// Insert Match Position for user and maker
				userMatchPositionDetails.PositionSize = executionPositionSize
				c.InsertMatchedPosition(userMatchPositionDetails)
				c.updateUserWalletAmountWithDeduction(order.UserId, executionPositionSize*executionPrice)
				makerMatchPositionDetails.PositionSize = executionPositionSize
				c.InsertMatchedPosition(makerMatchPositionDetails)

				// Update Maker Order from orderbook
				if (makerOrder.PositionSize) == 0 {
					c.DeleteOrderByPrimaryKey(makerOrder.UserId, makerOrder.Price, makerOrder.Side, makerOrder.Market)
				} else {
					c.updateOrderByPrimaryKey(makerOrder.UserId, makerOrder.Price, makerOrder.Side, makerOrder.Market, makerOrder.PositionSize)
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
func (c *TarantoolClient) matchingEngineForShortOrder(order *OrderStruct, orderBook []*OrderStruct) bool {
	market := order.Market

	if len(orderBook) == 0 {
		return false
	}

	sortedOrderBook := c.SortOrderBook(orderBook)
	if len(sortedOrderBook) > 0 {
		if sortedOrderBook[len(sortedOrderBook)-1][0].Price < order.Price {
			c.InsertNewOrder(order)
			c.updateUserWalletAmountWithDeduction(order.UserId, order.PositionSize*order.Price)
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
				checkUserMargin := c.checkAccountMargin(userExecutionDetails, false)
				if !checkUserMargin {
					c.logger.Info("user margin less than 10%", "orderUser", order.UserId)
					return false
				}

				makerExecutionDetails.ExecutionPositionSize = executionPositionSize
				checkMakerMargin := c.checkAccountMargin(makerExecutionDetails, true)
				if !checkMakerMargin {
					c.logger.Info("user margin less than 10%", "maker", makerOrder.UserId)
					c.deleteMakerOrderFromOrderBook(makerOrder)
					continue
				}

				order.PositionSize = order.PositionSize - executionPositionSize

				// Insert Match Position for user and maker
				userMatchPositionDetails.PositionSize = executionPositionSize
				c.InsertMatchedPosition(userMatchPositionDetails)
				c.updateUserWalletAmountWithDeduction(order.UserId, executionPositionSize*executionPrice)
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
				checkUserMargin := c.checkAccountMargin(userExecutionDetails, false)
				if !checkUserMargin {
					c.logger.Info("user margin less than 10%", "orderUser", order.UserId)
					return false
				}

				makerExecutionDetails.ExecutionPositionSize = executionPositionSize
				checkMakerMargin := c.checkAccountMargin(makerExecutionDetails, true)
				if !checkMakerMargin {
					c.logger.Info("user margin less than 10%", "maker", makerOrder.UserId)
					c.deleteMakerOrderFromOrderBook(makerOrder)
					continue
				}

				makerOrder.PositionSize = makerOrder.PositionSize - executionPositionSize
				order.PositionSize = 0

				// Insert Match Position for user and maker
				userMatchPositionDetails.PositionSize = executionPositionSize
				c.InsertMatchedPosition(userMatchPositionDetails)
				c.updateUserWalletAmountWithDeduction(order.UserId, executionPositionSize*executionPrice)
				makerMatchPositionDetails.PositionSize = executionPositionSize
				c.InsertMatchedPosition(makerMatchPositionDetails)

				// Update Maker Order from orderbook
				if (makerOrder.PositionSize) == 0 {
					c.DeleteOrderByPrimaryKey(makerOrder.UserId, makerOrder.Price, makerOrder.Side, makerOrder.Market)
				} else {
					c.updateOrderByPrimaryKey(makerOrder.UserId, makerOrder.Price, makerOrder.Side, makerOrder.Market, makerOrder.PositionSize)
				}

				// Update Market Price
				c.UpdateMarketPrice(market, executionPrice)
				c.logger.Info("matchingEngineForShortOrder", "executionPositionSize", executionPositionSize, "executionPrice", executionPrice, "orderUser", order.UserId, "maker", makerOrder.UserId)
			}

		}
	}
	return true
}

func (c *TarantoolClient) updateUserWalletAmountWithDeduction(userId string, amountToDeduct float64) {
	balance := c.GetUserWalletBalance(userId)
	balance = balance - amountToDeduct
	c.UpdateUserWalletBalance(userId, balance)
}

func (c *TarantoolClient) checkAccountMargin(executionDetail *ExecutionDetailsStruct, isMaker bool) bool {

	userWalletBalance := c.GetUserWalletBalance(executionDetail.UserId)

	if !isMaker && (userWalletBalance == 0 || (userWalletBalance < executionDetail.ExecutionPositionSize*executionDetail.ExecutionPrice)) {
		// User has insufficient balance
		// Maker wallet balance is already deducted when writing into orderbook.
		return false
	}

	// Account for wallet balance after deduct.
	if !isMaker {
		// Money is deducted when it first enters the orderbook
		userWalletBalance = userWalletBalance - executionDetail.ExecutionPositionSize*executionDetail.ExecutionPrice
	}

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

	accountEquity := c.calculateAccountEquity(userWalletBalance, userPositions)
	totalAccountNotional := c.calculateTotalAccountNotional(userPositions)

	accountMargin := c.calculateAccountMargin(accountEquity, totalAccountNotional)

	return accountMargin >= 0.1
}

func (c *TarantoolClient) SortOrderBook(orderbooks []*OrderStruct) [][]*OrderStruct {
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
			return orders[i].createdTime < orders[j].createdTime
		})
		sortedOrders = append(sortedOrders, orders)
	}
	return sortedOrders
}

func (c *TarantoolClient) getBidOrderBook(currentOrder *OrderStruct) []*OrderStruct {
	return c.getOrdersByMarketAndSide(currentOrder.Market, "1")
}

func (c *TarantoolClient) getAskOrderBook(currentOrder *OrderStruct) []*OrderStruct {
	return c.getOrdersByMarketAndSide(currentOrder.Market, "-1")
}

func (c *TarantoolClient) deleteMakerOrderFromOrderBook(makerOrder *OrderStruct) {
	userOrderBook := c.getAllOrders()
	for _, order := range userOrderBook {
		if order.UserId == makerOrder.UserId {
			refundAmount := order.PositionSize * order.Price
			walletBalance := c.GetUserWalletBalance(order.UserId)
			c.UpdateUserWalletBalance(order.UserId, walletBalance+refundAmount)
			c.DeleteOrderByPrimaryKey(order.UserId, order.Price, order.Side, order.Market)
		}
	}
}
