package repositories

import (
	"context"
	"fmt"
	"rango/auth"
	"rango/core"
	"rango/db/generated"
)

type PGAuthRepository struct {
	queries generated.Queries
}

// AddUser implements auth.AuthRepository.
func (p *PGAuthRepository) AddUser(ctx context.Context, params auth.UserWithPassword) error {
	err := p.queries.AddUser(ctx, generated.AddUserParams{
		Username: params.Username,
		Password: params.Password,
	})

	return err
}

// GetSession implements auth.AuthRepository.
func (p *PGAuthRepository) GetSession(ctx context.Context, sessionID string) (auth.Session, error) {
	session, err := p.queries.GetSession(ctx, sessionID)

	if err != nil {
		return auth.Session{}, fmt.Errorf("Failed to get session: %w", err)
	}

	return auth.Session{
		ID: session.ID,
		User: core.User{
			Username: session.Username,
			ID:       core.IDType(session.UserID.Int32),
		},
		CreatedAt:    session.CreatedAt,
		RevokedAt:    session.RevokedAt,
		ExpiresAt:    session.ExpiresAt,
		LastActiveAt: session.LastActiveAt,
	}, nil
}

// GetUserByUsername implements auth.AuthRepository.
func (p *PGAuthRepository) GetUserByUsername(ctx context.Context, username string) (auth.UserWithPassword, error) {
	user, err := p.queries.GetUserByUsername(ctx, username)

	if err != nil {
		return auth.UserWithPassword{}, fmt.Errorf("Failed to get user: %w", err)
	}

	return auth.UserWithPassword{
		User: core.User{
			Username: user.Username,
			ID:       core.IDType(user.ID),
		},
		Password: user.Password,
	}, nil
}
