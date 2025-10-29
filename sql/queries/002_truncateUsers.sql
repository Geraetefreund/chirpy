-- name: TruncateUsers :exec
TRUNCATE TABLE users RESTART IDENTITY CASCADE;
