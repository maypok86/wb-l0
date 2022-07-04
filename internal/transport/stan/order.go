package stan

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/maypok86/wb-l0/internal/entity"
	"github.com/maypok86/wb-l0/pkg/logger"
	"github.com/nats-io/stan.go"
)

func (r Router) newOrder(ctx context.Context) error {
	return r.natsStreaming.Subscribe("orders", func(msg *stan.Msg) {
		order := &entity.Order{}
		if err := json.Unmarshal(msg.Data, order); err != nil {
			logger.Error(fmt.Errorf("can not unmarshal order: %w", err))
			return
		}
		_, err := r.orderUsecase.CreateOrder(ctx, order, time.Hour)
		if err != nil {
			logger.Error(fmt.Errorf("can not create order: %w", err))
			return
		}
		logger.Info("successful create order")
	})
}
