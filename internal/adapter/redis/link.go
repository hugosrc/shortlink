package redis

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v9"
)

type RedisCaching struct {
	rdb *redis.Client
}

func NewRedisCaching(rdb *redis.Client) *RedisCaching {
	return &RedisCaching{
		rdb: rdb,
	}
}

func (c *RedisCaching) Get(ctx context.Context, hash string) (string, error) {
	url, err := c.rdb.Get(ctx, hash).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}

	return url, nil
}

func (c *RedisCaching) Set(ctx context.Context, hash string, originalURL string) error {
	if err := c.rdb.Set(ctx, hash, originalURL, 0).Err(); err != nil {
		return err
	}

	return nil
}

func (c *RedisCaching) Del(ctx context.Context, hash string) error {
	if err := c.rdb.Del(ctx, hash).Err(); err != nil {
		return err
	}

	return nil
}
