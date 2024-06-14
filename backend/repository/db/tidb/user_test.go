package tidb

import (
	"errors"
	"fmt"
	"testing"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
)

func TestUser(t *testing.T) {
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

	// Check if duplicated error occurs
	err = dbrepo.AddUser(*uObj)
	if !errors.Is(err, ErrDuplicated) {
		t.Fatal("Expected ErrDuplicated")
	}

	// Fetch user from db and check
	fuObj, err := dbrepo.FindUser(uObj.UserId)
	if err != nil {
		t.Fatal(err)
	}

	if !lib.CompareStructs(*uObj, *fuObj) {
		t.Fatal("Fetched user doesn't match")
	}
}

func TestUserRoom(t *testing.T) {
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

	// Add a new user room
	ur := domain.UserRoom{
		UserId:    "u1",
		Room:      *rObj,
		ReadUntil: 5,
	}

	err = dbrepo.AddUserRoom(ur)
	if err != nil {
		t.Fatal(err)
	}

	// Check if duplicated error occurs
	err = dbrepo.AddUserRoom(ur)
	if !errors.Is(err, ErrDuplicated) {
		t.Fatal("Expected ErrDuplicated")
	}

	// Fetch user from db and check
	urs, err := dbrepo.FindUserRooms("u1")
	if err != nil {
		t.Fatal(err)
	}
	if len(urs) != 1 {
		t.Fatal("n of rooms got", len(urs))
	}

	if !lib.CompareStructs(urs[0], ur) {
		fmt.Println(urs[0].Room)
		fmt.Println(ur.Room)
		t.Fatal("Fetched user doesn't match")
	}
}
