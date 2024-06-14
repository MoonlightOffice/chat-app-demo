package tidb

import (
	"errors"

	"chat-app-demo/repository/db/types"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrNotExist   = types.ErrNotExist
	ErrDuplicated = types.ErrDuplicated
)

func isErrDuplicate(err error) bool {
	const ErrCodeDuplicate = 1062

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeDuplicate {
		return true
	}

	return false
}
