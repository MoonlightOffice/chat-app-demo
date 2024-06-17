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

func keyMessage(roomId string, msgId int64) string {
	return fmt.Sprintf("m.%s.%d", roomId, msgId)
}

func (rc *RedisCache) FindMessage(roomId string, msgId int64) (*domain.Message, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	key := keyMessage(roomId, msgId)

	// Fetch from cache
	mObjJson, err := rc.rdb.Get(ctx, key).Result()
	if err == nil {
		var mObj domain.Message
		err = json.Unmarshal([]byte(mObjJson), &mObj)
		if err != nil {
			return nil, lib.ErrBuilder(err)
		}

		return &mObj, nil
	} else if isNotRedisNil(err) {
		return nil, lib.ErrBuilder(err)
	}

	// Fetch from db
	dbrepo, err := db.NewDBRepository()
	if err != nil {
		return nil, lib.ErrBuilder(err)
	}
	defer dbrepo.Close()

	mObj, err := dbrepo.FindMessage(roomId, msgId)
	if err != nil {
		if errors.Is(err, db.ErrNotExist) {
			return nil, ErrNotExist
		}

		return nil, lib.ErrBuilder(err)
	}

	// Set to Redis cache
	b, _ := json.Marshal(mObj)
	_, err = rc.rdb.Set(ctx, key, b, 0).Result()
	if err != nil {
		return nil, lib.ErrBuilder(err)
	}

	return mObj, nil
}

func (rc *RedisCache) DeleteMessage(roomId string, msgId int64) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	key := keyMessage(roomId, msgId)
	_, err := rc.rdb.Del(ctx, key).Result()
	if err != nil {
		return lib.ErrBuilder(err)
	}

	return nil
}
