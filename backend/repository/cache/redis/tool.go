package redis

import (
	"context"

	"chat-app-demo/config"
	"chat-app-demo/repository/db/tidb"
)

func DeleteAll() error {
	err := tidb.DeleteAll()
	if err != nil {
		return err
	}

	// Disable this function in production
	if config.AppConfig.Mode == config.ModeProduction {
		return nil
	}

	rdb := getClient()
	_, err = rdb.FlushAll(context.TODO()).Result()

	return err
}
