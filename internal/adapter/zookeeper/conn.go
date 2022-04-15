package zookeeper

import (
	"errors"
	"time"

	"github.com/go-zookeeper/zk"
)

const (
	zookeeperDefaultData = "10000"
	zookeeperCounterPath = "/shorturl_seed"
	zookeeperServer      = "127.0.0.1:2181"
)

func New() (*zk.Conn, error) {
	conn, _, err := zk.Connect([]string{zookeeperServer}, time.Second*5)
	if err != nil {
		return nil, err
	}

	if _, err := conn.Create(
		zookeeperCounterPath,
		[]byte(zookeeperDefaultData),
		0,
		zk.WorldACL(zk.PermAll),
	); err != nil && !errors.Is(err, zk.ErrNodeExists) {
		return nil, err
	}

	return conn, nil
}
