package cache

import "chat-app-demo/domain"

type ICache interface {
	/* UserSession */

	// Handlable errors: ErrNotExist
	FindSession(userId, sessionId string) (*domain.Session, error)

	DeleteSession(userId, sessionId string) error

	/* Login */

	// Handlable errors: ErrNotExist
	FindLogin(userId string) (*domain.Login, error)

	DeleteLogin(userId string) error

	/* UserRoom */

	FindUserRooms(userId string) ([]domain.UserRoom, error)

	DeleteUserRooms(userId string) error

	/* Room */

	// Handlable errors: ErrNotExist
	FindRoom(roomId string) (*domain.Room, error)

	DeleteRoom(roomId string) error

	/* Participants */

	FindParticipants(roomId string) ([]domain.Participant, error)

	DeleteParticipants(roomId string) error

	/* Messages */

	// Handlable errors: ErrNotExist
	FindMessage(roomId string, msgId int64) (*domain.Message, error)

	DeleteMessage(roomId string, msgId int64) error
}
