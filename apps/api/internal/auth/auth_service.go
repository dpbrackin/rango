package auth

import (
	"context"
	"fmt"
	"rango/api/internal/core"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepository AuthRepository
	orgRepository  core.OrgatizationRepository
	clock          core.Clock
}

type NewAuthServiceParams struct {
	AuthRepository AuthRepository
	OrgRepository  core.OrgatizationRepository
	Clock          core.Clock
}

func NewAuthService(params NewAuthServiceParams) *AuthService {
	return &AuthService{
		authRepository: params.AuthRepository,
		orgRepository:  params.OrgRepository,
		clock:          params.Clock,
	}
}

func (srv *AuthService) AuthenticateWithPassword(ctx context.Context, creds PasswordCredentials) (core.User, error) {
	var user core.User

	dbUser, err := srv.authRepository.GetUserByUsername(ctx, creds.Username)

	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(creds.Password))

	if err != nil {
		return user, err
	}

	user = core.User{
		Username: dbUser.Username,
		ID:       dbUser.ID,
	}

	return user, nil
}

func (srv *AuthService) Register(ctx context.Context, creds PasswordCredentials) (core.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)

	if err != nil {
		return core.User{}, err
	}

	org := core.Organization{
		ID:   core.IDType(uuid.New()),
		Name: fmt.Sprintf("%s's Organization", creds.Username),
	}

	user, err := srv.authRepository.RegisterUser(ctx, UserWithPassword{
		User: core.User{
			Username: creds.Username,
			Org:      org,
		},
		Password: string(hashedPassword),
	})

	if err != nil {
		return core.User{}, fmt.Errorf("Failed to create user: %w", err)
	}

	return user, nil
}

func (srv *AuthService) AuthenticateSession(ctx context.Context, sessionID string) (core.User, error) {
	session, err := srv.authRepository.GetSession(ctx, sessionID)

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

func (srv *AuthService) CreateSession(ctx context.Context, user core.User) (*Session, error) {
	session, err := NewSession(user)

	if err != nil {
		return nil, err
	}

	err = srv.authRepository.CreateSession(ctx, *session)

	if err != nil {
		return nil, err
	}

	createdSession, err := srv.authRepository.GetSession(ctx, session.ID)

	return &createdSession, err
}
