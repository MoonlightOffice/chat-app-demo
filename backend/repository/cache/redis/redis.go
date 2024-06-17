package redis

import (
	"chat-app-demo/config"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb *redis.ClusterClient
}

func NewRedisCache() *RedisCache {
	return &RedisCache{rdb: getClient()}
}

func getClient() *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: config.AppConfig.RedisCluster,
	})
}
