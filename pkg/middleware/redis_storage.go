package middleware

import (
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
)

type RedisStorage struct {
	client   *redis.Client
	ctx      context.Context
	limiters map[string]*rate.Limiter
}

func NewRedisStorage() *RedisStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Endereço do Redis
	})
	return &RedisStorage{
		client:   rdb,
		ctx:      context.Background(),
		limiters: make(map[string]*rate.Limiter),
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

func (rs *RedisStorage) GetLimiter(key string) *rate.Limiter {
	return rs.limiters[key]
}

func (rs *RedisStorage) SetLimiter(key string, limiter *rate.Limiter) {
	rs.limiters[key] = limiter
}
