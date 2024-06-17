package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
	"chat-app-demo/repository/db"
)

func keyUserRooms(userId string) string {
	return fmt.Sprintf("ur.%s", userId)
}

func (rc *RedisCache) FindUserRooms(userId string) ([]domain.UserRoom, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	key := keyUserRooms(userId)

	// Fetch from cache
	urObjsJson, err := rc.rdb.Get(ctx, key).Result()
	if err == nil {
		var urObjs []domain.UserRoom
		err = json.Unmarshal([]byte(urObjsJson), &urObjs)
		if err != nil {
			return nil, lib.ErrBuilder(err)
		}

		return urObjs, nil
	} else if isNotRedisNil(err) {
		return nil, lib.ErrBuilder(err)
	}

	// Fetch from db
	dbrepo, err := db.NewDBRepository()
	if err != nil {
		return nil, lib.ErrBuilder(err)
	}
	defer dbrepo.Close()

	urObjs, err := dbrepo.FindUserRooms(userId)
	if err != nil {
		if errors.Is(err, db.ErrNotExist) {
			return nil, ErrNotExist
		}

		return nil, lib.ErrBuilder(err)
	}

	// Set to Redis cache
	b, _ := json.Marshal(urObjs)
	_, err = rc.rdb.Set(ctx, key, b, 0).Result()
	if err != nil {
		return nil, lib.ErrBuilder(err)
	}

	return urObjs, nil
}

func (rc *RedisCache) DeleteUserRooms(userId string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	key := keyUserRooms(userId)
	_, err := rc.rdb.Del(ctx, key).Result()
	if err != nil {
		return lib.ErrBuilder(err)
	}

	return nil
}
