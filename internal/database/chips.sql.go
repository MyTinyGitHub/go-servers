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
	UserID uuid.UUID
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

const deleteChirpOfUser = `-- name: DeleteChirpOfUser :exec
DELETE FROM chirps
WHERE id = $1
AND user_id = $2
`

type DeleteChirpOfUserParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteChirpOfUser(ctx context.Context, arg DeleteChirpOfUserParams) error {
	_, err := q.db.ExecContext(ctx, deleteChirpOfUser, arg.ID, arg.UserID)
	return err
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

const getAllChirpsAsc = `-- name: GetAllChirpsAsc :many
SELECT id, created_at, updated_at, body, user_id FROM chirps
ORDER BY created_at asc
`

func (q *Queries) GetAllChirpsAsc(ctx context.Context) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getAllChirpsAsc)
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

const getAllChirpsByUserIdAsc = `-- name: GetAllChirpsByUserIdAsc :many
SELECT id, created_at, updated_at, body, user_id FROM chirps
WHERE
  user_id = $1
ORDER BY created_at asc
`

func (q *Queries) GetAllChirpsByUserIdAsc(ctx context.Context, userID uuid.UUID) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getAllChirpsByUserIdAsc, userID)
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

const getAllChirpsByUserIdDesc = `-- name: GetAllChirpsByUserIdDesc :many
SELECT id, created_at, updated_at, body, user_id FROM chirps
WHERE
  user_id = $1
ORDER BY created_at desc
`

func (q *Queries) GetAllChirpsByUserIdDesc(ctx context.Context, userID uuid.UUID) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getAllChirpsByUserIdDesc, userID)
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

const getAllChirpsDesc = `-- name: GetAllChirpsDesc :many
SELECT id, created_at, updated_at, body, user_id FROM chirps
ORDER BY created_at desc
`

func (q *Queries) GetAllChirpsDesc(ctx context.Context) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getAllChirpsDesc)
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

func (q *Queries) GetChirpById(ctx context.Context, id uuid.UUID) (Chirp, error) {
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

const userHasChirp = `-- name: UserHasChirp :one
SELECT true FROM chirps
WHERE 
  ID = $1
  AND user_id = $2
`

type UserHasChirpParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) UserHasChirp(ctx context.Context, arg UserHasChirpParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, userHasChirp, arg.ID, arg.UserID)
	var column_1 bool
	err := row.Scan(&column_1)
	return column_1, err
}
