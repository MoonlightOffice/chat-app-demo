package redis

import (
	"context"
	"testing"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
	"chat-app-demo/repository/db"

	"github.com/redis/go-redis/v9"
)

func TestUserRooms(t *testing.T) {
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

	urObj := domain.UserRoom{
		UserId:    uObj.UserId,
		Room:      *rObj,
		ReadUntil: 3,
	}
	err = dbrepo.AddUserRoom(urObj)
	if err != nil {
		t.Fatal(err)
	}

	// Get data
	furObjs, err := rc.FindUserRooms(uObj.UserId)
	if err != nil {
		t.Fatal(err)
	}
	if len(furObjs) != 1 {
		t.Fatal("Expected 1, got", len(furObjs))
	}
	if !lib.CompareStructs(urObj, furObjs[0]) {
		t.Fatal("doesn't match")
	}

	// Confirm it's cached
	_, err = rc.rdb.Get(context.Background(), keyUserRooms(uObj.UserId)).Result()
	if err != nil {
		t.Fatal(err)
	}

	// Delete cache
	err = rc.DeleteUserRooms(uObj.UserId)
	if err != nil {
		t.Fatal(err)
	}

	// Confirm cached data is deleted
	_, err = rc.rdb.Get(context.Background(), keyUserRooms(uObj.UserId)).Result()
	if err != redis.Nil {
		t.Fatal("Expected redis.Nil, got:", err)
	}
}
