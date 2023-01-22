package zookeeper

import (
	"errors"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/hugosrc/shortlink/internal/util"
	"github.com/spf13/viper"
)

var (
	zookeeperCounterPath = "/shortlink"
)

func New(conf viper.Viper) (*zk.Conn, error) {
	conn, _, err := zk.Connect([]string{conf.GetString("ZOOKEEPER_SERVER")}, time.Second*5)
	if err != nil {
		return nil, util.WrapErrorf(err, util.ErrCodeUnknown, "error connecting to zookeeper")
	}

	zookeeperCounterPath = conf.GetString("ZOOKEEPER_COUNTER_PATH")

	_, err = conn.Create(
		conf.GetString("ZOOKEEPER_COUNTER_PATH"),
		[]byte(conf.GetString("ZOOKEEPER_COUNTER_DEFAULT_VALUE")),
		0,
		zk.WorldACL(zk.PermAll),
	)

	if err != nil && !errors.Is(err, zk.ErrNodeExists) {
		return nil, util.WrapErrorf(err, util.ErrCodeUnknown, "error creating path in zookeeper")
	}

	return conn, nil
}
