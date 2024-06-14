package tidb

import (
	"errors"
	"testing"
	"time"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
)

func TestMessage(t *testing.T) {
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

	// Add a message
	msgObj := domain.Message{
		RoomId:    rObj.RoomId,
		MessageId: 1,
		UserId:    uObj.UserId,
		Content:   "First message",
		SentAt:    time.Now(),
	}

	err = dbrepo.AddMessage(msgObj)
	if err != nil {
		t.Fatal(err)
	}

	err = dbrepo.AddMessage(msgObj)
	if !errors.Is(err, ErrDuplicated) {
		t.Fatal("Expected ErrDuplicated")
	}

	// Get message
	fmsgObj, err := dbrepo.FindMessage(msgObj.RoomId, msgObj.MessageId)
	if err != nil {
		t.Fatal(err)
	}
	if !lib.CompareStructs(*fmsgObj, msgObj) {
		t.Fatal("message doesn't match")
	}

	// Delete the message
	err = dbrepo.DeleteMessage(msgObj.RoomId, msgObj.MessageId)
	if err != nil {
		t.Fatal(err)
	}

	_, err = dbrepo.FindMessage(msgObj.RoomId, msgObj.MessageId)
	if !errors.Is(err, ErrNotExist) {
		t.Fatal("Expected ErrNotExist", err)
	}
}
