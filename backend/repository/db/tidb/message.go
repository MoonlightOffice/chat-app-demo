package tidb

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
)

func (repo TiDBRepository) AddMessage(msgObj domain.Message) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	stmt := `INSERT INTO messages (room_id, message_id, user_id, content, sent_at) VALUES (?, ?, ?, ?, ?)`
	_, err := repo.crud.ExecContext(
		ctx,
		stmt,
		msgObj.RoomId,
		msgObj.MessageId,
		msgObj.UserId,
		msgObj.Content,
		msgObj.SentAt.UnixMilli(),
	)
	if err != nil {
		if isErrDuplicate(err) {
			return lib.ErrBuilder(ErrDuplicated)
		}

		return lib.ErrBuilder(err)
	}

	return nil
}

func (repo TiDBRepository) FindMessage(roomId string, msgId int64) (*domain.Message, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	var (
		userId  string
		content string
		sentAt  int64
	)

	stmt := `SELECT user_id, content, sent_at FROM messages WHERE room_id = ? AND message_id = ? FOR UPDATE`
	err := repo.crud.QueryRowContext(ctx, stmt, roomId, msgId).Scan(
		&userId,
		&content,
		&sentAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, lib.ErrBuilder(ErrNotExist)
		}

		return nil, lib.ErrBuilder(err)
	}

	return &domain.Message{
		RoomId:    roomId,
		MessageId: msgId,
		UserId:    userId,
		Content:   content,
		SentAt:    time.UnixMilli(sentAt),
	}, nil
}

func (repo TiDBRepository) DeleteMessage(roomId string, msgId int64) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	stmt := `DELETE FROM messages WHERE room_id = ? AND message_id = ?`
	_, err := repo.crud.ExecContext(
		ctx,
		stmt,
		roomId,
		msgId,
	)
	if err != nil {
		return lib.ErrBuilder(err)
	}

	return nil
}
