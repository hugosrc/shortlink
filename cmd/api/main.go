package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v9"
	"github.com/go-zookeeper/zk"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/hugosrc/shortlink/config"
	"github.com/hugosrc/shortlink/internal/adapter/cassandra"
	"github.com/hugosrc/shortlink/internal/adapter/cassandra/repository"
	kafkaAdapter "github.com/hugosrc/shortlink/internal/adapter/kafka"
	"github.com/hugosrc/shortlink/internal/adapter/keycloak"
	redisAdapter "github.com/hugosrc/shortlink/internal/adapter/redis"
	"github.com/hugosrc/shortlink/internal/adapter/zookeeper"
	"github.com/hugosrc/shortlink/internal/core/service"
	"github.com/hugosrc/shortlink/internal/handler/rest"
	"github.com/jxskiss/base62"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("couldn't initialize zap logger: %v", err)
		os.Exit(1)
	}

	config, err := config.Init()
	if err != nil {
		logger.Error("couldn't initialize configuration", zap.Error(err))
		os.Exit(1)
	}

	cassandraConn, err := cassandra.New(config)
	if err != nil {
		logger.Error("couldn't connect to cassandra", zap.Error(err))
		os.Exit(1)
	}

	redisConn, err := redisAdapter.New(config)
	if err != nil {
		logger.Error("couldn't connect to redis", zap.Error(err))
		os.Exit(1)
	}

	zookeeperConn, err := zookeeper.New(config)
	if err != nil {
		logger.Error("couldn't connect to zookeeper", zap.Error(err))
		os.Exit(1)
	}

	keycloakAuth := keycloak.NewOpenIDAuth(config)

	kafkaProducer, err := kafkaAdapter.NewProducer(config)
	if err != nil {
		logger.Error("couldn't connect to kafka", zap.Error(err))
		os.Exit(1)
	}

	logMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("received request", zap.String("method", r.Method), zap.String("uri", r.RequestURI))
			next.ServeHTTP(w, r)
		})
	}

	server := newServer(serverConf{
		Address:      fmt.Sprintf(":%d", 3000),
		Auth:         keycloakAuth,
		Cassandra:    cassandraConn,
		Redis:        redisConn,
		Zookeeper:    zookeeperConn,
		Kafka:        kafkaProducer,
		MetricsTopic: config.GetString("KAFKA_METRICS_PRODUCER_TOPIC_NAME"),
		Middlewares:  []func(next http.Handler) http.Handler{logMiddleware},
	})

	go func() {
		logger.Info("server started")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("couldn't start the server", zap.Error(err))
		}
	}()

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		if err := redisConn.Close(); err != nil {
			logger.Error("error closing redis connection", zap.Error(err))
		}

		cassandraConn.Close()
		zookeeperConn.Close()
		kafkaProducer.Close()
		cancel()
		close(done)
	}()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("couldn't gracefully shutdown the server", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("shutdown performed successfully")
}

type serverConf struct {
	Address      string
	Auth         *keycloak.OpenIDAuth
	Cassandra    *gocql.Session
	Redis        *redis.Client
	Zookeeper    *zk.Conn
	Kafka        *kafka.Producer
	MetricsTopic string
	Middlewares  []func(next http.Handler) http.Handler
}

func newServer(conf serverConf) *http.Server {
	r := mux.NewRouter()

	for _, middleware := range conf.Middlewares {
		r.Use(middleware)
	}

	counter := zookeeper.NewCounter(conf.Zookeeper)
	_ = counter.UpdateCounterBase() // TODO: goroutine

	encoder := base62.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	caching := redisAdapter.NewRedisCaching(conf.Redis)
	repo := repository.NewLinkRepository(conf.Cassandra)

	service := service.NewLinkService(counter, encoder, caching, repo)

	metricsProducer := kafkaAdapter.NewKafkaMetricsProducer(conf.MetricsTopic, conf.Kafka)

	rest.NewLinkHandler(conf.Auth, metricsProducer, service).Register(r)

	return &http.Server{
		Addr:              conf.Address,
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
