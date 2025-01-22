package main

import (
	"context"
	"fmt"
	"github.com/SemyonTolkachyov/amb/src/common/db"
	"github.com/SemyonTolkachyov/amb/src/common/event"
	"github.com/SemyonTolkachyov/amb/src/querysrv/search"
	"github.com/gorilla/mux"
	"github.com/sethvargo/go-envconfig"
	"github.com/sethvargo/go-retry"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Port                 string `env:"PORT"`
	PostgresDB           string `env:"POSTGRES_DB"`
	PostgresUser         string `env:"POSTGRES_USER"`
	PostgresPassword     string `env:"POSTGRES_PASSWORD"`
	NatsAddress          string `env:"NATS_ADDRESS"`
	ElasticsearchAddress string `env:"ELASTICSEARCH_ADDRESS"`
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/messages", listMessagesHandler).
		Methods(http.MethodGet)
	router.HandleFunc("/search", searchMessagesHandler).
		Methods(http.MethodGet)
	router.Use(mux.CORSMethodMiddleware(router))
	return
}

func main() {
	ctx := context.Background()
	var cfg Config
	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to PostgreSQL
	err = retry.Do(ctx, retry.NewConstant(2*time.Second), func(ctx context.Context) error {
		addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
		repo, err := db.NewPostgres(ctx, addr)
		if err != nil {
			log.Println(err)
			return retry.RetryableError(err)
		}
		db.SetRepository(repo)
		return nil
	})
	if err != nil {
		return
	}
	defer db.Close(ctx)

	// Connect to ElasticSearch
	err = retry.Do(ctx, retry.NewConstant(2*time.Second), func(ctx context.Context) error {
		es, err := search.NewElastic(fmt.Sprintf("http://%s", cfg.ElasticsearchAddress))
		if err != nil {
			log.Println(err)
			return retry.RetryableError(err)
		}
		search.SetRepository(es)
		return nil
	})
	defer search.Close()

	// Connect to Nats
	err = retry.Do(ctx, retry.NewConstant(2*time.Second), func(ctx context.Context) error {
		es, err := event.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))
		if err != nil {
			log.Println(err)
			return retry.RetryableError(err)
		}

		err = es.OnMessageCreated(onMessageCreated)
		if err != nil {
			log.Println(err)
			return retry.RetryableError(err)
		}

		event.SetEventStore(es)
		return nil
	})
	defer event.Close()

	// Run HTTP server
	router := newRouter()
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router); err != nil {
		log.Fatal(err)
	}
}
