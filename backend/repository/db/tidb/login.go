package tidb

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
)

func (repo TiDBRepository) AddLogin(lObj domain.Login) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	stmt := `INSERT INTO logins (user_id, login_code, created_at) VALUES (?, ?, ?)`
	_, err := repo.crud.ExecContext(
		ctx,
		stmt,
		lObj.UserId,
		lObj.LoginCode,
		lObj.CreatedAt.UnixMilli(),
	)
	if err != nil {
		if isErrDuplicate(err) {
			return lib.ErrBuilder(ErrDuplicated)
		}

		return lib.ErrBuilder(err)
	}

	return nil
}

func (repo TiDBRepository) FindLogin(userId string) (*domain.Login, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	var (
		loginCode string
		createdAt int64
	)

	stmt := `SELECT login_code, created_at FROM logins WHERE user_id = ? FOR UPDATE`
	err := repo.crud.QueryRowContext(ctx, stmt, userId).Scan(
		&loginCode,
		&createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, lib.ErrBuilder(ErrNotExist)
		}

		return nil, lib.ErrBuilder(err)
	}

	return &domain.Login{
		UserId:    userId,
		LoginCode: loginCode,
		CreatedAt: time.UnixMilli(createdAt),
	}, nil
}

func (repo TiDBRepository) UpdateLogin(lObj domain.Login) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	stmt := `UPDATE logins SET login_code = ?, created_at = ? WHERE user_id = ?`
	_, err := repo.crud.ExecContext(
		ctx,
		stmt,
		lObj.LoginCode,
		lObj.CreatedAt.UnixMilli(),
		lObj.UserId,
	)
	if err != nil {
		return lib.ErrBuilder(err)
	}

	return nil
}

func (repo TiDBRepository) DeleteLogin(userId string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	stmt := `DELETE FROM logins WHERE user_id = ?`
	_, err := repo.crud.ExecContext(
		ctx,
		stmt,
		userId,
	)
	if err != nil {
		return lib.ErrBuilder(err)
	}

	return nil
}
