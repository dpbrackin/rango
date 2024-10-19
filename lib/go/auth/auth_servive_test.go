package auth_test

import (
	"context"
	"fmt"
	"rango/auth"
	"rango/core"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthRepository struct {
	mock.Mock
}

type mockClock struct {
	mock.Mock
}

// Now implements core.Clock.
func (m *mockClock) Now() time.Time {
	args := m.Called()

	return args.Get(0).(time.Time)
}

// AddUser implements auth.AuthRepository.
func (m *mockAuthRepository) AddUser(ctx context.Context, params auth.UserWithPassword) error {
	args := m.Called(ctx, params)

	return args.Error(0)
}

// GetSession implements auth.AuthRepository.
func (m *mockAuthRepository) GetSession(ctx context.Context, sessionID string) (auth.Session, error) {
	args := m.Called(ctx, sessionID)
	// Check if the returned value is of type auth.Session
	session, ok := args.Get(0).(auth.Session)
	if !ok {
		return auth.Session{}, fmt.Errorf("expected auth.Session, but got %T", args.Get(0))
	}

	return session, args.Error(1)
}

// GetUserByUsername implements auth.AuthRepository.
func (m *mockAuthRepository) GetUserByUsername(ctx context.Context, username string) (auth.UserWithPassword, error) {
	args := m.Called(ctx, username)
	// Check if the returned value is of type auth.Session
	user, ok := args.Get(0).(auth.UserWithPassword)
	if !ok {
		return auth.UserWithPassword{}, fmt.Errorf("expected auth.Session, but got %T", args.Get(0))
	}

	return user, args.Error(1)
}

func TestValidAuthenticateSession(t *testing.T) {
	validSession := auth.Session{
		ID: "session1",
		User: core.User{
			Username: "user1",
			ID:       core.IDType(1),
		},
		CreatedAt:    time.Date(2024, 12, 12, 0, 0, 0, 0, time.UTC),
		RevokedAt:    time.Time{},
		ExpiresAt:    time.Date(2025, 12, 12, 0, 0, 0, 0, time.UTC),
		LastActiveAt: time.Date(2024, 12, 12, 0, 0, 0, 0, time.UTC),
		IsRevoked:    false,
	}

	repository := new(mockAuthRepository)
	repository.On("GetSession", mock.Anything, "session1").Return(validSession, nil)

	fakeClock := new(mockClock)
	fakeClock.On("Now").Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

	service := auth.NewAuthService(auth.NewAuthServiceParams{
		Repository: repository,
		Clock:      fakeClock,
	})

	ctx := context.Background()
	res, err := service.AuthenticateSession(ctx, "session1")

	assert.Nil(t, err)
	assert.Equal(t, validSession.User, res)
}
