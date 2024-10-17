// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Session struct {
	ID           string
	UserID       pgtype.Int4
	CreatedAt    time.Time
	RevokedAt    time.Time
	ExpiresAt    time.Time
	LastActiveAt time.Time
}

type User struct {
	ID       int32
	Username string
	Password string
}
