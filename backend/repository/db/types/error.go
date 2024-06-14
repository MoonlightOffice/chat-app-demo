package types

import (
	"errors"
)

var (
	ErrNotExist   = errors.New("record does not exist")
	ErrDuplicated = errors.New("duplicated record")
)
