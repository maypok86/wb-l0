package postgres

import (
	"fmt"
	"time"
)

type ConnectionConfig struct {
	host     string
	port     string
	dbname   string
	username string
	password string
	sslmode  string
}

func NewConnectionConfig(host, port, dbname, username, password, sslmode string) ConnectionConfig {
	return ConnectionConfig{
		host:     host,
		port:     port,
		dbname:   dbname,
		username: username,
		password: password,
		sslmode:  sslmode,
	}
}

func (cc ConnectionConfig) getDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cc.username,
		cc.password,
		cc.host,
		cc.port,
		cc.dbname,
		cc.sslmode,
	)
}

type poolConfig struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
}

func getDefaultPoolConfig() *poolConfig {
	return &poolConfig{
		maxPoolSize:  1,
		connAttempts: 10,
		connTimeout:  time.Second,
	}
}

type Option func(*poolConfig)

func MaxPoolSize(maxPoolSize int) Option {
	return func(c *poolConfig) {
		c.maxPoolSize = maxPoolSize
	}
}

func ConnAttempts(connAttempts int) Option {
	return func(c *poolConfig) {
		c.connAttempts = connAttempts
	}
}

func ConnTimeout(connTimeout time.Duration) Option {
	return func(c *poolConfig) {
		c.connTimeout = connTimeout
	}
}
