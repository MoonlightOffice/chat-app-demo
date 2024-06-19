package pubsub

import (
	"chat-app-demo/repository/pubsub/redis"
)

var PubSub IPubSub = redis.NewRedisPubSub()
