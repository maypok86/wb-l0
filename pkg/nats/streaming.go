package nats

import (
	"fmt"
	"sync"

	"github.com/maypok86/wb-l0/pkg/logger"
	"github.com/nats-io/stan.go"
)

type Streaming struct {
	mutex         sync.Mutex
	conn          stan.Conn
	subscriptions map[string]stan.Subscription
}

func NewStreaming(config Config) (*Streaming, error) {
	natsURL := fmt.Sprintf("nats://%s:%s", config.Host, config.Port)
	conn, err := stan.Connect(config.ClusterID, config.ClientID, stan.NatsURL(natsURL))
	if err != nil {
		return nil, fmt.Errorf("can not connect to nats-streaming-server: %w", err)
	}
	return &Streaming{
		conn:          conn,
		subscriptions: make(map[string]stan.Subscription),
	}, nil
}

func (s *Streaming) Subscribe(channel string, msgHandler stan.MsgHandler) error {
	subscription, err := s.conn.Subscribe(channel, msgHandler, stan.StartWithLastReceived())
	if err != nil {
		return fmt.Errorf("can not subscribe: %w", err)
	}
	s.mutex.Lock()
	s.subscriptions[channel] = subscription
	s.mutex.Unlock()
	return nil
}

func (s *Streaming) Unsubscribe(channel string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	subscription := s.subscriptions[channel]
	if subscription != nil {
		if err := subscription.Unsubscribe(); err != nil {
			logger.Error(fmt.Errorf("can not unsubscribe: %w", err))
		}
		delete(s.subscriptions, channel)
	} else {
		logger.Warnf("not found subscription for channel %s", channel)
	}
}

func (s *Streaming) UnsubscribeAll() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for ch, sub := range s.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			logger.Error(fmt.Errorf("can not unsubscribe: %w", err))
		}
		delete(s.subscriptions, ch)
	}
}

func (s *Streaming) Close() error {
	return s.conn.Close()
}
