package redis

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/hugosrc/shortlink/internal/util"
	"github.com/spf13/viper"
)

func New(conf *viper.Viper) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.GetString("REDIS_SERVER"),
		Password: conf.GetString("REDIS_PASSWORD"),
		DB:       conf.GetInt("REDIS_DATABASE"),
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, util.WrapErrorf(err, util.ErrCodeUnknown, "error connecting to redis server")
	}

	return rdb, nil
}
