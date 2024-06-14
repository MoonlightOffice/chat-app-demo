package tidb

import (
	"context"
	"time"

	"chat-app-demo/domain"
	"chat-app-demo/lib"
)

func (repo TiDBRepository) AddParticipant(rpObj domain.Participant) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	stmt := `INSERT INTO participants (room_id, user_id, joined_at) VALUES (?, ?, ?)`
	_, err := repo.crud.ExecContext(
		ctx,
		stmt,
		rpObj.RoomId,
		rpObj.UserId,
		rpObj.JoinedAt.UnixMilli(),
	)
	if err != nil {
		if isErrDuplicate(err) {
			return lib.ErrBuilder(ErrDuplicated)
		}

		return lib.ErrBuilder(err)
	}

	return nil
}

func (repo TiDBRepository) FindParticipants(roomId string) ([]domain.Participant, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	stmt := `SELECT user_id, joined_at FROM participants WHERE room_id = ? FOR UPDATE`
	rows, err := repo.crud.QueryContext(ctx, stmt, roomId)
	if err != nil {
		return nil, lib.ErrBuilder(err)
	}
	defer rows.Close()

	participants := make([]domain.Participant, 0)
	for rows.Next() {
		var (
			userId   string
			joinedAt int64
		)
		err = rows.Scan(&userId, &joinedAt)
		if err != nil {
			return nil, lib.ErrBuilder(err)
		}

		participants = append(participants, domain.Participant{
			RoomId:   roomId,
			UserId:   userId,
			JoinedAt: time.UnixMilli(joinedAt),
		})
	}

	return participants, nil
}

func (repo TiDBRepository) DeleteParticipant(roomId, userId string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	stmt := `DELETE FROM participants WHERE room_id = ? AND user_id = ?`
	_, err := repo.crud.ExecContext(
		ctx,
		stmt,
		roomId,
		userId,
	)
	if err != nil {
		return lib.ErrBuilder(err)
	}

	return nil
}
