package api

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maypok86/wb-l0/internal/config"
	"github.com/maypok86/wb-l0/internal/transport/http"
	"github.com/maypok86/wb-l0/pkg/httpserver"
	"github.com/maypok86/wb-l0/pkg/logger"
)

type App struct {
	ctx        context.Context
	httpServer httpserver.Server
}

func New(ctx context.Context) (App, error) {
	cfg := config.Get()
	handler := http.NewHandler()
	return App{
		ctx: ctx,
		httpServer: httpserver.New(
			handler.GetHTTPHandler(),
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
	return a.httpServer.Stop(a.ctx, httpShutdownTimeout)
}
