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

func (c *TarantoolClient) convertToInt64(str string) int64 {
	intValue, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Println("Error converting string to int64:", err)
		return 0
	}
	return intValue
}

func (c *TarantoolClient) convertToInt(data interface{}) int {
	if intValue, ok := data.(int); ok {
		return intValue
	}
	strnum := fmt.Sprintf("%v", data)
	intNum, err := strconv.Atoi(strnum)

	if err != nil {
		return 0
	}
	return intNum
}
