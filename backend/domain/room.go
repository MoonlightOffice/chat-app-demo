package domain

import (
	"chat-app-demo/lib"
	"time"
)

type Room struct {
	RoomId      string    `json:"roomId"`
	CreatedAt   time.Time `json:"createdAt"`
	Name        string    `json:"name"`
	LastMessage *Message  `json:"lastMessage"`
}

func NewRoom(name string) (*Room, bool) {
	n := len(name)
	if n == 0 || n > 300 {
		return nil, false
	}

	roomId := lib.GenerateID("r")
	now := time.Now()

	return &Room{
		RoomId:      roomId,
		CreatedAt:   now,
		Name:        name,
		LastMessage: nil,
	}, true
}

type Participant struct {
	RoomId   string    `json:"roomId"`
	UserId   string    `json:"userId"`
	JoinedAt time.Time `json:"joinedAt"`
}

type Message struct {
	RoomId    string    `json:"roomId"`
	MessageId int64     `json:"messageId"`
	UserId    string    `json:"userId"`
	Content   string    `json:"content"`
	SentAt    time.Time `json:"sentAt"`
}
