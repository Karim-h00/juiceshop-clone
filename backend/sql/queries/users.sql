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

-- name: GetPasswordByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserRole :one
SELECT role FROM users
WHERE id = $1;