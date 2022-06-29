package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"
)

type Config struct {
	host           string
	port           string
	maxHeaderBytes int
	readTimeout    time.Duration
	writeTimeout   time.Duration
}

func NewConfig(host, port string, maxHeaderBytes int, readTimeout, writeTimeout time.Duration) Config {
	return Config{
		host:           host,
		port:           port,
		maxHeaderBytes: maxHeaderBytes,
		readTimeout:    readTimeout,
		writeTimeout:   writeTimeout,
	}
}

type Server struct {
	httpServer *http.Server
}

func New(httpHandler http.Handler, cfg Config) Server {
	return Server{
		httpServer: &http.Server{
			Addr:           net.JoinHostPort(cfg.host, cfg.port),
			Handler:        httpHandler,
			MaxHeaderBytes: cfg.maxHeaderBytes << 20,
			ReadTimeout:    cfg.readTimeout,
			WriteTimeout:   cfg.writeTimeout,
		},
	}
}

func (s Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s Server) Stop(ctx context.Context, shutdownTimeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}
