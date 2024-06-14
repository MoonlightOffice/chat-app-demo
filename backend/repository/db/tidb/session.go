package tidb

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
)

func (repo TiDBRepository) AddSession(sObj domain.Session) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	stmt := `INSERT INTO user_sessions (user_id, session_id, session_code, issued_at) VALUES (?, ?, ?, ?)`
	_, err := repo.crud.ExecContext(
		ctx,
		stmt,
		sObj.UserId,
		sObj.SessionId,
		sObj.SessionCode,
		sObj.IssuedAt.UnixMilli(),
	)
	if err != nil {
		if isErrDuplicate(err) {
			return lib.ErrBuilder(ErrDuplicated)
		}

		return lib.ErrBuilder(err)
	}

	return nil
}

func (repo TiDBRepository) FindSession(userId, sessionId string) (*domain.Session, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	var (
		sessionCode string
		issuedAt    int64
	)

	stmt := `SELECT session_code, issued_at FROM user_sessions WHERE user_id = ? AND session_id = ? FOR UPDATE`
	err := repo.crud.QueryRowContext(ctx, stmt, userId, sessionId).Scan(
		&sessionCode,
		&issuedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, lib.ErrBuilder(ErrNotExist)
		}

		return nil, lib.ErrBuilder(err)
	}

	return &domain.Session{
		UserId:      userId,
		SessionId:   sessionId,
		SessionCode: sessionCode,
		IssuedAt:    time.UnixMilli(issuedAt),
	}, nil
}

func (repo TiDBRepository) DeleteSession(userId, sessionId string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	stmt := `DELETE FROM user_sessions WHERE user_id = ? AND session_id = ?`
	_, err := repo.crud.ExecContext(
		ctx,
		stmt,
		userId,
		sessionId,
	)
	if err != nil {
		return lib.ErrBuilder(err)
	}

	return nil
}
