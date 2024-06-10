package domain

import "time"

type User struct {
	UserId    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewUser(id string) (*User, bool) {
	n := len(id)
	if n == 0 || n > 20 {
		return nil, false
	}

	return &User{UserId: id, CreatedAt: time.Now()}, true
}

type UserRoom struct {
	UserId    string `json:"userId"`
	Room      `json:"room"`
	ReadUntil uint64 `json:"readUntil"` // The last message's id the user has read in the room
}
