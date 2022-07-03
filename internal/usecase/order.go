package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/maypok86/wb-l0/internal/entity"
)

//go:generate mockgen -source=order.go -destination=mock_test.go -package=usecase_test

type OrderCache interface {
	Set(string, *entity.Order, time.Duration) error
	Get(string) (*entity.Order, error)
}

type OrderRepository interface {
	CreateOrder(context.Context, *entity.Order) (*entity.Order, error)
	GetOrderByID(context.Context, string) (*entity.Order, error)
	GetAllOrders(context.Context) ([]*entity.Order, error)
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

func (ou OrderUsecase) CreateOrder(ctx context.Context, order *entity.Order, ttl time.Duration) (*entity.Order, error) {
	order, err := ou.repository.CreateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("can not create in repository: %w", err)
	}
	if err := ou.cache.Set(order.OrderUID, order, ttl); err != nil {
		return nil, fmt.Errorf("can not create order in cache: %w", err)
	}
	return order, nil
}

func (ou OrderUsecase) GetOrderByID(ctx context.Context, orderUID string) (*entity.Order, error) {
	order, err := ou.cache.Get(orderUID)
	if err != nil {
		order, err = ou.repository.GetOrderByID(ctx, orderUID)
		if err != nil {
			return nil, fmt.Errorf("can not get order by id: %w", err)
		}
	}
	return order, nil
}

func (ou OrderUsecase) LoadDBToCache(ctx context.Context, ttl time.Duration) error {
	orders, err := ou.repository.GetAllOrders(ctx)
	if err != nil {
		return fmt.Errorf("can not get all orders: %w", err)
	}

	for _, order := range orders {
		if err := ou.cache.Set(order.OrderUID, order, ttl); err != nil {
			return fmt.Errorf("can not set order to cache: %w", err)
		}
	}

	return nil
}
