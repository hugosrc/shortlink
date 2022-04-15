package zookeeper

import (
	"strconv"
	"sync"

	"github.com/go-zookeeper/zk"
)

const (
	counterRange = 100000
)

type ZookeeperCounter struct {
	mu   sync.Mutex
	conn *zk.Conn
	base int
	v    int
}

func NewCounter(conn *zk.Conn) *ZookeeperCounter {
	return &ZookeeperCounter{
		conn: conn,
	}
}

func (c *ZookeeperCounter) Inc() (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inc := c.v + (c.base * counterRange)

	c.v += 1
	if c.v == counterRange {
		if err := c.UpdateCounterBase(); err != nil {
			return 0, err
		}
	}

	return inc, nil
}

func (c *ZookeeperCounter) UpdateCounterBase() error {
	bytes, _, err := c.conn.Get(zookeeperCounterPath)
	if err != nil {
		return err
	}

	base, err := strconv.Atoi(string(bytes))
	if err != nil {
		return err
	}

	c.v = 0
	c.base = base

	if _, err := c.conn.Set(zookeeperCounterPath, []byte(strconv.Itoa(base+1)), -1); err != nil {
		return err
	}

	return nil
}
