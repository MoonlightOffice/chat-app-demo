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

func keyRoom(roomId string) string {
	return fmt.Sprintf("r.%s", roomId)
}

func (rc *RedisCache) FindRoom(roomId string) (*domain.Room, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	key := keyRoom(roomId)

	// Fetch from cache
	rObjJson, err := rc.rdb.Get(ctx, key).Result()
	if err == nil {
		var rObj domain.Room
		err = json.Unmarshal([]byte(rObjJson), &rObj)
		if err != nil {
			return nil, lib.ErrBuilder(err)
		}

		return &rObj, nil
	} else if isNotRedisNil(err) {
		return nil, lib.ErrBuilder(err)
	}

	// Fetch from db
	dbrepo, err := db.NewDBRepository()
	if err != nil {
		return nil, lib.ErrBuilder(err)
	}
	defer dbrepo.Close()

	rObj, err := dbrepo.FindRoom(roomId)
	if err != nil {
		if errors.Is(err, db.ErrNotExist) {
			return nil, ErrNotExist
		}

		return nil, lib.ErrBuilder(err)
	}

	// Set to Redis cache
	b, _ := json.Marshal(rObj)
	_, err = rc.rdb.Set(ctx, key, b, 0).Result()
	if err != nil {
		return nil, lib.ErrBuilder(err)
	}

	return rObj, nil
}

func (rc *RedisCache) DeleteRoom(roomId string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	key := keyRoom(roomId)
	_, err := rc.rdb.Del(ctx, key).Result()
	if err != nil {
		return lib.ErrBuilder(err)
	}

	return nil
}
