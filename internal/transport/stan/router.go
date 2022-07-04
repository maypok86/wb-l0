package stan

import (
	"context"
	"time"

	"github.com/maypok86/wb-l0/internal/entity"
	"github.com/maypok86/wb-l0/pkg/nats"
)

type OrderUsecase interface {
	CreateOrder(context.Context, *entity.Order, time.Duration) (*entity.Order, error)
}

type Router struct {
	natsStreaming *nats.Streaming
	orderUsecase  OrderUsecase
}

func NewRouter(natsStreaming *nats.Streaming, orderUsecase OrderUsecase) Router {
	return Router{
		natsStreaming: natsStreaming,
		orderUsecase:  orderUsecase,
	}
}

func (r Router) Init(ctx context.Context) error {
	return r.newOrder(ctx)
}
