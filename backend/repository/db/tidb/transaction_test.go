package tidb

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
)

func TestTransaction(t *testing.T) {
	defer DeleteAll()

	repo, err := NewTiDBRepository()
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()

	tx, err := repo.BeginTx()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.RollbackTx()

	// Prepare mock data
	uObj, _ := domain.NewUser("u1")

	// Check manual rollback
	err = tx.AddUser(*uObj)
	if err != nil {
		t.Fatal(err)
	}

	err = tx.RollbackTx() // check this rollback reverts the last AddUser()
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.FindUser(uObj.UserId)
	if !errors.Is(err, ErrNotExist) {
		t.Fatal("Expected user not found")
	}

	// Check commit
	tx, err = repo.BeginTx()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.RollbackTx()

	err = tx.AddUser(*uObj)
	if err != nil {
		t.Fatal(err)
	}

	err = tx.CommitTx()
	if err != nil {
		t.Fatal(err)
	}

	err = tx.RollbackTx() // Check this rollback doesn't revert the last commit
	if err == nil {
		t.Fatal("expected error: already committed")
	}

	fetchedUserObj, err := repo.FindUser(uObj.UserId)
	if err != nil {
		t.Fatal(err)
	}

	ok := lib.CompareStructs(*uObj, *fetchedUserObj)
	if !ok {
		t.Fatal("user objects don't match")
	}
}

func TestCutInLine(t *testing.T) {
	defer DeleteAll()

	// Prepare mock data
	uObj, _ := domain.NewUser("u1")

	repo, err := NewTiDBRepository()
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()

	err = repo.AddUser(*uObj)
	if err != nil {
		t.Fatal(err)
	}

	lObj := domain.NewLogin(uObj.UserId)
	err = repo.AddLogin(*lObj)
	if err != nil {
		t.Fatal(err)
	}

	// Check cut in line

	tx, err := repo.BeginTx()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.RollbackTx()

	flObj, err := tx.FindLogin(lObj.UserId)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan error, 1)
	go func(ch chan error) {
		repo2, err := NewTiDBRepository()
		if err != nil {
			ch <- err
		}
		tx2, err := repo2.BeginTx()
		if err != nil {
			ch <- err
		}
		defer tx2.RollbackTx()

		flObjCutIn, err := tx2.FindLogin(lObj.UserId)
		if err != nil {
			ch <- err
		}
		if flObjCutIn.CreatedAt.UnixMilli() != 0 {
			ch <- fmt.Errorf("UnixMilli: %d", flObjCutIn.CreatedAt.UnixMilli())
		}

		ch <- nil
	}(ch)

	flObj.CreatedAt = time.UnixMilli(0)
	err = tx.UpdateLogin(*flObj)
	if err != nil {
		t.Fatal(err)
	}

	err = tx.CommitTx()
	if err != nil {
		t.Fatal(err)
	}

	flObj2, err := repo.FindLogin(lObj.UserId)
	if err != nil {
		t.Fatal(err)
	}
	if flObj2.CreatedAt.UnixMilli() != 0 {
		t.Fatal("Expected 0")
	}

	err = <-ch
	if err != nil {
		t.Fatal(err)
	}
}
