// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: chips.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createChirp = `-- name: CreateChirp :one
INSERT INTO chirps(id, created_at, updated_at, user_id, body)
VALUES(
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING id, created_at, updated_at, body, user_id
`

type CreateChirpParams struct {
	UserID uuid.NullUUID
	Body   string
}

func (q *Queries) CreateChirp(ctx context.Context, arg CreateChirpParams) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, createChirp, arg.UserID, arg.Body)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}

const deleteChirps = `-- name: DeleteChirps :many
DELETE FROM chirps
RETURNING id, created_at, updated_at, body, user_id
`

func (q *Queries) DeleteChirps(ctx context.Context) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, deleteChirps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllChirps = `-- name: GetAllChirps :many
SELECT id, created_at, updated_at, body, user_id FROM chirps
ORDER BY created_at
`

func (q *Queries) GetAllChirps(ctx context.Context) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getAllChirps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getChirpById = `-- name: GetChirpById :one
SELECT id, created_at, updated_at, body, user_id FROM chirps
WHERE 
  ID = $1
ORDER BY created_at
`

func (q *Queries) GetChirpById(ctx context.Context, id uuid.NullUUID) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, getChirpById, id)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}
