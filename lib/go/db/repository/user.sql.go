// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const addUser = `-- name: AddUser :exec
INSERT INTO users(username, password) VALUES($1,$2)
`

type AddUserParams struct {
	Username string
	Password string
}

func (q *Queries) AddUser(ctx context.Context, arg AddUserParams) error {
	_, err := q.db.Exec(ctx, addUser, arg.Username, arg.Password)
	return err
}

const getSession = `-- name: GetSession :one
SELECT sessions.id, sessions.user_id, sessions.created_at, sessions.revoked_at, sessions.expires_at, sessions.last_active_at, users.username
FROM sessions JOIN users on users.id = sessions.user_id
WHERE sessions.id = $1
`

type GetSessionRow struct {
	ID           string
	UserID       pgtype.Int4
	CreatedAt    time.Time
	RevokedAt    time.Time
	ExpiresAt    time.Time
	LastActiveAt time.Time
	Username     string
}

func (q *Queries) GetSession(ctx context.Context, id string) (GetSessionRow, error) {
	row := q.db.QueryRow(ctx, getSession, id)
	var i GetSessionRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.CreatedAt,
		&i.RevokedAt,
		&i.ExpiresAt,
		&i.LastActiveAt,
		&i.Username,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT
  id, username, password
FROM
  users
WHERE
  username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(&i.ID, &i.Username, &i.Password)
	return i, err
}
