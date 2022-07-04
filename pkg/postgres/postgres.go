package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	Builder sq.StatementBuilderType
	Pool    *pgxpool.Pool
}

func New(ctx context.Context, connectionConfig ConnectionConfig, opts ...Option) (*Postgres, error) {
	cfg := getDefaultPoolConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	instance := &Postgres{}
	instance.Builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	dsn := connectionConfig.getDSN()
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}
	poolCfg.MaxConns = int32(cfg.maxPoolSize)

	for cfg.connAttempts > 0 {
		instance.Pool, err = pgxpool.ConnectConfig(ctx, poolCfg)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", cfg.connAttempts)
		time.Sleep(cfg.connTimeout)

		cfg.connAttempts--
	}

	if err != nil {
		return nil, errors.New("all attempts are exceeded. Unable to connect to instance")
	}

	return instance, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
