package api

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maypok86/wb-l0/internal/cache"
	"github.com/maypok86/wb-l0/internal/config"
	"github.com/maypok86/wb-l0/internal/repository"
	"github.com/maypok86/wb-l0/internal/transport/http"
	"github.com/maypok86/wb-l0/internal/transport/stan"
	"github.com/maypok86/wb-l0/internal/usecase"
	"github.com/maypok86/wb-l0/pkg/httpserver"
	"github.com/maypok86/wb-l0/pkg/logger"
	"github.com/maypok86/wb-l0/pkg/nats"
	"github.com/maypok86/wb-l0/pkg/postgres"
)

type App struct {
	ctx           context.Context
	httpServer    httpserver.Server
	db            *postgres.Postgres
	natsStreaming *nats.Streaming
}

func New(ctx context.Context) (App, error) {
	cfg := config.Get()

	memoryCache := cache.NewMemoryCache()
	postgresInstance, err := postgres.New(
		ctx,
		postgres.NewConnectionConfig(
			cfg.Postgres.Host,
			cfg.Postgres.Port,
			cfg.Postgres.DBName,
			cfg.Postgres.User,
			cfg.Postgres.Password,
			cfg.Postgres.SSLMode,
		),
	)
	if err != nil {
		return App{}, fmt.Errorf("can not connect to postgres: %w", err)
	}
	orderRepository := repository.NewOrderPostgresRepository(postgresInstance)
	orderUsecase := usecase.NewOrderUsecase(memoryCache, orderRepository)

	if err := orderUsecase.LoadDBToCache(ctx, time.Hour); err != nil {
		return App{}, fmt.Errorf("can not load db to cache: %w", err)
	}

	natsStreaming, err := nats.NewStreaming(nats.NewConfig(cfg.STAN.Host, cfg.STAN.Port, cfg.STAN.ClusterID, cfg.STAN.ClientID))
	if err != nil {
		return App{}, fmt.Errorf("can not connect to stan-streaming-server: %w", err)
	}
	router := stan.NewRouter(natsStreaming, orderUsecase)
	if err := router.Init(ctx); err != nil {
		return App{}, fmt.Errorf("can not init nats router: %w", err)
	}

	handler := http.NewHandler(orderUsecase)
	return App{
		ctx: ctx,
		httpServer: httpserver.New(
			handler.Init(),
			httpserver.NewConfig(
				cfg.HTTP.Host,
				cfg.HTTP.Port,
				cfg.HTTP.MaxHeaderBytes,
				cfg.HTTP.ReadTimeout,
				cfg.HTTP.WriteTimeout,
			),
		),
	}, nil
}

func (a App) Run() error {
	eChan := make(chan error)
	interrupt := make(chan os.Signal, 1)

	logger.Info("Http server is starting")
	go func() {
		if err := a.httpServer.Start(); err != nil {
			eChan <- fmt.Errorf("failed to listen and serve: %w", err)
		}
	}()

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-eChan:
		return fmt.Errorf("wb-l0 started failed: %w", err)
	case <-interrupt:
	}

	const httpShutdownTimeout = 5 * time.Second
	if err := a.httpServer.Stop(a.ctx, httpShutdownTimeout); err != nil {
		return err
	}
	a.db.Close()
	a.natsStreaming.UnsubscribeAll()
	return a.natsStreaming.Close()
}
