package tidb

import (
	"errors"
	"testing"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
)

func TestSession(t *testing.T) {
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

	// Add a session
	sObj := domain.NewSession(uObj.UserId)

	err = dbrepo.AddSession(*sObj)
	if err != nil {
		t.Fatal(err)
	}

	err = dbrepo.AddSession(*sObj)
	if !errors.Is(err, ErrDuplicated) {
		t.Fatal("Expected ErrDuplicated")
	}

	// Get message
	fsObj, err := dbrepo.FindSession(sObj.UserId, sObj.SessionId)
	if err != nil {
		t.Fatal(err)
	}
	if !lib.CompareStructs(*fsObj, *sObj) {
		t.Fatal("message doesn't match")
	}

	// Delete the message
	err = dbrepo.DeleteSession(sObj.UserId, sObj.SessionId)
	if err != nil {
		t.Fatal(err)
	}

	_, err = dbrepo.FindSession(sObj.UserId, sObj.SessionId)
	if !errors.Is(err, ErrNotExist) {
		t.Fatal("Expected ErrNotExist", err)
	}
}
