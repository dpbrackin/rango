// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package generated

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Document struct {
	ID     int32
	UserID pgtype.Int4
	Source string
}

type Session struct {
	ID           string
	UserID       pgtype.Int4
	CreatedAt    pgtype.Timestamptz
	RevokedAt    pgtype.Timestamptz
	ExpiresAt    pgtype.Timestamptz
	LastActiveAt pgtype.Timestamptz
}

type User struct {
	ID       int32
	Username string
	Password string
}