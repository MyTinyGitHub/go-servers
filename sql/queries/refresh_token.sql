-- name: CreateRefreshToken :one
INSERT INTO refresh_token(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES(
  $1,
  NOW(),
  NOW(),
  $2,
  NOW() + (60 * INTERVAL '1 DAY'),
  NULL
)
RETURNING *;


-- name: GetRefreshTokenByToken :one
SELECT
  *
FROM
  refresh_token 
WHERE
  token = $1 
  AND expires_at > NOW()
  AND revoked_at IS NULL;

-- name: GetUserByRefreshToken :one
SELECT
  *
FROM
  users 
WHERE
  id = (select user_id from refresh_token where token = $1);

-- name: RevokeTokenByToken :exec
UPDATE
  refresh_token 
SET
  revoked_at = NOW(),
  updated_at = NOW() 
WHERE
  token = $1;

-- name: UpdateExpiresAtForRevoked :exec
UPDATE 
  refresh_token 
SET
  expires_at = NOW() + (60 * INTERVAL '1 DAY'),
  updated_at = NOW() 
WHERE
  token = $1;
