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
	InsertMatchedPosition(position *PositionStruct) error
	GetPosition(userID string, market string, side string) (*PositionStruct, error)
	GetAllPositions() ([]*PositionStruct, error)
	DeletePosition(userID string, market string, side string) error
}

const (
	positionSpace   = "positions"
	userMarketIndex = "user_market_side_index"
)

func (c *TarantoolClient) InsertMatchedPosition(position *PositionStruct) error {
	currentPosition, err := c.GetPosition(position.UserID, position.Market, position.Side)
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

// @github help me write a insert position function
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

// @github help me write a get position function
func (c *TarantoolClient) GetPosition(userID string, market string, side string) (*PositionStruct, error) {
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
