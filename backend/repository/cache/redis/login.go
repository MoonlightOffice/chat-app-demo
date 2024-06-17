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

func keyLogin(userId string) string {
	return fmt.Sprintf("l.%s", userId)
}

func (rc *RedisCache) FindLogin(userId string) (*domain.Login, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	key := keyLogin(userId)

	// Fetch from cache
	lObjJson, err := rc.rdb.Get(ctx, key).Result()
	if err == nil {
		var lObj domain.Login
		err = json.Unmarshal([]byte(lObjJson), &lObj)
		if err != nil {
			return nil, lib.ErrBuilder(err)
		}

		return &lObj, nil
	} else if isNotRedisNil(err) {
		return nil, lib.ErrBuilder(err)
	}

	// Fetch from db
	dbrepo, err := db.NewDBRepository()
	if err != nil {
		return nil, lib.ErrBuilder(err)
	}
	defer dbrepo.Close()

	lObj, err := dbrepo.FindLogin(userId)
	if err != nil {
		if errors.Is(err, db.ErrNotExist) {
			return nil, ErrNotExist
		}

		return nil, lib.ErrBuilder(err)
	}

	// Set to Redis cache
	b, _ := json.Marshal(lObj)
	_, err = rc.rdb.Set(ctx, key, b, 0).Result()
	if err != nil {
		return nil, lib.ErrBuilder(err)
	}

	return lObj, nil
}

func (rc *RedisCache) DeleteLogin(userId string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	key := keyLogin(userId)
	_, err := rc.rdb.Del(ctx, key).Result()
	if err != nil {
		return lib.ErrBuilder(err)
	}

	return nil
}
