-- name: CreateChirp :one
INSERT INTO chirps(id, created_at, updated_at, user_id, body)
VALUES(
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING *;

-- name: GetAllChirpsAsc :many
SELECT * FROM chirps
ORDER BY created_at asc;

-- name: GetAllChirpsDesc :many
SELECT * FROM chirps
ORDER BY created_at desc;

-- name: GetAllChirpsByUserIdAsc :many
SELECT * FROM chirps
WHERE
  user_id = $1
ORDER BY created_at asc;

-- name: GetAllChirpsByUserIdDesc :many
SELECT * FROM chirps
WHERE
  user_id = $1
ORDER BY created_at desc;

-- name: GetChirpById :one
SELECT * FROM chirps
WHERE 
  ID = $1
ORDER BY created_at;

-- name: UserHasChirp :one
SELECT true FROM chirps
WHERE 
  ID = $1
  AND user_id = $2;

-- name: DeleteChirps :many
DELETE FROM chirps
RETURNING *;

-- name: DeleteChirpOfUser :exec
DELETE FROM chirps
WHERE id = $1
AND user_id = $2;


