package matching_engine

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/wen-git-acc/orderbook/pkg/logger"
	"github.com/wen-git-acc/orderbook/pkg/tarantool_pkg"
)

type MockTarantoolClient struct {
	mock.Mock
}

func (m *MockTarantoolClient) CalculateAccountMargin(accountEquity float64, totalAccountNotional float64) float64 {
	args := m.Called(accountEquity, totalAccountNotional)
	return args.Get(0).(float64)
}

func (m *MockTarantoolClient) CalculateTotalAccountNotional(positions []*tarantool_pkg.PositionStruct) float64 {
	args := m.Called(positions)
	return args.Get(0).(float64)
}

func (m *MockTarantoolClient) CalculateAccountEquity(walletBalance float64, positions []*tarantool_pkg.PositionStruct) float64 {
	args := m.Called(walletBalance, positions)
	return args.Get(0).(float64)
}

func (m *MockTarantoolClient) GetMarketPriceByMarket(market string) float64 {
	args := m.Called(market)
	return args.Get(0).(float64)
}

func (m *MockTarantoolClient) UpdateMarketPrice(market string, marketPrice float64) {
	m.Called(market, marketPrice)
}

func (m *MockTarantoolClient) GetOrderBook(market string) *tarantool_pkg.SimplifiedOrderBook {
	args := m.Called(market)
	return args.Get(0).(*tarantool_pkg.SimplifiedOrderBook)
}

func (m *MockTarantoolClient) DeleteOrderByPrimaryKey(userId string, price float64, side string, market string) error {
	args := m.Called(userId, price, side, market)
	return args.Error(0)
}

func (m *MockTarantoolClient) GetOrderByPrimaryKey(userId string, price float64, side string, market string) *tarantool_pkg.OrderStruct {
	args := m.Called(userId, price, side, market)
	return args.Get(0).(*tarantool_pkg.OrderStruct)
}

func (m *MockTarantoolClient) UpdateOrderByPrimaryKey(userId string, price float64, side string, market string, positionSize float64) error {
	args := m.Called(userId, price, side, market, positionSize)
	return args.Error(0)
}

func (m *MockTarantoolClient) GetAllOrders() []*tarantool_pkg.OrderStruct {
	args := m.Called()
	return args.Get(0).([]*tarantool_pkg.OrderStruct)
}

func (m *MockTarantoolClient) GetOrdersByMarketAndSide(market string, side string) []*tarantool_pkg.OrderStruct {
	args := m.Called(market, side)
	return args.Get(0).([]*tarantool_pkg.OrderStruct)
}

func (m *MockTarantoolClient) GetBidOrderBook(currentOrder *tarantool_pkg.OrderStruct) []*tarantool_pkg.OrderStruct {
	args := m.Called(currentOrder)
	return args.Get(0).([]*tarantool_pkg.OrderStruct)
}

func (m *MockTarantoolClient) GetAskOrderBook(currentOrder *tarantool_pkg.OrderStruct) []*tarantool_pkg.OrderStruct {
	args := m.Called(currentOrder)
	return args.Get(0).([]*tarantool_pkg.OrderStruct)
}

func (m *MockTarantoolClient) InsertNewOrder(order *tarantool_pkg.OrderStruct) error {
	m.Called(order)
	return nil
}

func (m *MockTarantoolClient) InsertPosition(position *tarantool_pkg.PositionStruct) error {
	args := m.Called(position)
	return args.Error(0)
}

func (m *MockTarantoolClient) DeletePosition(userID string, market string, side string) error {
	args := m.Called(userID, market, side)
	return args.Error(0)
}

func (m *MockTarantoolClient) InsertMatchedPosition(position *tarantool_pkg.PositionStruct) error {
	args := m.Called(position)
	return args.Error(0)
}

func (m *MockTarantoolClient) GetAllPositions() ([]*tarantool_pkg.PositionStruct, error) {
	args := m.Called()
	return args.Get(0).([]*tarantool_pkg.PositionStruct), args.Error(1)
}

func (m *MockTarantoolClient) GetUserPositions(userID string) ([]*tarantool_pkg.PositionStruct, error) {
	args := m.Called(userID)
	return args.Get(0).([]*tarantool_pkg.PositionStruct), args.Error(1)
}

func (m *MockTarantoolClient) GetNetPositionSizeByValidatingPosition(order *tarantool_pkg.OrderStruct) (*tarantool_pkg.OrderStruct, error) {
	args := m.Called(order)
	return args.Get(0).(*tarantool_pkg.OrderStruct), args.Error(1)
}

func (m *MockTarantoolClient) IsUserRegistered(userID string) bool {
	args := m.Called(userID)
	return args.Bool(0)
}

func (m *MockTarantoolClient) GetUserWalletBalance(userID string) float64 {
	args := m.Called(userID)
	return args.Get(0).(float64)
}

func (m *MockTarantoolClient) UpdateUserWalletBalance(userID string, balance float64) error {
	args := m.Called(userID, balance)
	return args.Error(0)
}

func (m *MockTarantoolClient) CreateUserWalletBalance(userID string, balance float64) error {
	args := m.Called(userID, balance)
	return args.Error(0)
}

func TestMatchingEngineForLongOrder(t *testing.T) {
	// Create a mock Tarantool client
	mockClient := new(MockTarantoolClient)

	positionUser1Mock := []*tarantool_pkg.PositionStruct{
		{
			UserID:       "user1",
			Market:       "BTC-USD",
			PositionSize: 10,
			AvgPrice:     100,
			Side:         "1",
		},
		{
			UserID:       "user1",
			Market:       "BTC-USD",
			PositionSize: 10,
			AvgPrice:     100,
			Side:         "1",
		},
	}

	// Define test cases
	testCases := []struct {
		name            string
		order           *tarantool_pkg.OrderStruct
		orderBook       []*tarantool_pkg.OrderStruct
		sortedOrderBook [][]*tarantool_pkg.OrderStruct
		usersPosition   []*tarantool_pkg.PositionStruct
		expectedResult  bool
	}{
		{
			name: "Empty order book",
			order: &tarantool_pkg.OrderStruct{
				Price:        100,
				UserId:       "user1",
				PositionSize: 10,
				Market:       "BTC-USD",
				Side:         "1",
			},
			orderBook:       []*tarantool_pkg.OrderStruct{},
			sortedOrderBook: [][]*tarantool_pkg.OrderStruct{},
			expectedResult:  true,
			usersPosition:   positionUser1Mock,
		},
		{
			name: "Partial match",
			order: &tarantool_pkg.OrderStruct{
				Price:        100,
				UserId:       "user1",
				PositionSize: 10,
				Market:       "BTC-USD",
				Side:         "1",
			},
			orderBook: []*tarantool_pkg.OrderStruct{
				{
					Price:        100,
					UserId:       "user2",
					PositionSize: 5,
					Market:       "BTC-USD",
					Side:         "-1",
				},
			},
			sortedOrderBook: [][]*tarantool_pkg.OrderStruct{{
				{
					Price:        100,
					UserId:       "user2",
					PositionSize: 5,
					Market:       "BTC-USD",
					Side:         "-1",
				},
			}},
			expectedResult: true,
			usersPosition:  positionUser1Mock,
		},
		{
			name: "Full match",
			order: &tarantool_pkg.OrderStruct{
				Price:        100,
				UserId:       "user1",
				PositionSize: 10,
				Market:       "BTC-USD",
				Side:         "1",
			},
			orderBook: []*tarantool_pkg.OrderStruct{
				{
					Price:        100,
					UserId:       "user2",
					PositionSize: 10,
					Market:       "BTC-USD",
					Side:         "-1",
				},
			},
			sortedOrderBook: [][]*tarantool_pkg.OrderStruct{{
				{
					Price:        100,
					UserId:       "user2",
					PositionSize: 10,
					Market:       "BTC-USD",
					Side:         "-1",
				},
			}},
			expectedResult: true,
			usersPosition:  positionUser1Mock,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup the mock expectations
			mockClient.On("GetUserWalletBalance", mock.Anything).Return(2.0)
			mockClient.On("GetUserPositions", mock.Anything).Return(tc.usersPosition, nil)
			mockClient.On("CalculateAccountMargin", mock.Anything, mock.Anything).Return(0.1)
			mockClient.On("CalculateTotalAccountNotional", mock.Anything).Return(10.0)
			mockClient.On("CalculateAccountEquity", mock.Anything, mock.Anything).Return(10.0)
			mockClient.On("UpdateOrderByPrimaryKey", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			mockClient.On("InsertNewOrder", tc.order).Return()
			mockClient.On("GetUserWalletBalance", mock.Anything).Return(12.0)
			mockClient.On("InsertMatchedPosition", mock.Anything).Return(nil)
			mockClient.On("DeleteOrderByPrimaryKey", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			mockClient.On("UpdateMarketPrice", mock.Anything, mock.Anything).Return()

			// Initialize matching engine
			matchingEngine := &MatchingEngine{
				tarantool: mockClient,
				logger:    logger.NewLoggerClient("testing-mock"),
			}
			// Call the function under test
			result := matchingEngine.MatchingEngineForLongOrder(tc.order, tc.orderBook)

			// Assert the results
			if result != tc.expectedResult {
				t.Errorf("expected %v, got %v", tc.expectedResult, result)
			}
		})
	}
}

func TestMatchingEngineForShortOrder(t *testing.T) {
	// Create a mock Tarantool client
	mockClient := new(MockTarantoolClient)

	positionUser1Mock := []*tarantool_pkg.PositionStruct{
		{
			UserID:       "user1",
			Market:       "BTC-USD",
			PositionSize: 10,
			AvgPrice:     100,
			Side:         "1",
		},
		{
			UserID:       "user1",
			Market:       "BTC-USD",
			PositionSize: 10,
			AvgPrice:     100,
			Side:         "1",
		},
	}

	// Define test cases
	testCases := []struct {
		name            string
		order           *tarantool_pkg.OrderStruct
		orderBook       []*tarantool_pkg.OrderStruct
		sortedOrderBook [][]*tarantool_pkg.OrderStruct
		usersPosition   []*tarantool_pkg.PositionStruct
		expectedResult  bool
	}{
		{
			name: "Empty order book",
			order: &tarantool_pkg.OrderStruct{
				Price:        100,
				UserId:       "user1",
				PositionSize: 10,
				Market:       "BTC-USD",
				Side:         "-1",
			},
			orderBook:       []*tarantool_pkg.OrderStruct{},
			sortedOrderBook: [][]*tarantool_pkg.OrderStruct{},
			expectedResult:  true,
			usersPosition:   positionUser1Mock,
		},
		{
			name: "Partial match",
			order: &tarantool_pkg.OrderStruct{
				Price:        100,
				UserId:       "user1",
				PositionSize: 10,
				Market:       "BTC-USD",
				Side:         "-1",
			},
			orderBook: []*tarantool_pkg.OrderStruct{
				{
					Price:        100,
					UserId:       "user2",
					PositionSize: 5,
					Market:       "BTC-USD",
					Side:         "1",
				},
			},
			sortedOrderBook: [][]*tarantool_pkg.OrderStruct{{
				{
					Price:        100,
					UserId:       "user2",
					PositionSize: 5,
					Market:       "BTC-USD",
					Side:         "1",
				},
			}},
			expectedResult: true,
			usersPosition:  positionUser1Mock,
		},
		{
			name: "Full match",
			order: &tarantool_pkg.OrderStruct{
				Price:        100,
				UserId:       "user1",
				PositionSize: 10,
				Market:       "BTC-USD",
				Side:         "-1",
			},
			orderBook: []*tarantool_pkg.OrderStruct{
				{
					Price:        100,
					UserId:       "user2",
					PositionSize: 10,
					Market:       "BTC-USD",
					Side:         "1",
				},
			},
			sortedOrderBook: [][]*tarantool_pkg.OrderStruct{{
				{
					Price:        100,
					UserId:       "user2",
					PositionSize: 10,
					Market:       "BTC-USD",
					Side:         "1",
				},
			}},
			expectedResult: true,
			usersPosition:  positionUser1Mock,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup the mock expectations
			mockClient.On("GetUserWalletBalance", mock.Anything).Return(2.0)
			mockClient.On("GetUserPositions", mock.Anything).Return(tc.usersPosition, nil)
			mockClient.On("CalculateAccountMargin", mock.Anything, mock.Anything).Return(0.1)
			mockClient.On("CalculateTotalAccountNotional", mock.Anything).Return(10.0)
			mockClient.On("CalculateAccountEquity", mock.Anything, mock.Anything).Return(10.0)
			mockClient.On("UpdateOrderByPrimaryKey", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			mockClient.On("InsertNewOrder", tc.order).Return()
			mockClient.On("GetUserWalletBalance", mock.Anything).Return(12.0)
			mockClient.On("InsertMatchedPosition", mock.Anything).Return(nil)
			mockClient.On("DeleteOrderByPrimaryKey", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			mockClient.On("UpdateMarketPrice", mock.Anything, mock.Anything).Return()

			// Initialize matching engine
			matchingEngine := &MatchingEngine{
				tarantool: mockClient,
				logger:    logger.NewLoggerClient("testing-mock"),
			}
			// Call the function under test
			result := matchingEngine.MatchingEngineForShortOrder(tc.order, tc.orderBook)

			// Assert the results
			if result != tc.expectedResult {
				t.Errorf("expected %v, got %v", tc.expectedResult, result)
			}
		})
	}
}
