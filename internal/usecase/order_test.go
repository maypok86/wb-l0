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

func fakeTTL(t *testing.T) time.Duration {
	t.Helper()

	var ttl time.Duration
	require.NoError(t, faker.FakeData(&ttl))
	return ttl
}

func TestOrderUsecase_CreateOrder(t *testing.T) {
	order, cache, repo := mockOrderUsecase(t)
	validOrder := validOrder(t)
	ttl := fakeTTL(t)

	tests := []struct {
		name   string
		mock   func()
		order  *entity.Order
		result *entity.Order
		err    error
	}{
		{
			name: "error in repo",
			mock: func() {
				repo.EXPECT().CreateOrder(context.Background(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			order:  nil,
			result: nil,
			err:    errors.New("can not create in repository: some error"),
		},
		{
			name: "error in cache",
			mock: func() {
				repo.EXPECT().CreateOrder(context.Background(), gomock.Any()).Return(validOrder, nil)
				cache.EXPECT().Set(validOrder.OrderUID, gomock.Any(), gomock.Any()).Return(errors.New("some error"))
			},
			order:  validOrder,
			result: nil,
			err:    errors.New("can not create order in cache: some error"),
		},
		{
			name: "success",
			mock: func() {
				repo.EXPECT().CreateOrder(context.Background(), gomock.Any()).Return(validOrder, nil)
				cache.EXPECT().Set(validOrder.OrderUID, gomock.Any(), gomock.Any()).Return(nil)
			},
			order:  validOrder,
			result: validOrder,
			err:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := order.CreateOrder(context.Background(), tt.order, ttl)
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

	orderEntity := validOrder(t)

	tests := []struct {
		name   string
		mock   func()
		result *entity.Order
		err    error
	}{
		{
			name: "contains in cache",
			mock: func() {
				cache.EXPECT().Get(orderEntity.OrderUID).Return(orderEntity, nil)
			},
			result: orderEntity,
			err:    nil,
		},
		{
			name: "error in repo",
			mock: func() {
				cache.EXPECT().Get(orderEntity.OrderUID).Return(nil, errors.New("cache error"))
				repo.EXPECT().
					GetOrderByID(context.Background(), orderEntity.OrderUID).
					Return(nil, errors.New("repo error"))
			},
			result: nil,
			err:    errors.New("can not get order by id: repo error"),
		},
		{
			name: "contains in repo",
			mock: func() {
				cache.EXPECT().Get(orderEntity.OrderUID).Return(nil, errors.New("cache error"))
				repo.EXPECT().GetOrderByID(context.Background(), orderEntity.OrderUID).Return(orderEntity, nil)
			},
			result: orderEntity,
			err:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := order.GetOrderByID(context.Background(), orderEntity.OrderUID)
			require.Equal(t, tt.result, result)
			if tt.err != nil {
				require.Equal(t, tt.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestOrderUsecase_LoadDBToCache(t *testing.T) {
	order, cache, repo := mockOrderUsecase(t)

	ttl := fakeTTL(t)
	const length = 10
	orders := validOrders(t, length)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(orders), func(i, j int) {
		orders[i], orders[j] = orders[j], orders[i]
	})
	index := rand.Intn(len(orders))

	tests := []struct {
		name string
		mock func()
		err  error
	}{
		{
			name: "error in repo",
			mock: func() {
				repo.EXPECT().GetAllOrders(context.Background()).Return(nil, errors.New("some error"))
			},
			err: errors.New("can not get all orders: some error"),
		},
		{
			name: "error in cache",
			mock: func() {
				repo.EXPECT().GetAllOrders(context.Background()).Return(orders, nil)
				for i, order := range orders {
					if i == index {
						cache.EXPECT().Set(order.OrderUID, order, ttl).Return(errors.New("some error"))
					} else {
						cache.EXPECT().Set(order.OrderUID, order, ttl).Return(nil)
					}
				}
			},
			err: errors.New("can not set order to cache: some error"),
		},
		{
			name: "valid orders",
			mock: func() {
				repo.EXPECT().GetAllOrders(context.Background()).Return(orders, nil)
				for _, order := range orders {
					cache.EXPECT().Set(order.OrderUID, order, ttl).Return(nil)
				}
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := order.LoadDBToCache(context.Background(), ttl)
			if tt.err != nil {
				require.Equal(t, tt.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}
