package tidb

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
)

func (repo TiDBRepository) AddRoom(rObj domain.Room) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	lastMessageId := int64(0)
	if rObj.LastMessage != nil {
		lastMessageId = rObj.LastMessage.MessageId
	}

	stmt := `INSERT INTO rooms (room_id, created_at, name, last_message_id) VALUES (?, ?, ?, ?)`
	_, err := repo.crud.ExecContext(
		ctx,
		stmt,
		rObj.RoomId,
		rObj.CreatedAt.UnixMilli(),
		rObj.Name,
		lastMessageId,
	)
	if err != nil {
		if isErrDuplicate(err) {
			return lib.ErrBuilder(ErrDuplicated)
		}

		return lib.ErrBuilder(err)
	}

	return nil
}

func (repo TiDBRepository) FindRoom(roomId string) (*domain.Room, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	tx, err := repo.crud.Begin()
	if err != nil {
		return nil, lib.ErrBuilder(err)
	}
	defer tx.Rollback()

	// Fetch room
	var (
		createdAt     int64
		name          string
		lastMessageId int64
	)
	stmt := `SELECT created_at, name, last_message_id FROM rooms WHERE room_id = ? FOR UPDATE`
	err = tx.QueryRowContext(ctx, stmt, roomId).Scan(
		&createdAt,
		&name,
		&lastMessageId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, lib.ErrBuilder(ErrNotExist)
		}

		return nil, lib.ErrBuilder(err)
	}
	room := domain.Room{
		RoomId:    roomId,
		CreatedAt: time.UnixMilli(createdAt),
		Name:      name,
	}

	// Fetch last message
	var (
		userId  string
		content string
		sentAt  int64
	)
	stmt = `SELECT user_id, content, sent_at FROM messages WHERE room_id = ? AND message_id = ? FOR UPDATE`
	err = tx.QueryRowContext(ctx, stmt, roomId, lastMessageId).Scan(
		&userId,
		&content,
		&sentAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) && lastMessageId == 0 {
			// No last message means the room has just been created and nobody has
			// sent any messages.
			return &room, nil
		}

		return nil, lib.ErrBuilder(err)
	}
	*room.LastMessage = domain.Message{
		RoomId:    roomId,
		MessageId: lastMessageId,
		UserId:    userId,
		Content:   content,
		SentAt:    time.UnixMilli(sentAt),
	}

	return &room, nil
}

func (repo TiDBRepository) UpdateRoom(rObj domain.Room) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	lastMessageId := int64(0)
	if rObj.LastMessage != nil {
		lastMessageId = rObj.LastMessage.MessageId
	}

	stmt := `UPDATE rooms SET name = ?, last_message_id = ? WHERE room_id = ?`
	_, err := repo.crud.ExecContext(
		ctx,
		stmt,
		rObj.Name,
		lastMessageId,
		rObj.RoomId,
	)
	if err != nil {
		return lib.ErrBuilder(err)
	}

	return nil
}

func (repo TiDBRepository) DeleteRoom(roomId string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	tx, err := repo.crud.Begin()
	if err != nil {
		return lib.ErrBuilder(err)
	}
	defer tx.Rollback()

	// Fetch all participants
	stmt := `SELECT user_id FROM participants WHERE room_id = ? FOR UPDATE`
	rows, err := tx.QueryContext(ctx, stmt, roomId)
	if err != nil {
		return lib.ErrBuilder(err)
	}
	defer rows.Close()

	participantIDs := make([]string, 0)
	for rows.Next() {
		var userId string
		err = rows.Scan(&userId)
		if err != nil {
			return lib.ErrBuilder(err)
		}
		participantIDs = append(participantIDs, userId)
	}

	// table: user_rooms
	stmt = `DELETE FROM user_rooms WHERE user_id IN ` + serializeSlice(participantIDs)
	_, err = tx.ExecContext(ctx, stmt)
	if err != nil {
		return lib.ErrBuilder(err)
	}

	// table: participants
	stmt = `DELETE FROM participants WHERE room_id = ?`
	_, err = tx.ExecContext(ctx, stmt, roomId)
	if err != nil {
		return lib.ErrBuilder(err)
	}

	// table: messages
	affectedRows := int64(-1)
	for affectedRows != 0 {
		stmt = `DELETE FROM rooms WHERE room_id = ? LIMIT 1000`
		result, err := tx.ExecContext(ctx, stmt, roomId)
		if err != nil {
			return lib.ErrBuilder(err)
		}
		affectedRows, err = result.RowsAffected()
		if err != nil {
			return lib.ErrBuilder(err)
		}
	}

	// table: rooms
	stmt = `DELETE FROM rooms WHERE room_id = ?`
	_, err = tx.ExecContext(ctx, stmt, roomId)
	if err != nil {
		return lib.ErrBuilder(err)
	}

	err = tx.Commit()
	if err != nil {
		return lib.ErrBuilder(err)
	}

	return nil
}
