package domain

import (
	"chat-app-demo/lib"
	"time"
)

type Login struct {
	UserId    string    `json:"userId"`
	LoginCode string    `json:"loginCode"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewLogin(userId string) *Login {
	return &Login{
		UserId:    userId,
		LoginCode: lib.GenerateID("l"),
		CreatedAt: time.Now(),
	}
}
