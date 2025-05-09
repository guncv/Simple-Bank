-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    full_name,
    email
) VALUES (
    sqlc.arg(username),
    sqlc.arg(hashed_password),
    sqlc.arg(full_name),
    sqlc.arg(email)
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = sqlc.arg(username)
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET 
    hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
    password_change_at = COALESCE(sqlc.narg(password_change_at), password_change_at),
    full_name = COALESCE(sqlc.narg(full_name), full_name),
    email = COALESCE(sqlc.narg(email), email)
WHERE 
    username = sqlc.arg(username) RETURNING *;