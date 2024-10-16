package matching_engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wen-git-acc/orderbook/pkg/tarantool_pkg"
)

type MockTarantoolClient struct {
	mock.Mock
	*tarantool_pkg.TarantoolClient
}

func (m *MockTarantoolClient) sortOrderBook(orderbooks []*tarantool_pkg.OrderStruct) [][]*tarantool_pkg.OrderStruct {
	args := m.Called(orderbooks)
	return args.Get(0).([][]*tarantool_pkg.OrderStruct)
}

func (m *MockTarantoolClient) InsertNewOrder(order *tarantool_pkg.OrderStruct) {
	m.Called(order)
}

func (m *MockTarantoolClient) CheckAccountMargin(executionDetail *tarantool_pkg.ExecutionDetailsStruct) bool {
	args := m.Called(executionDetail)
	return args.Bool(0)
}

func (m *MockTarantoolClient) InsertMatchedPosition(position *tarantool_pkg.PositionStruct) error {
	args := m.Called(position)
	return args.Error(0)
}

func (m *MockTarantoolClient) DeleteOrderByPrimaryKey(userId string, price float64, side string, market string) error {
	args := m.Called(userId, price, side, market)
	return args.Error(0)
}

func (m *MockTarantoolClient) UpdateMarketPrice(market string, price float64) {
	m.Called(market, price)
}

func (m *MockTarantoolClient) UpdateOrderByPrimaryKey(userId string, price float64, side string, market string, positionSize int) error {
	args := m.Called(userId, price, side, market, positionSize)
	return args.Error(0)
}

func (m *MockTarantoolClient) GetUserWalletBalance(userId string) float64 {
	args := m.Called(userId)
	return args.Get(0).(float64)
}

func (m *MockTarantoolClient) GetUserPositions(userId string) ([]*tarantool_pkg.PositionStruct, error) {
	args := m.Called(userId)
	return args.Get(0).([]*tarantool_pkg.PositionStruct), args.Error(1)
}

func TestMatchingEngineForLongOrder(t *testing.T) {
	// Create a mock Tarantool client
	mockClient := new(MockTarantoolClient)

	mockLongOrder := &tarantool_pkg.OrderStruct{
		Price:        100,
		UserId:       "user1",
		PositionSize: 10,
		Market:       "BTC-USD",
		Side:         "1",
	}

	fullMatchShortOrder := &tarantool_pkg.OrderStruct{
		Price:        100,
		UserId:       "user2",
		PositionSize: 10,
		Market:       "BTC-USD",
		Side:         "-1",
	}

	partialMatchShortOrder := &tarantool_pkg.OrderStruct{
		Price:        100,
		UserId:       "user2",
		PositionSize: 5,
		Market:       "BTC-USD",
		Side:         "-1",
	}

	executionDetailMock := &tarantool_pkg.ExecutionDetailsStruct{
		UserId:                "user1",
		Market:                "BTC-USD",
		Side:                  "1",
		ExecutionPositionSize: 10,
		ExecutionPrice:        100,
	}

	// Define test cases
	testCases := []struct {
		name            string
		order           *tarantool_pkg.OrderStruct
		orderBook       []*tarantool_pkg.OrderStruct
		sortedOrderBook [][]*tarantool_pkg.OrderStruct
		expectedResult  bool
	}{
		{
			name:            "Empty order book",
			order:           mockLongOrder,
			orderBook:       []*tarantool_pkg.OrderStruct{},
			sortedOrderBook: [][]*tarantool_pkg.OrderStruct{},
			expectedResult:  false,
		},
		{
			name:            "Partial match",
			order:           mockLongOrder,
			orderBook:       []*tarantool_pkg.OrderStruct{partialMatchShortOrder},
			sortedOrderBook: [][]*tarantool_pkg.OrderStruct{{partialMatchShortOrder}},
			expectedResult:  true,
		},
		{
			name:            "Full match",
			order:           fullMatchShortOrder,
			orderBook:       []*tarantool_pkg.OrderStruct{fullMatchShortOrder},
			sortedOrderBook: [][]*tarantool_pkg.OrderStruct{{fullMatchShortOrder}},
			expectedResult:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup the mock expectations
			mockClient.On("CheckAccountMargin", executionDetailMock).Return(true)
			mockClient.On("sortOrderBook", tc.orderBook).Return(tc.sortedOrderBook)
			mockClient.On("InsertNewOrder", tc.order).Return()
			mockClient.On("GetUserWalletBalance", mock.Anything).Return(12.0)
			mockClient.On("InsertMatchedPosition", mock.Anything).Return(nil)
			mockClient.On("DeleteOrderByPrimaryKey", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			mockClient.On("UpdateMarketPrice", mock.Anything, mock.Anything).Return()
			mockClient.On("deleteOrderFromOrderBook", mock.Anything).Return()

			//Initialize matching engine
			matchingEngine := &MatchingEngine{
				tarantool: mockClient,
			}
			// Call the function under test
			result := mockClient.MatchingEngineForLongOrder(tc.order, tc.orderBook)

			// Assert the results
			assert.Equal(t, tc.expectedResult, result)
			mockClient.AssertExpectations(t)
		})
	}
}
