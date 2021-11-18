package sdkredis

import (
	"fmt"
	"github.com/cesc1802/core-service/config"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisOpt struct {
	name     string
	host     string
	port     string
	password string
}

type Redis struct {
	Rdb       *redis.Client
	CommonDur time.Duration
	*RedisOpt
}

func NewRedis(name string, cfg *config.RedisConfig) *Redis {
	return &Redis{
		RedisOpt: &RedisOpt{
			name:     name,
			host:     cfg.Host,
			port:     cfg.Port,
			password: cfg.Password,
		},
	}
}

func (r *Redis) Name() string {
	return r.name
}
func (r *Redis) Start() error {
	r.Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.host, r.port),
		Password: r.password,
	})
	return nil
}

func (r *Redis) Stop() error {
	if r.Rdb != nil {
		if err := r.Rdb.Close(); err != nil {
			return err
		}
	}
	return nil
}
