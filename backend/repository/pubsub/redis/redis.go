package redis

import (
	"chat-app-demo/config"

	"github.com/redis/go-redis/v9"
)

type RedisPubSub struct {
	rdb *redis.ClusterClient
}

func NewRedisPubSub() *RedisPubSub {
	return &RedisPubSub{rdb: getClient()}
}

func getClient() *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: config.AppConfig.RedisCluster,
	})
}
