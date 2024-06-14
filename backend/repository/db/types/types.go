package types

import "chat-app-demo/domain"

type DBRepository interface {
	DBMethods
	BeginTx() (DBTxRepository, error)

	// Close (or release) db connection
	Close()
}

type DBTxRepository interface {
	DBMethods
	CommitTx() error
	RollbackTx() error
}

type DBMethods interface {
	/* User */

	// Handlable errors: ErrDuplicated
	AddUser(uObj domain.User) error

	// Handlable errors: ErrNotExist
	FindUser(userId string) (*domain.User, error)

	/* UserSession */

	// Handlable errors: ErrDuplicated
	AddSession(sObj domain.Session) error

	// Handlable errors: ErrNotExist
	FindSession(userId, sessionId string) (*domain.Session, error)

	DeleteSession(userId, sessionId string) error

	/* Login */

	// Handlable errors: ErrDuplicated
	AddLogin(lObj domain.Login) error

	// Handlable errors: ErrNotExist
	FindLogin(userId string) (*domain.Login, error)

	UpdateLogin(lObj domain.Login) error

	DeleteLogin(userId string) error

	/* UserRoom */

	// Handlable errors: ErrDuplicated
	AddUserRoom(ur domain.UserRoom) error

	// This method is not transaction-safe.
	FindUserRooms(userId string) ([]domain.UserRoom, error)

	UpdateUserRoom(ur domain.UserRoom) error

	/* Room */

	// Handlable errors: ErrDuplicated
	AddRoom(rObj domain.Room) error

	// Handlable errors: ErrNotExist
	FindRoom(roomId string) (*domain.Room, error)

	UpdateRoom(ur domain.Room) error

	DeleteRoom(roomId string) error

	/* Participants */

	// Handlable errors: ErrDuplicated
	AddParticipant(rpObj domain.Participant) error

	FindParticipants(roomId string) ([]domain.Participant, error)

	DeleteParticipant(roomId, userId string) error

	/* Messages */

	// Handlable errors: ErrDuplicated
	AddMessage(msgObj domain.Message) error

	// Handlable errors: ErrNotExist
	FindMessage(roomId string, msgId int64) (*domain.Message, error)

	DeleteMessage(roomId string, msgId int64) error
}
