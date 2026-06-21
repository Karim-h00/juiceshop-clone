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

-- name: GetPasswordByUserID :one
SELECT hashed_password FROM users
WHERE ID = $1;

-- name: GetUserRole :one
SELECT role FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET
    username = $1,
    email = $2,
    updated_at = now()
WHERE id = $3
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET
    hashed_password = $1,
    updated_at = now()
WHERE id = $2;

-- name: GetUserByID :one
SELECT * from users
WHERE id = $1;

-- name: GetAllUsers :many
SELECT * from users
ORDER BY username
LIMIT $1
OFFSET $2;

-- name: SearchUsers :many
SELECT * FROM users
WHERE username ILIKE '%' || $1 || '%'
OR email ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 
OFFSET $3;

-- name: UpdateUserRole :exec
UPDATE users
SET
    role = $1,
    updated_at = now()
WHERE id = $2;

-- name: DeleteUserByID :exec
DELETE FROM users WHERE id = $1;