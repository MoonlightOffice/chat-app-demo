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

func TestParticipants(t *testing.T) {
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

	rpObj := domain.Participant{
		RoomId:   rObj.RoomId,
		UserId:   uObj.UserId,
		JoinedAt: time.Now(),
	}
	err = dbrepo.AddParticipant(rpObj)
	if err != nil {
		t.Fatal(err)
	}

	// Get data
	frpObjs, err := rc.FindParticipants(rObj.RoomId)
	if err != nil {
		t.Fatal(err)
	}
	if len(frpObjs) != 1 {
		t.Fatal("Expected 1, got", len(frpObjs))
	}
	if !lib.CompareStructs(rpObj, frpObjs[0]) {
		t.Fatal("doesn't match")
	}

	// Confirm it's cached
	_, err = rc.rdb.Get(context.Background(), keyParticipants(rObj.RoomId)).Result()
	if err != nil {
		t.Fatal(err)
	}

	// Delete cache
	err = rc.DeleteParticipants(rObj.RoomId)
	if err != nil {
		t.Fatal(err)
	}

	// Confirm cached data is deleted
	_, err = rc.rdb.Get(context.Background(), keyParticipants(rObj.RoomId)).Result()
	if err != redis.Nil {
		t.Fatal("Expected redis.Nil, got:", err)
	}
}
