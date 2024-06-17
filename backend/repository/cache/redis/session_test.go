package redis

import (
	"context"
	"testing"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
	"chat-app-demo/repository/db"

	"github.com/redis/go-redis/v9"
)

func TestSession(t *testing.T) {
	defer DeleteAll()

	rc := NewRedisCache()

	// Prepare data
	dbrepo, err := db.NewDBRepository()
	if err != nil {
		t.Fatal(err)
	}
	defer dbrepo.Close()

	uObj, _ := domain.NewUser("u1")
	err = dbrepo.AddUser(*uObj)
	if err != nil {
		t.Fatal(err)
	}

	sObj := domain.NewSession(uObj.UserId)
	err = dbrepo.AddSession(*sObj)
	if err != nil {
		t.Fatal(err)
	}

	// Get data
	fsObj, err := rc.FindSession(sObj.UserId, sObj.SessionId)
	if err != nil {
		t.Fatal(err)
	}
	if !lib.CompareStructs(*sObj, *fsObj) {
		t.Fatal("doesn't match")
	}

	// Confirm it's cached
	_, err = rc.rdb.Get(context.Background(), keySession(sObj.UserId, sObj.SessionId)).Result()
	if err != nil {
		t.Fatal(err)
	}

	// Delete cache
	err = rc.DeleteSession(sObj.UserId, sObj.SessionId)
	if err != nil {
		t.Fatal(err)
	}

	// Confirm cached data is deleted
	_, err = rc.rdb.Get(context.Background(), keySession(sObj.UserId, sObj.SessionId)).Result()
	if err != redis.Nil {
		t.Fatal("Expected redis.Nil, got:", err)
	}
}
