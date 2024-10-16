package tarantool_pkg

import (
	"fmt"

	"github.com/tarantool/go-tarantool/v2"
)

type PositionStruct struct {
	UserID       string  `json:"user_id"`
	Market       string  `json:"market"`
	PositionSize float64 `json:"position_size"`
	AvgPrice     float64 `json:"avg_price"`
	Side         string  `json:"side"`
}

type TarantoolPositionConnInterface interface {
	InsertPosition(position *PositionStruct) error
	DeletePosition(userID string, market string, side string) error
	InsertMatchedPosition(position *PositionStruct) error
	GetAllPositions() ([]*PositionStruct, error)
	GetUserPositions(userID string) ([]*PositionStruct, error)
	GetNetPositionSizeByValidatingPosition(order *OrderStruct) (*OrderStruct, error)
}

const (
	positionSpace   = "positions"
	userMarketIndex = "user_market_side_index"
)

func (c *TarantoolClient) InsertMatchedPosition(position *PositionStruct) error {
	currentPosition, err := c.getPosition(position.UserID, position.Market, position.Side)
	if err != nil {
		c.logger.Error("failed to get current position", err)
		return err
	}

	// Update the average price and position size
	if currentPosition != nil {
		newPositionSize := currentPosition.PositionSize + position.PositionSize
		newAvgPrice := (currentPosition.AvgPrice*currentPosition.PositionSize + position.AvgPrice*position.PositionSize) / newPositionSize
		position.PositionSize = newPositionSize
		position.AvgPrice = newAvgPrice

	}

	return c.InsertPosition(position)

}

func (c *TarantoolClient) InsertPosition(position *PositionStruct) error {
	conn := c.conn
	_, err := conn.Do(
		tarantool.NewUpsertRequest(positionSpace).
			Tuple([]interface{}{
				position.UserID,
				position.Market,
				position.PositionSize,
				position.AvgPrice,
				position.Side,
			}).
			Operations(tarantool.NewOperations().
				Assign(2, position.PositionSize).
				Assign(3, position.AvgPrice),
			),
	).Get()
	if err != nil {
		c.logger.Error("failed to get current position", err)
		return err
	}
	return nil
}

func (c *TarantoolClient) getPosition(userID string, market string, side string) (*PositionStruct, error) {
	conn := c.conn
	result, err := conn.Do(
		tarantool.NewSelectRequest(positionSpace).
			Index(userMarketIndex).
			Iterator(tarantool.IterEq).
			Key([]interface{}{userID, market, side}),
	).Get()

	if err != nil {
		return nil, fmt.Errorf("failed to get position: %w", err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	positions := c.getTransformPositionList(result)
	var position *PositionStruct
	if len(positions) > 0 {
		position = positions[0]
	}

	return position, nil
}

func (c *TarantoolClient) getTransformPositionList(data []interface{}) []*PositionStruct {
	positions := []*PositionStruct{}
	for _, item := range data {
		data, ok := item.([]interface{})
		if !ok {
			return nil
		}
		positions = append(positions, c.transformToPositionStruct(data))
	}
	return positions
}

func (c *TarantoolClient) getTransformPositionListByUserId(data []interface{}, userId string) []*PositionStruct {
	positions := []*PositionStruct{}
	for _, item := range data {
		data, ok := item.([]interface{})
		if !ok {
			return nil
		}

		position := c.transformToPositionStruct(data)
		if position.UserID == userId {
			positions = append(positions, c.transformToPositionStruct(data))
		}
	}
	return positions
}

func (c *TarantoolClient) transformToPositionStruct(data []interface{}) *PositionStruct {
	return &PositionStruct{
		UserID:       data[0].(string),
		Market:       data[1].(string),
		PositionSize: c.convertToFloat64(data[2]),
		AvgPrice:     c.convertToFloat64(data[3]),
		Side:         data[4].(string),
	}
}

// get all positions
func (c *TarantoolClient) GetAllPositions() ([]*PositionStruct, error) {
	conn := c.conn
	result, err := conn.Do(
		tarantool.NewSelectRequest(positionSpace).
			Iterator(tarantool.IterAll).
			Key([]interface{}{}),
	).Get()

	if err != nil {
		return nil, fmt.Errorf("failed to get all positions: %w", err)
	}

	positions := c.getTransformPositionList(result)

	return positions, nil
}

func (c *TarantoolClient) DeletePosition(userID string, market string, side string) error {
	conn := c.conn
	_, err := conn.Do(
		tarantool.NewDeleteRequest(positionSpace).
			Index(userMarketIndex).
			Key([]interface{}{userID, market, side}),
	).Get()
	if err != nil {
		return fmt.Errorf("failed to delete position: %w", err)
	}
	return nil
}

func (c *TarantoolClient) GetUserPositions(userID string) ([]*PositionStruct, error) {
	conn := c.conn
	result, err := conn.Do(
		tarantool.NewSelectRequest(positionSpace).
			Iterator(tarantool.IterAll).
			Key([]interface{}{}),
	).Get()

	if err != nil {
		return nil, fmt.Errorf("failed to get user positions: %w", err)
	}

	positions := c.getTransformPositionListByUserId(result, userID)

	return positions, nil
}

func (c *TarantoolClient) GetNetPositionSizeByValidatingPosition(order *OrderStruct) (*OrderStruct, error) {
	// Get all the positions for the user
	userId := order.UserId
	positions, err := c.GetUserPositions(userId)

	if err != nil {
		c.logger.Error("failed to get retrieve existings positions", err)
		return order, err
	}

	oppositeSide := "1"

	if order.Side == oppositeSide {
		oppositeSide = "-1"
	}

	newPositionSize := 0.0
	netPositionSize := 0.0
	isOppositePositionFound := false
	for _, position := range positions {
		if position.Market == order.Market && position.Side == oppositeSide {
			userWalletBalance := c.GetUserWalletBalance(userId)
			isOppositePositionFound = true
			// Update the order
			netPositionSize = position.PositionSize - order.PositionSize
			if netPositionSize <= 0 {
				newPositionSize = -netPositionSize

				// Close the current position as it is opposite direction
				if err := c.DeletePosition(userId, position.Market, position.Side); err != nil {
					c.logger.Error("failed to delete position", err)
					return order, err
				}

				//Refund the amount to user
				// Calculate the profit and loss (PnL)
				marketPrice := c.GetMarketPriceByMarket(position.Market)

				// If is long position, the side factor is 1, else -1
				sideFactor := 1.0
				if position.Side == "-1" {
					sideFactor = -1.0
				}

				pnl := sideFactor * (marketPrice - position.AvgPrice) * position.PositionSize

				// Update pnl to user balance as position closed
				c.UpdateUserWalletBalance(userId, userWalletBalance+pnl)

			} else {
				newPositionSize = 0
				// Reduce the current position with order size assign newPosition = 0
				// New position will be the net position and need to deduct and refund the amount to user based on current order size
				// Update the current position
				if err := c.InsertPosition(&PositionStruct{
					UserID:       userId,
					Market:       position.Market,
					PositionSize: netPositionSize,
					AvgPrice:     position.AvgPrice,
					Side:         position.Side,
				}); err != nil {
					c.logger.Error("failed to insert position", err)
					return order, err
				}

				//Calculate the profit and lost (PnL)
				marketPrice := c.GetMarketPriceByMarket(position.Market)

				sideFactor := 1.0
				if position.Side == "-1" {
					sideFactor = -1.0
				}

				pnl := sideFactor * (marketPrice - position.AvgPrice) * order.PositionSize

				// Update pnl to user balance as position closed
				c.UpdateUserWalletBalance(userId, userWalletBalance+pnl)
			}
		}
	}

	if isOppositePositionFound {
		order.PositionSize = newPositionSize
	}

	// If order position size is 0, it will not trigger the matching engine.
	return order, nil
}
