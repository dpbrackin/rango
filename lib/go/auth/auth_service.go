package auth

import (
	"context"
	"fmt"
	"rango/core"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repository AuthRepository
	clock      core.Clock
}

type NewAuthServiceParams struct {
	Repository AuthRepository
	Clock      core.Clock
}

func NewAuthService(params NewAuthServiceParams) *AuthService {
	return &AuthService{
		repository: params.Repository,
		clock:      params.Clock,
	}
}

func (srv *AuthService) AuthenticateWithPassword(ctx context.Context, creds PasswordCredentials) (core.User, error) {
	var user core.User

	dbUser, err := srv.repository.GetUserByUsername(ctx, creds.Username)

	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(creds.Password))

	if err != nil {
		return user, err
	}

	user = core.User{
		Username: dbUser.Username,
	}

	return user, nil
}

func (srv *AuthService) Register(ctx context.Context, creds PasswordCredentials) (core.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)

	if err != nil {
		return core.User{}, err
	}

	err = srv.repository.AddUser(ctx, UserWithPassword{
		User: core.User{
			Username: creds.Username,
		},
		Password: string(hashedPassword),
	})

	if err != nil {
		return core.User{}, err
	}

	user := core.User{
		Username: creds.Username,
	}

	return user, nil
}

func (srv *AuthService) AuthenticateSession(ctx context.Context, sessionID string) (core.User, error) {
	session, err := srv.repository.GetSession(ctx, sessionID)

	if err != nil {
		return core.User{}, fmt.Errorf("Failed to get session: %w", err)
	}

	now := srv.clock.Now()

	isRevoked := session.IsRevoked && now.After(session.RevokedAt)
	isExpired := now.After(session.ExpiresAt)

	if isRevoked || isExpired {
		return core.User{}, fmt.Errorf("Session expired")
	}

	return session.User, nil
}
