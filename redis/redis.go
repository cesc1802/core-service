package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cesc1802/core-service/config"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	Rdb       *redis.Client
	CommonDur time.Duration
}

func NewRedis(c config.Config) (*Redis, error) {
	commonDur, err := time.ParseDuration(c.Redis.Ttl.Common)
	if err != nil {
		return nil, err
	}
	return &Redis{
		Rdb: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port),
			Password: c.Redis.Password,
			DB:       c.Redis.Database,
		}),
		CommonDur: commonDur,
	}, nil
}

func (r *Redis) Close() error {
	if err := r.Rdb.Close(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) GetStruct(ctx context.Context, key string, val interface{}) error {
	result, err := r.Rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	if len(result) > 0 {
		err = json.Unmarshal([]byte(result), val)
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func (r *Redis) SetStruct(ctx context.Context, key string, val interface{}, exp time.Duration) error {
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = r.Rdb.Set(ctx, key, string(jsonBytes), exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) ScanStruct(ctx context.Context, key string, val interface{}) error {
	keys, _, err := r.Rdb.Scan(ctx, 0, key, 0).Result()
	if err != nil {
		return err
	}
	res := make([]interface{}, 0)
	if len(keys) > 0 {
		for _, key := range keys {
			err = r.GetStruct(ctx, key, val)
			if err != nil {
				return err
			}

			res = append(res, val)
		}
	}
	return nil
}
