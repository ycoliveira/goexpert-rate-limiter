package middleware

import (
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStorage() *RedisStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	return &RedisStorage{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (rs *RedisStorage) IsBlocked(key string) bool {
	val, err := rs.client.Get(rs.ctx, key).Result()
	if err == redis.Nil {
		return false
	}
	if err != nil {
		panic(err)
	}
	return val == "blocked"
}

func (rs *RedisStorage) Block(key string, duration time.Duration) {
	err := rs.client.Set(rs.ctx, key, "blocked", duration).Err()
	if err != nil {
		panic(err)
	}
}
