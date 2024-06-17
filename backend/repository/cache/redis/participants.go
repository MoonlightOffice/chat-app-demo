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

func keyParticipants(roomId string) string {
	return fmt.Sprintf("rp.%s", roomId)
}

func (rc *RedisCache) FindParticipants(roomId string) ([]domain.Participant, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	key := keyParticipants(roomId)

	// Fetch from cache
	rpObjsJson, err := rc.rdb.Get(ctx, key).Result()
	if err == nil {
		var rpObjs []domain.Participant
		err = json.Unmarshal([]byte(rpObjsJson), &rpObjs)
		if err != nil {
			return nil, lib.ErrBuilder(err)
		}

		return rpObjs, nil
	} else if isNotRedisNil(err) {
		return nil, lib.ErrBuilder(err)
	}

	// Fetch from db
	dbrepo, err := db.NewDBRepository()
	if err != nil {
		return nil, lib.ErrBuilder(err)
	}
	defer dbrepo.Close()

	rpObjs, err := dbrepo.FindParticipants(roomId)
	if err != nil {
		if errors.Is(err, db.ErrNotExist) {
			return nil, ErrNotExist
		}

		return nil, lib.ErrBuilder(err)
	}

	// Set to Redis cache
	b, _ := json.Marshal(rpObjs)
	_, err = rc.rdb.Set(ctx, key, b, 0).Result()
	if err != nil {
		return nil, lib.ErrBuilder(err)
	}

	return rpObjs, nil
}

func (rc *RedisCache) DeleteParticipants(roomId string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	key := keyParticipants(roomId)
	_, err := rc.rdb.Del(ctx, key).Result()
	if err != nil {
		return lib.ErrBuilder(err)
	}

	return nil
}
