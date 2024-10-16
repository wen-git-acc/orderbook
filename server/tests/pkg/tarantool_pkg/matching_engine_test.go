// package tarantool_pkg_test
package a

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/wen-git-acc/orderbook/pkg/tarantool_pkg"
// )

// type MockTarantoolClient struct {
// 	mock.Mock
// 	tarantool_pkg.TarantoolClient
// }

// func (m *MockTarantoolClient) SortOrderBook(orderbooks []*tarantool_pkg.OrderStruct) [][]*tarantool_pkg.OrderStruct {
// 	args := m.Called(orderbooks)
// 	return args.Get(0).([][]*tarantool_pkg.OrderStruct)
// }

// func (m *MockTarantoolClient) InsertNewOrder(order *tarantool_pkg.OrderStruct) {
// 	m.Called(order)
// }

// func (m *MockTarantoolClient) updateUserWalletAmount(userId string, amountToDeduct float64) {
// 	m.Called(userId, amountToDeduct)
// }

// func (m *MockTarantoolClient) checkAccountMargin(executionDetail *tarantool_pkg.ExecutionDetailsStruct) bool {
// 	args := m.Called(executionDetail)
// 	return args.Bool(0)
// }

// func (m *MockTarantoolClient) InsertMatchedPosition(position *tarantool_pkg.PositionStruct) {
// 	m.Called(position)
// }

// func (m *MockTarantoolClient) DeleteOrderByPrimaryKey(userId string, price float64, side string, market string) {
// 	m.Called(userId, price, side, market)
// }

// func (m *MockTarantoolClient) UpdateMarketPrice(market string, price float64) {
// 	m.Called(market, price)
// }

// func TestMatchingEngineForShortOrder(t *testing.T) {
// 	mockClient := new(MockTarantoolClient)
// 	order := &tarantool_pkg.OrderStruct{
// 		UserId:       "user1",
// 		Market:       "BTC-USD",
// 		Price:        50000,
// 		PositionSize: 1,
// 	}

// 	t.Run("Empty Order Book", func(t *testing.T) {
// 		mockClient.On("SortOrderBook", mock.Anything).Return([][]*tarantool_pkg.OrderStruct{})
// 		result := mockClient.MatchingEngineForShortOrder(order, []*tarantool_pkg.OrderStruct{})
// 		assert.False(t, result)
// 		mockClient.AssertExpectations(t)
// 	})

// 	t.Run("Last Element Price Smaller", func(t *testing.T) {
// 		orderBook := []*tarantool_pkg.OrderStruct{
// 			{Price: 49000},
// 		}
// 		mockClient.On("SortOrderBook", orderBook).Return([][]*tarantool_pkg.OrderStruct{{orderBook[0]}})
// 		mockClient.On("InsertNewOrder", order).Return()
// 		mockClient.On("updateUserWalletAmount", order.UserId, order.PositionSize*order.Price).Return()

// 		result := mockClient.MatchingEngineForShortOrder(order, orderBook)
// 		assert.True(t, result)
// 		mockClient.AssertExpectations(t)
// 	})

// 	t.Run("Partial Match", func(t *testing.T) {
// 		orderBook := []*tarantool_pkg.OrderStruct{
// 			{UserId: "maker1", Price: 50000, PositionSize: 0.5},
// 		}
// 		mockClient.On("SortOrderBook", orderBook).Return([][]*tarantool_pkg.OrderStruct{{orderBook[0]}})
// 		mockClient.On("checkAccountMargin", mock.Anything).Return(true)
// 		mockClient.On("InsertMatchedPosition", mock.Anything).Return()
// 		mockClient.On("DeleteOrderByPrimaryKey", "maker1", 50000.0, "-1", "BTC-USD").Return()
// 		mockClient.On("UpdateMarketPrice", "BTC-USD", 50000.0).Return()

// 		result := mockClient.MatchingEngineForShortOrder(order, orderBook)
// 		assert.True(t, result)
// 		mockClient.AssertExpectations(t)
// 	})

// 	t.Run("Full Match", func(t *testing.T) {
// 		orderBook := []*tarantool_pkg.OrderStruct{
// 			{UserId: "maker1", Price: 50000, PositionSize: 1},
// 		}
// 		mockClient.On("SortOrderBook", orderBook).Return([][]*tarantool_pkg.OrderStruct{{orderBook[0]}})
// 		mockClient.On("checkAccountMargin", mock.Anything).Return(true)
// 		mockClient.On("InsertMatchedPosition", mock.Anything).Return()
// 		mockClient.On("UpdateOrderByPrimaryKey", "maker1", 50000.0, "-1", "BTC-USD", 0.0).Return()
// 		mockClient.On("UpdateMarketPrice", "BTC-USD", 50000.0).Return()

// 		result := mockClient.MatchingEngineForShortOrder(order, orderBook)
// 		assert.True(t, result)
// 		mockClient.AssertExpectations(t)
// 	})
// }
