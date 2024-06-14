package tidb

import (
	"context"
	"fmt"
	"time"

	"chat-app-demo/config"
)

func serializeSlice(slice []string) string {
	if len(slice) == 0 {
		return "('')"
	}

	serialized := "("

	for i, item := range slice {
		serialized += fmt.Sprintf("'%s'", item)

		if i != len(slice)-1 {
			serialized += ","
		}
	}

	serialized += ")"

	return serialized
}

func DeleteAll() error {
	// Disable this function in production
	if config.AppConfig.Mode == config.ModeProduction {
		return nil
	}

	// Prepare DB connection
	dbrepo, err := NewTiDBRepository()
	if err != nil {
		return fmt.Errorf("crdb.TiDBRepository.DeleteAll(): %w", err)
	}
	dbrepo.Close() // Do not defer

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	// Define tables to delete data from
	tables := []string{
		"logins",
		"messages",
		"participants",
		"user_rooms",
		"user_sessions",
		"rooms",
		"users",
	}

	for _, table := range tables {
		stmt := fmt.Sprintf(`DELETE FROM %s WHERE true`, table)
		_, err := pool.ExecContext(ctx, stmt)
		if err != nil {
			return err
		}
	}

	return nil
}
