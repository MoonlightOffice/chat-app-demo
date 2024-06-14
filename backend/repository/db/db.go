package db

import (
	"chat-app-demo/lib"
	"chat-app-demo/repository/db/tidb"
	"chat-app-demo/repository/db/types"
)

// Do not forget to defer .Close()
func NewDBRepository() (types.DBRepository, error) {
	repo, err := tidb.NewTiDBRepository()
	if err != nil {
		return nil, lib.ErrBuilder(err)
	}

	return repo, nil
}
