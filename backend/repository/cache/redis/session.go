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

func keySession(userId, sessionId string) string {
	return fmt.Sprintf("us.%s.%s", userId, sessionId)
}

func (rc *RedisCache) FindSession(userId, sessionId string) (*domain.Session, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	key := keySession(userId, sessionId)

	// Fetch from cache
	sObjJson, err := rc.rdb.Get(ctx, key).Result()
	if err == nil {
		var sObj domain.Session
		err = json.Unmarshal([]byte(sObjJson), &sObj)
		if err != nil {
			return nil, lib.ErrBuilder(err)
		}

		return &sObj, nil
	} else if isNotRedisNil(err) {
		return nil, lib.ErrBuilder(err)
	}

	// Fetch from db
	dbrepo, err := db.NewDBRepository()
	if err != nil {
		return nil, lib.ErrBuilder(err)
	}
	defer dbrepo.Close()

	sObj, err := dbrepo.FindSession(userId, sessionId)
	if err != nil {
		if errors.Is(err, db.ErrNotExist) {
			return nil, ErrNotExist
		}

		return nil, lib.ErrBuilder(err)
	}

	// Set to Redis cache
	b, _ := json.Marshal(sObj)
	_, err = rc.rdb.Set(ctx, key, b, 0).Result()
	if err != nil {
		return nil, lib.ErrBuilder(err)
	}

	return sObj, nil
}

func (rc *RedisCache) DeleteSession(userId, sessionId string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	key := keySession(userId, sessionId)
	_, err := rc.rdb.Del(ctx, key).Result()
	if err != nil {
		return lib.ErrBuilder(err)
	}

	return nil
}
