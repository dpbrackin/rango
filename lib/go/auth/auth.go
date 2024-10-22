package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"rango/core"
	"time"
)

type Session struct {
	ID           string
	User         core.User
	CreatedAt    time.Time
	RevokedAt    time.Time
	ExpiresAt    time.Time
	LastActiveAt time.Time
	IsRevoked    bool
}

type UserWithPassword struct {
	core.User
	Password string
}

type AuthRepository interface {
	GetUserByUsername(ctx context.Context, username string) (UserWithPassword, error)
	AddUser(ctx context.Context, params UserWithPassword) error
	GetSession(ctx context.Context, sessionID string) (Session, error)
	CreateSession(ctx context.Context, session Session) error
}

type PasswordCredentials struct {
	Username string
	Password string
}

func NewSessionID(n int) (string, error) {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		ret[i] = charset[num.Int64()]
	}

	return string(ret), nil
}

func NewSession(user core.User) (*Session, error) {
	sessionTTL := time.Hour * 24 * 30
	sessionID, err := NewSessionID(32)

	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %w", err)
	}
	return &Session{
		ID:           sessionID,
		User:         user,
		CreatedAt:    time.Time{},
		RevokedAt:    time.Time{},
		ExpiresAt:    time.Now().Add(sessionTTL),
		LastActiveAt: time.Time{},
		IsRevoked:    false,
	}, nil
}
