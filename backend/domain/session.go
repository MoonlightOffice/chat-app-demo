package domain

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"chat-app-demo/lib"
)

const (
	sessionLifespan = 180 * 24 * time.Hour
)

type Session struct {
	UserId      string    `json:"userId"`
	SessionId   string    `json:"sessionId"`
	SessionCode string    `json:"sessionCode"`
	IssuedAt    time.Time `json:"issuedAt"`
}

func NewSession(userId string) *Session {
	now := time.Now().UTC()

	sessionId := lib.GenerateID("s")
	sessionCode := base64.RawURLEncoding.EncodeToString(lib.GenerateKey(32))

	return &Session{
		UserId:      userId,
		SessionId:   sessionId,
		SessionCode: sessionCode,
		IssuedAt:    now,
	}
}

type jwtPayload struct {
	UserId      string `json:"uid"`
	SessionId   string `json:"sid"`
	SessionCode string `json:"sec"`
	IssuedAt    int64  `json:"iat"`
}

func GetSessionFromToken(token string) (Session, bool) {
	payload, err := lib.FromJWT(token)
	if err != nil {
		return Session{}, false
	}

	var pl jwtPayload
	err = json.Unmarshal(payload, &pl)
	if err != nil {
		return Session{}, false
	}

	s := Session{
		UserId:      pl.UserId,
		SessionId:   pl.SessionId,
		SessionCode: pl.SessionCode,
		IssuedAt:    time.UnixMilli(pl.IssuedAt),
	}

	if len(s.UserId) == 0 {
		return Session{}, false
	}
	if len(s.SessionId) == 0 {
		return Session{}, false
	}
	if len(s.SessionCode) == 0 {
		return Session{}, false
	}

	return s, true
}

func (s Session) ToToken() string {
	payload, err := json.Marshal(jwtPayload{
		UserId:      s.UserId,
		SessionId:   s.SessionId,
		SessionCode: s.SessionCode,
		IssuedAt:    s.IssuedAt.UnixMilli(),
	})
	if err != nil {
		panic(fmt.Sprintf("Session.ToToken() has failed: %v\n", err))
	}

	token, err := lib.ToJWT(payload)
	if err != nil {
		panic(fmt.Sprintf("Session.ToToken() has failed: %v\n", err))
	}

	return token
}

func (s Session) IsExpired() bool {
	now := time.Now().UTC()

	return now.Sub(s.IssuedAt) > sessionLifespan
}
