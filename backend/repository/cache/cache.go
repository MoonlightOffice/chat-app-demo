package cache

import (
	"chat-app-demo/repository/cache/redis"
)

var Cache ICache = redis.NewRedisCache()
