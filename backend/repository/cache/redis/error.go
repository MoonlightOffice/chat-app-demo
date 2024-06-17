package redis

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

var (
	ErrNotExist = errors.New("record does not exist")
)

func isNotRedisNil(err error) bool {
	return err != redis.Nil
}
