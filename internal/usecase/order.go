package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/maypok86/wb-l0/internal/entity"
)

//go:generate mockgen -source=order.go -destination=mock_test.go -package=usecase_test

type OrderCache interface {
	Set(int, *entity.Order, int64) error
	Get(int) (*entity.Order, error)
}

type OrderRepository interface {
	CreateOrder(context.Context, []byte) (int, error)
	GetOrderByID(context.Context, int) ([]byte, error)
	GetAllOrders(context.Context) ([][]byte, error)
}

type OrderUsecase struct {
	cache      OrderCache
	repository OrderRepository
}

func NewOrderUsecase(cache OrderCache, repository OrderRepository) OrderUsecase {
	return OrderUsecase{
		cache:      cache,
		repository: repository,
	}
}

func (ou OrderUsecase) CreateOrder(ctx context.Context, data []byte, ttl int64) (*entity.Order, error) {
	id, err := ou.repository.CreateOrder(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("can not create in repository: %w", err)
	}
	order := &entity.Order{}
	if err := json.Unmarshal(data, order); err != nil {
		return nil, fmt.Errorf("can not parse json order: %w", err)
	}
	if err := ou.cache.Set(id, order, ttl); err != nil {
		return nil, fmt.Errorf("can not create order in cache: %w", err)
	}
	return order, nil
}

func (ou OrderUsecase) GetOrderByID(ctx context.Context, id int) (*entity.Order, error) {
	order, err := ou.cache.Get(id)
	if err != nil {
		data, err := ou.repository.GetOrderByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("can not get order by id: %w", err)
		}
		order = &entity.Order{}
		if err := json.Unmarshal(data, order); err != nil {
			return nil, fmt.Errorf("can not parse json order: %w", err)
		}
	}
	return order, nil
}

func (ou OrderUsecase) GetAllOrders(ctx context.Context) ([]*entity.Order, error) {
	const defaultOrdersCapacity = 64
	orders := make([]*entity.Order, 0, defaultOrdersCapacity)
	orderBytes, err := ou.repository.GetAllOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("can not get all orders from repo: %w", err)
	}
	for _, orderByte := range orderBytes {
		order := &entity.Order{}
		if err := json.Unmarshal(orderByte, order); err != nil {
			return nil, fmt.Errorf("can not parse json order: %w", err)
		}
		orders = append(orders, order)
	}
	return orders, nil
}
