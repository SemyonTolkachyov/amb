package main

import (
	"context"
	"fmt"
	"github.com/SemyonTolkachyov/amb/src/common/event"
	"github.com/sethvargo/go-envconfig"
	"github.com/sethvargo/go-retry"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Port        string `env:"PORT"`
	NatsAddress string `env:"NATS_ADDRESS"`
}

func main() {
	ctx := context.Background()
	var cfg Config
	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to Nats
	hub := newHub()
	err = retry.Do(ctx, retry.NewConstant(2*time.Second), func(ctx context.Context) error {
		es, err := event.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))
		if err != nil {
			log.Println(err)
			return retry.RetryableError(err)
		}

		// Push messages to clients
		err = es.OnMessageCreated(func(m event.MessageCreatedEvent) {
			log.Printf("Message received: %v\n", m)
			hub.broadcast(newMessageCreatedModel(m.ID, m.Body, m.CreatedAt), nil)
		})
		if err != nil {
			log.Println(err)
			return retry.RetryableError(err)
		}

		event.SetEventStore(es)
		return nil
	})
	defer event.Close()

	// Run WebSocket server
	go hub.run()
	http.HandleFunc("/pusher", hub.handleWebSocket)
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
