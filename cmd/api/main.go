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

	"github.com/go-redis/redis/v9"
	"github.com/go-zookeeper/zk"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/hugosrc/shortlink/config"
	"github.com/hugosrc/shortlink/internal/adapter/cassandra"
	"github.com/hugosrc/shortlink/internal/adapter/cassandra/repository"
	"github.com/hugosrc/shortlink/internal/adapter/keycloak"
	redisAdapter "github.com/hugosrc/shortlink/internal/adapter/redis"
	"github.com/hugosrc/shortlink/internal/adapter/zookeeper"
	"github.com/hugosrc/shortlink/internal/core/service"
	"github.com/hugosrc/shortlink/internal/handler/rest"
	"github.com/jxskiss/base62"
)

func main() {
	c, err := run(fmt.Sprintf(":%d", 3000))
	if err != nil {
		log.Fatalf("couldn't start the server: %s", err.Error())
	}

	if err := <-c; err != nil {
		log.Fatalf("an error occurred during execution: %s", err.Error())
	}
}

type serverConf struct {
	Address   string
	Auth      *keycloak.OpenIDAuth
	DB        *gocql.Session
	RDB       *redis.Client
	Zookeeper *zk.Conn
}

func newServer(conf serverConf) (*http.Server, error) {
	r := mux.NewRouter()

	counter := zookeeper.NewCounter(conf.Zookeeper)
	if err := counter.UpdateCounterBase(); err != nil {
		return nil, err
	}

	encoder := base62.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	caching := redisAdapter.NewRedisCaching(conf.RDB)
	repo := repository.NewLinkRepository(conf.DB)

	svc := service.NewLinkService(counter, encoder, caching, repo)

	rest.NewLinkHandler(conf.Auth, svc).Register(r)

	return &http.Server{
		Addr:              conf.Address,
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}, nil
}

func run(addr string) (<-chan error, error) {
	conf := config.Init()

	dbSession, err := cassandra.New(conf)
	if err != nil {
		return nil, err
	}

	rdbConn, err := redisAdapter.New(conf)
	if err != nil {
		return nil, err
	}

	zkConn, err := zookeeper.New(conf)
	if err != nil {
		return nil, err
	}

	auth := keycloak.NewOpenIDAuth(conf)

	srv, err := newServer(serverConf{
		Address:   addr,
		Auth:      auth,
		DB:        dbSession,
		RDB:       rdbConn,
		Zookeeper: zkConn,
	})
	if err != nil {
		return nil, err
	}

	c := make(chan error, 1)
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		log.Println("shutdown signal received")

		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			dbSession.Close()
			rdbConn.Close()
			zkConn.Close()
			stop()
			cancel()
			close(c)
		}()

		srv.SetKeepAlivesEnabled(false)

		if err := srv.Shutdown(timeout); err != nil {
			c <- err
		}

		log.Println("successful shutdown")
	}()

	go func() {
		log.Println("server started")

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			c <- err
		}
	}()

	return c, nil
}
