package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/maypok86/wb-l0/internal/entity"
	"github.com/nats-io/stan.go"
)

var (
	delay     int
	clusterID string
	clientID  string
	channel   string
)

func init() {
	flag.IntVar(&delay, "delay", 0, "delay between publications")
	flag.StringVar(&clusterID, "cluster", "wb-cluster", "ClusterID for nats-streaming-server")
	flag.StringVar(&clientID, "client", "wb-client", "ClientID for nats-streaming-server")
	flag.StringVar(&channel, "channel", "orders", "Channel for subscribe")
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func publishOrder(channel string, conn stan.Conn, eChan chan<- error) {
	rand.Seed(time.Now().Unix())
	for {
		order := &entity.Order{}
		if err := faker.FakeData(order); err != nil {
			eChan <- err
		}

		data, err := json.Marshal(order)
		if err != nil {
			eChan <- err
		}

		i := rand.Intn(10)
		if i == 0 {
			data = []byte("<>")
			order.OrderUID = "not valid json"
		}

		if err := conn.Publish(channel, data); err != nil {
			eChan <- err
		}

		if i == 0 {
			log.Println("Not valid json published")
		} else {
			log.Printf("Order with uid '%s' published\n", order.OrderUID)
		}

		<-time.After(time.Duration(delay) * time.Second)
	}
}

func run() error {
	conn, err := stan.Connect(clusterID, clientID, stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		return err
	}
	defer conn.Close()

	interrupt := make(chan os.Signal, 1)
	eChan := make(chan error)

	go publishOrder(channel, conn, eChan)

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-eChan:
		return err
	case <-interrupt:
		return nil
	}
}
