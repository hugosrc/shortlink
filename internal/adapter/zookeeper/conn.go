package zookeeper

import (
	"errors"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/spf13/viper"
)

var (
	zookeeperCounterPath = "/shortlink"
)

func New(conf viper.Viper) (*zk.Conn, error) {
	conn, _, err := zk.Connect([]string{conf.GetString("ZOOKEEPER_SERVER")}, time.Second*5)
	if err != nil {
		return nil, err
	}

	zookeeperCounterPath = conf.GetString("ZOOKEEPER_COUNTER_PATH")

	if _, err := conn.Create(
		conf.GetString("ZOOKEEPER_COUNTER_PATH"),
		[]byte(conf.GetString("ZOOKEEPER_COUNTER_DEFAULT_VALUE")),
		0,
		zk.WorldACL(zk.PermAll),
	); err != nil && !errors.Is(err, zk.ErrNodeExists) {
		return nil, err
	}

	return conn, nil
}
