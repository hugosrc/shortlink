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

	"github.com/go-zookeeper/zk"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/hugosrc/shortlink/internal/adapter/cassandra"
	"github.com/hugosrc/shortlink/internal/adapter/cassandra/repository"
	"github.com/hugosrc/shortlink/internal/adapter/zookeeper"
	"github.com/hugosrc/shortlink/internal/core/service"
	"github.com/hugosrc/shortlink/internal/handler/rest"
	"github.com/jxskiss/base62"
)

func main() {
	c, err := run(fmt.Sprintf(":%d", 8000)) // TODO: Get Port From Env Variables
	if err != nil {
		log.Fatalf("couldn't start the server: %s", err.Error())
	}

	if err := <-c; err != nil {
		log.Fatalf("an error occurred during execution: %s", err.Error())
	}
}

type serverConf struct {
	Address   string
	DB        *gocql.Session
	Zookeeper *zk.Conn
}

func newServer(conf serverConf) (*http.Server, error) {
	r := mux.NewRouter()

	counter := zookeeper.NewCounter(conf.Zookeeper)
	if err := counter.UpdateCounterBase(); err != nil {
		return nil, err
	}

	encoder := base62.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	repo := repository.NewLinkRepository(conf.DB)

	svc := service.NewLinkService(counter, encoder, repo)

	rest.NewLinkHandler(svc).Register(r)

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
	dbSession, err := cassandra.New()
	if err != nil {
		return nil, err
	}

	zkConn, err := zookeeper.New()
	if err != nil {
		return nil, err
	}

	srv, err := newServer(serverConf{
		Address:   addr,
		DB:        dbSession,
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
