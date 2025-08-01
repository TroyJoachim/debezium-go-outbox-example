// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: users.sql

package DAL

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, email)
VALUES ($1, $2)
RETURNING id, username, email, created_at
`

type CreateUserParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Username, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, username, email, created_at
FROM users
WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id pgtype.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.CreatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, username, email, created_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserEmail = `-- name: UpdateUserEmail :one
UPDATE users
SET email = $2
WHERE id = $1
RETURNING id, username, email, created_at
`

type UpdateUserEmailParams struct {
	ID    pgtype.UUID `json:"id"`
	Email string      `json:"email"`
}

func (q *Queries) UpdateUserEmail(ctx context.Context, arg UpdateUserEmailParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUserEmail, arg.ID, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.CreatedAt,
	)
	return i, err
}
