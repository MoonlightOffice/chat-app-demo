package tidb

import (
	"errors"
	"testing"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
)

func TestRoom(t *testing.T) {
	defer DeleteAll()

	dbrepo, err := NewTiDBRepository()
	if err != nil {
		t.Fatal(err)
	}
	defer dbrepo.Close()

	// Add a new user
	uObj, _ := domain.NewUser("u1")
	err = dbrepo.AddUser(*uObj)
	if err != nil {
		t.Fatal(err)
	}

	// Add a new room
	rObj, _ := domain.NewRoom("some room")
	err = dbrepo.AddRoom(*rObj)
	if err != nil {
		t.Fatal(err)
	}

	// Find the room
	frObj, err := dbrepo.FindRoom(rObj.RoomId)
	if err != nil {
		t.Fatal(err)
	}
	if !lib.CompareStructs(*rObj, *frObj) {
		t.Fatal("Room object doesn't match")
	}

	// Update the room name
	rObj.Name = "Another room"
	err = dbrepo.UpdateRoom(*rObj)
	if err != nil {
		t.Fatal(err)
	}
	frObj, err = dbrepo.FindRoom(rObj.RoomId)
	if err != nil {
		t.Fatal(err)
	}
	if !lib.CompareStructs(*rObj, *frObj) {
		t.Fatal("Room object doesn't match")
	}
	if frObj.Name != "Another room" {
		t.Fatal("Room name is incorrect")
	}

	// Delete the room
	err = dbrepo.DeleteRoom(rObj.RoomId)
	if err != nil {
		t.Fatal(err)
	}
	_, err = dbrepo.FindRoom(rObj.RoomId)
	if !errors.Is(err, ErrNotExist) {
		t.Fatal(err)
	}
}
