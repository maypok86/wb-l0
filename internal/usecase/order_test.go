package usecase_test

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/maypok86/wb-l0/internal/entity"
	"github.com/maypok86/wb-l0/internal/usecase"
	"github.com/stretchr/testify/require"
)

func mockOrderUsecase(t *testing.T) (usecase.OrderUsecase, *MockOrderCache, *MockOrderRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockOrderRepository(mockCtl)
	cache := NewMockOrderCache(mockCtl)
	order := usecase.NewOrderUsecase(cache, repo)

	return order, cache, repo
}

func readFile(t *testing.T, filename string) []byte {
	t.Helper()

	data, err := ioutil.ReadFile(filename)
	require.NoError(t, err)
	return data
}

func invalidOrderData(t *testing.T) []byte {
	t.Helper()

	return readFile(t, "fixtures/invalid_model.json")
}

func validOrderData(t *testing.T) []byte {
	t.Helper()

	return readFile(t, "fixtures/model.json")
}

func validOrder(t *testing.T) *entity.Order {
	t.Helper()

	data := validOrderData(t)
	order := &entity.Order{}
	require.NoError(t, json.Unmarshal(data, order))
	return order
}

func validOrders(t *testing.T, length int) []*entity.Order {
	t.Helper()

	orders := make([]*entity.Order, length)
	require.NoError(t, faker.FakeData(&orders))
	return orders
}

func getJSONData(t *testing.T, v any) []byte {
	t.Helper()

	data, err := json.Marshal(v)
	require.NoError(t, err)
	return data
}

func getOrdersData(t *testing.T, orders []*entity.Order) [][]byte {
	t.Helper()

	data := make([][]byte, len(orders))
	for i, order := range orders {
		data[i] = getJSONData(t, order)
	}
	return data
}

func TestOrderUsecase_CreateOrder(t *testing.T) {
	order, cache, repo := mockOrderUsecase(t)
	const id = 0
	const ttl = 60

	tests := []struct {
		name   string
		mock   func()
		data   []byte
		result *entity.Order
		err    error
	}{
		{
			name: "error in repo",
			mock: func() {
				repo.EXPECT().CreateOrder(context.Background(), gomock.Any()).Return(id, errors.New("some error"))
			},
			data:   nil,
			result: nil,
			err:    errors.New("can not create in repository: some error"),
		},
		{
			name: "not valid json",
			mock: func() {
				repo.EXPECT().CreateOrder(context.Background(), gomock.Any()).Return(id, nil)
			},
			data:   []byte("<>"),
			result: nil,
			err:    errors.New("can not parse json order: invalid character '<' looking for beginning of value"),
		},
		{
			name: "not valid json order data",
			mock: func() {
				repo.EXPECT().CreateOrder(context.Background(), gomock.Any()).Return(id, nil)
			},
			data:   invalidOrderData(t),
			result: nil,
			err:    errors.New("can not parse json order: json: cannot unmarshal number into Go struct field Order.locale of type string"),
		},
		{
			name: "error in cache",
			mock: func() {
				repo.EXPECT().CreateOrder(context.Background(), gomock.Any()).Return(id, nil)
				cache.EXPECT().Set(id, gomock.Any(), gomock.Any()).Return(errors.New("some error"))
			},
			data:   validOrderData(t),
			result: nil,
			err:    errors.New("can not create order in cache: some error"),
		},
		{
			name: "success",
			mock: func() {
				repo.EXPECT().CreateOrder(context.Background(), gomock.Any()).Return(id, nil)
				cache.EXPECT().Set(id, gomock.Any(), gomock.Any()).Return(nil)
			},
			data:   validOrderData(t),
			result: validOrder(t),
			err:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := order.CreateOrder(context.Background(), tt.data, ttl)
			require.Equal(t, tt.result, result)
			if tt.err != nil {
				require.Equal(t, tt.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestOrderUsecase_GetOrderByID(t *testing.T) {
	order, cache, repo := mockOrderUsecase(t)
	const id = 0

	invalidOrderData := invalidOrderData(t)
	validOrderData := validOrderData(t)
	validOrder := validOrder(t)

	tests := []struct {
		name   string
		mock   func()
		result *entity.Order
		err    error
	}{
		{
			name: "contains in cache",
			mock: func() {
				cache.EXPECT().Get(id).Return(validOrder, nil)
			},
			result: validOrder,
			err:    nil,
		},
		{
			name: "error in repo",
			mock: func() {
				cache.EXPECT().Get(id).Return(nil, errors.New("cache error"))
				repo.EXPECT().GetOrderByID(context.Background(), id).Return(nil, errors.New("repo error"))
			},
			result: nil,
			err:    errors.New("can not get order by id: repo error"),
		},
		{
			name: "not valid json",
			mock: func() {
				cache.EXPECT().Get(id).Return(nil, errors.New("cache error"))
				repo.EXPECT().GetOrderByID(context.Background(), id).Return([]byte("<>"), nil)
			},
			result: nil,
			err:    errors.New("can not parse json order: invalid character '<' looking for beginning of value"),
		},
		{
			name: "not valid json order data",
			mock: func() {
				cache.EXPECT().Get(id).Return(nil, errors.New("cache error"))
				repo.EXPECT().GetOrderByID(context.Background(), id).Return(invalidOrderData, nil)
			},
			result: nil,
			err:    errors.New("can not parse json order: json: cannot unmarshal number into Go struct field Order.locale of type string"),
		},
		{
			name: "contains in repo",
			mock: func() {
				cache.EXPECT().Get(id).Return(nil, errors.New("cache error"))
				repo.EXPECT().GetOrderByID(context.Background(), id).Return(validOrderData, nil)
			},
			result: validOrder,
			err:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := order.GetOrderByID(context.Background(), id)
			require.Equal(t, tt.result, result)
			if tt.err != nil {
				require.Equal(t, tt.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestOrderUsecase_GetAllOrders(t *testing.T) {
	order, _, repo := mockOrderUsecase(t)

	const length = 10
	orders := validOrders(t, length)
	ordersData := getOrdersData(t, orders)
	invalidOrdersData := make([][]byte, len(ordersData)+1)
	copy(invalidOrdersData, ordersData)
	invalidOrdersData[len(invalidOrdersData)-1] = invalidOrderData(t)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(invalidOrdersData), func(i, j int) {
		invalidOrdersData[i], invalidOrdersData[j] = invalidOrdersData[j], invalidOrdersData[i]
	})

	tests := []struct {
		name   string
		mock   func()
		result []*entity.Order
		err    error
	}{
		{
			name: "error in repo",
			mock: func() {
				repo.EXPECT().GetAllOrders(context.Background()).Return(nil, errors.New("some error"))
			},
			result: nil,
			err:    errors.New("can not get all orders from repo: some error"),
		},
		{
			name: "invalid json order data",
			mock: func() {
				repo.EXPECT().GetAllOrders(context.Background()).Return(invalidOrdersData, nil)
			},
			result: nil,
			err:    errors.New("can not parse json order: json: cannot unmarshal number into Go struct field Order.locale of type string"),
		},
		{
			name: "valid orders",
			mock: func() {
				repo.EXPECT().GetAllOrders(context.Background()).Return(ordersData, nil)
			},
			result: orders,
			err:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := order.GetAllOrders(context.Background())
			if result != nil {
				require.Equal(t, getJSONData(t, tt.result), getJSONData(t, result))
			} else {
				require.Equal(t, tt.result, result)
			}
			if tt.err != nil {
				require.Equal(t, tt.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}
