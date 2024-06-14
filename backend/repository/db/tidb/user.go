package tidb

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
)

func (repo TiDBRepository) AddUser(uObj domain.User) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	stmt := `INSERT INTO users (user_id, created_at) VALUES (?, ?)`
	_, err := repo.crud.ExecContext(
		ctx,
		stmt,
		uObj.UserId,
		uObj.CreatedAt.UnixMilli(),
	)
	if err != nil {
		if isErrDuplicate(err) {
			return lib.ErrBuilder(ErrDuplicated)
		}

		return lib.ErrBuilder(err)
	}

	return nil
}

func (repo TiDBRepository) FindUser(userId string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	var createdAt int64

	stmt := `SELECT created_at FROM users WHERE user_id = ? FOR UPDATE`
	err := repo.crud.QueryRowContext(ctx, stmt, userId).Scan(&createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, lib.ErrBuilder(ErrNotExist)
		}

		return nil, lib.ErrBuilder(err)
	}

	return &domain.User{
		UserId:    userId,
		CreatedAt: time.UnixMilli(createdAt),
	}, nil
}

func (repo TiDBRepository) AddUserRoom(ur domain.UserRoom) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	stmt := `INSERT INTO user_rooms (user_id, room_id, read_until) VALUES (?, ?, ?)`
	_, err := repo.crud.ExecContext(
		ctx,
		stmt,
		ur.UserId,
		ur.RoomId,
		ur.ReadUntil,
	)
	if err != nil {
		if isErrDuplicate(err) {
			return lib.ErrBuilder(ErrDuplicated)
		}

		return lib.ErrBuilder(err)
	}

	return nil
}

func (repo TiDBRepository) FindUserRooms(userId string) ([]domain.UserRoom, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	// Fetch room ids
	stmt := `SELECT room_id, read_until FROM user_rooms WHERE user_id = ? FOR UPDATE`
	rows, err := repo.crud.QueryContext(ctx, stmt, userId)
	if err != nil {
		return nil, lib.ErrBuilder(err)
	}
	defer rows.Close()

	urs := make([]domain.UserRoom, 0)
	for rows.Next() {
		ur := domain.UserRoom{UserId: userId}

		err := rows.Scan(&ur.RoomId, &ur.ReadUntil)
		if err != nil {
			return nil, lib.ErrBuilder(err)
		}

		urs = append(urs, ur)
	}

	// Fetch individual room data
	for i, ur := range urs {
		rObj, err := repo.FindRoom(ur.RoomId)
		if err != nil {
			return nil, lib.ErrBuilder(err)
		}

		ur.Room = *rObj
		urs[i] = ur
	}

	return urs, nil
}

func (repo TiDBRepository) UpdateUserRoom(ur domain.UserRoom) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	stmt := `UPDATE user_rooms SET room_id = ?, read_until = ? WHERE user_id = ?`
	_, err := repo.crud.ExecContext(
		ctx,
		stmt,
		ur.RoomId,
		ur.ReadUntil,
		ur.UserId,
	)
	if err != nil {
		return lib.ErrBuilder(err)
	}

	return nil
}
