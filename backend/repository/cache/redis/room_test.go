package redis

import (
	"context"
	"testing"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
	"chat-app-demo/repository/db"

	"github.com/redis/go-redis/v9"
)

func TestRoom(t *testing.T) {
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

	rObj, _ := domain.NewRoom("sample room")
	err = dbrepo.AddRoom(*rObj)
	if err != nil {
		t.Fatal(err)
	}

	// Get data
	rObj, err = rc.FindRoom(rObj.RoomId)
	if err != nil {
		t.Fatal(err)
	}
	if !lib.CompareStructs(*rObj, *rObj) {
		t.Fatal("doesn't match")
	}

	// Confirm it's cached
	_, err = rc.rdb.Get(context.Background(), keyRoom(rObj.RoomId)).Result()
	if err != nil {
		t.Fatal(err)
	}

	// Delete cache
	err = rc.DeleteRoom(rObj.RoomId)
	if err != nil {
		t.Fatal(err)
	}

	// Confirm cached data is deleted
	_, err = rc.rdb.Get(context.Background(), keyRoom(rObj.RoomId)).Result()
	if err != redis.Nil {
		t.Fatal("Expected redis.Nil, got:", err)
	}
}
