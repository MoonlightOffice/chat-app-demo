package redis

import (
	"context"

	"chat-app-demo/config"
)

func DeleteAll() error {
	// Disable this function in production
	if config.AppConfig.Mode == config.ModeProduction {
		return nil
	}

	rdb := getClient()
	_, err := rdb.FlushAll(context.TODO()).Result()

	return err
}
