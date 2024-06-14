package domain

import (
	"testing"
	"time"

	"chat-app-demo/lib"
)

func TestNewSession(t *testing.T) {
	now := time.Now().UTC()
	userId := "u-1"

	s := NewSession(userId)

	if s.UserId != userId {
		t.Fatalf("Expected user id = %s, got %s\n", userId, s.UserId)
	}

	if s.IssuedAt.Sub(now) < 0 {
		t.Fatalf("Invalid CreatedAt: %v\n", s.IssuedAt)
	}

	if len(s.SessionCode) == 0 {
		t.Fatalf("Invalid session code: %s\n", s.SessionCode)
	}
}

func Test_SessionToken(t *testing.T) {
	s1 := NewSession("u-1")

	token := s1.ToToken()

	s2, ok := GetSessionFromToken(token)
	if !ok {
		t.Fatal("Invalid session token")
	}

	if !lib.CompareStructs(*s1, s2) {
		t.Fatal("s1 and s2 doesn't match")
	}
}

func TestSession_IsExpired(t *testing.T) {
	s := NewSession("u-1")
	s.IssuedAt = time.Now().Add(-sessionLifespan)

	if !s.IsExpired() {
		t.Fatalf("Expected session expired: %v\n", s.IssuedAt)
	}
}
