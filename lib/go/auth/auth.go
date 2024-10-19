package auth

import (
	"context"
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
}

type PasswordCredentials struct {
	Username string
	Password string
}
