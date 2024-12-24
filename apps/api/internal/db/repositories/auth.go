package repositories

import (
	"context"
	"fmt"
	"rango/api/internal/auth"
	"rango/api/internal/core"
	"rango/api/internal/db/generated"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type PGAuthRepository struct {
	queries *generated.Queries
	conn    *pgx.Conn
}

func NewPGAuthRepository(conn *pgx.Conn) *PGAuthRepository {
	queries := generated.New(conn)

	return &PGAuthRepository{
		queries: queries,
		conn:    conn,
	}
}

// RegisterUser implements auth.AuthRepository.
func (p *PGAuthRepository) RegisterUser(ctx context.Context, params auth.UserWithPassword) (core.User, error) {
	tx, err := p.conn.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return core.User{}, fmt.Errorf("Failed to start transaction for user registration: %w", err)
	}

	queries := p.queries.WithTx(tx)

	user, err := queries.AddUser(ctx, generated.AddUserParams{
		Username: params.Username,
		Password: params.Password,
	})

	if err != nil {
		return core.User{}, fmt.Errorf("Failed to create user: %w", err)
	}

	org, err := queries.CreateOrganization(ctx, generated.CreateOrganizationParams{
		ID: pgtype.UUID{
			Valid: true,
			Bytes: params.Org.ID,
		},
		Name: params.User.Org.Name,
	})

	if err != nil {
		return core.User{}, fmt.Errorf("Failed to create organization: %w", err)
	}

	_, err = queries.CreateMembership(ctx, generated.CreateMembershipParams{
		UserID: pgtype.UUID{
			Bytes: user.ID.Bytes,
			Valid: true,
		},
		OrgID: pgtype.UUID{
			Bytes: org.ID.Bytes,
			Valid: true,
		},
		IsDefault: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
	})

	if err != nil {
		return core.User{}, fmt.Errorf("Failed to create user membership: %w", err)
	}

	err = tx.Commit(ctx)

	if err != nil {
		return core.User{}, fmt.Errorf("Failed to commit changes: %w", err)
	}

	return core.User{
		ID:       user.ID.Bytes,
		Username: user.Username,
		Org: core.Organization{
			ID:   org.ID.Bytes,
			Name: org.Name,
		},
	}, nil
}

// AddUser implements auth.AuthRepository.
func (p *PGAuthRepository) AddUser(ctx context.Context, params auth.UserWithPassword) (core.User, error) {
	u, err := p.queries.AddUser(ctx, generated.AddUserParams{
		Username: params.Username,
		Password: params.Password,
	})

	if err != nil {
		return core.User{}, err
	}

	return core.User{
		ID:       u.ID.Bytes,
		Username: u.Username,
	}, nil
}

// GetSession implements auth.AuthRepository.
func (p *PGAuthRepository) GetSession(ctx context.Context, sessionID string) (auth.Session, error) {
	session, err := p.queries.GetSession(ctx, sessionID)

	if err != nil {
		return auth.Session{}, fmt.Errorf("Failed to get session: %w", err)
	}

	user, err := p.GetUserByUsername(ctx, session.Username)

	if err != nil {
		return auth.Session{}, fmt.Errorf("Failed to get user for session: %w", err)
	}

	return auth.Session{
		ID:           session.ID,
		User:         user.User,
		CreatedAt:    session.CreatedAt.Time,
		RevokedAt:    session.RevokedAt.Time,
		ExpiresAt:    session.ExpiresAt.Time,
		LastActiveAt: session.LastActiveAt.Time,
		IsRevoked:    session.RevokedAt.Valid,
	}, nil
}

// GetUserByUsername implements auth.AuthRepository.
func (p *PGAuthRepository) GetUserByUsername(ctx context.Context, username string) (auth.UserWithPassword, error) {
	user, err := p.queries.GetUserByUsername(ctx, username)

	if err != nil {
		return auth.UserWithPassword{}, fmt.Errorf("Failed to get user: %w", err)
	}

	org, err := p.queries.GetDefaultUserOrganization(ctx, user.ID)

	return auth.UserWithPassword{
		User: core.User{
			Username: user.Username,
			ID:       core.IDType(user.ID.Bytes),
			Org: core.Organization{
				ID:   org.ID.Bytes,
				Name: org.Name,
			},
		},
		Password: user.Password,
	}, nil
}

func (p *PGAuthRepository) CreateSession(ctx context.Context, session auth.Session) error {
	err := p.queries.CreateSession(ctx, generated.CreateSessionParams{
		ID: session.ID,
		UserID: pgtype.UUID{
			Bytes: session.User.ID,
			Valid: true,
		},
		ExpiresAt: pgtype.Timestamptz{
			Time:  session.ExpiresAt,
			Valid: true,
		},
	})

	return err
}
