package tarantool_pkg

import (
	"fmt"
	"strconv"
)

func (c *TarantoolClient) convertToFloat64(data interface{}) float64 {
	if balance, ok := data.(float64); ok {
		return balance
	}
	strnum := fmt.Sprintf("%v", data)
	floatNum, err := strconv.ParseFloat(strnum, 64)

	if err != nil {
		return 0
	}
	return floatNum
}
