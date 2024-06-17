package redis

import (
	"context"
	"testing"
	"time"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
	"chat-app-demo/repository/db"

	"github.com/redis/go-redis/v9"
)

func TestMessage(t *testing.T) {
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

	mObj := domain.Message{
		RoomId:    rObj.RoomId,
		MessageId: 1,
		UserId:    uObj.UserId,
		Content:   "some message",
		SentAt:    time.Now(),
	}
	err = dbrepo.AddMessage(mObj)
	if err != nil {
		t.Fatal(err)
	}

	// Get data
	fmObj, err := rc.FindMessage(mObj.RoomId, mObj.MessageId)
	if err != nil {
		t.Fatal(err)
	}
	if !lib.CompareStructs(mObj, *fmObj) {
		t.Fatal("doesn't match")
	}

	// Confirm it's cached
	_, err = rc.rdb.Get(context.Background(), keyMessage(mObj.RoomId, mObj.MessageId)).Result()
	if err != nil {
		t.Fatal(err)
	}

	// Delete cache
	err = rc.DeleteMessage(mObj.RoomId, mObj.MessageId)
	if err != nil {
		t.Fatal(err)
	}

	// Confirm cached data is deleted
	_, err = rc.rdb.Get(context.Background(), keyMessage(mObj.RoomId, mObj.MessageId)).Result()
	if err != redis.Nil {
		t.Fatal("Expected redis.Nil, got:", err)
	}
}
