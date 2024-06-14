package tidb

import (
	"errors"
	"testing"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
)

func TestLogin(t *testing.T) {
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

	// Add a Login
	lObj := domain.NewLogin(uObj.UserId)

	err = dbrepo.AddLogin(*lObj)
	if err != nil {
		t.Fatal(err)
	}

	err = dbrepo.AddLogin(*lObj)
	if !errors.Is(err, ErrDuplicated) {
		t.Fatal("Expected ErrDuplicated")
	}

	// Get login
	flObj, err := dbrepo.FindLogin(lObj.UserId)
	if err != nil {
		t.Fatal(err)
	}
	if !lib.CompareStructs(*flObj, *lObj) {
		t.Fatal("Login doesn't match")
	}

	// Delete the login
	err = dbrepo.DeleteLogin(lObj.UserId)
	if err != nil {
		t.Fatal(err)
	}

	_, err = dbrepo.FindLogin(lObj.UserId)
	if !errors.Is(err, ErrNotExist) {
		t.Fatal("Expected ErrNotExist", err)
	}
}
