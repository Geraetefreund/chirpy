-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING *;

-- name: LookUpUserByEmail :one
SELECT * FROM users WHERE email =$1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at, user_id)
VALUES (
  $1,
  NOW(),
  NOW(),
  $2,
  $3
  )
RETURNING *;

-- name: UpdateEmailAndPW :one
UPDATE users
SET email = $1,
  hashed_password = $2,
  updated_at = Now()
WHERE id = $3
RETURNING *;

-- name: UpgradeChirpyPlus :one
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1
RETURNING *;
