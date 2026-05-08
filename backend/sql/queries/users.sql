-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, username, email, hashed_password)
VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2,
    $3
)
RETURNING *;