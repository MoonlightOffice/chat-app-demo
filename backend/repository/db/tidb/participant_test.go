package tidb

import (
	"errors"
	"testing"
	"time"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
)

func TestParticipant(t *testing.T) {
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

	// Add a participant
	rp := domain.Participant{
		RoomId:   rObj.RoomId,
		UserId:   uObj.UserId,
		JoinedAt: time.Now(),
	}

	err = dbrepo.AddParticipant(rp)
	if err != nil {
		t.Fatal(err)
	}

	err = dbrepo.AddParticipant(rp)
	if !errors.Is(err, ErrDuplicated) {
		t.Fatal("Expected ErrDuplicated")
	}

	// Get participants
	rps, err := dbrepo.FindParticipants(rObj.RoomId)
	if err != nil {
		t.Fatal(err)
	}
	if !lib.CompareStructs(rps[0], rp) {
		t.Fatal("Participant doesn't match")
	}

	// Delete the participant
	err = dbrepo.DeleteParticipant(rp.RoomId, rp.UserId)
	if err != nil {
		t.Fatal(err)
	}

	rps, err = dbrepo.FindParticipants(rObj.RoomId)
	if err != nil {
		t.Fatal(err)
	}
	if len(rps) != 0 {
		t.Fatal("Expected 0")
	}
}
