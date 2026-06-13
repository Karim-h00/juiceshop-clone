-- name: GetAllJuice :many
SELECT * FROM juice 
ORDER BY created_at ASC;

-- name: GetJuiceByID :one
SELECT * FROM juice
WHERE id = $1;

-- name: GetJuiceDetails :one
SELECT * FROM juice
WHERE name = $1;

-- name: GetJuicesByIDs :many
SELECT name, id, price, stock FROM juice WHERE id = ANY(@ids::uuid[]);

-- name: GetJuiceByName :many
SELECT * FROM juice
WHERE name ILIKE $1;

-- name: DecrementJuiceStock :exec
UPDATE juice
SET stock = stock - $1
WHERE id = $2;

-- name: AddJuice :one
INSERT INTO juice (id, name, description, price, created_at, updated_at, stock)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    now(),
    now(),
    $4
)
RETURNING *;

-- name: UpdateJuice :one
UPDATE juice
SET
    name = $1,
    description = $2,
    price = $3,
    stock = $4,
    updated_at = now()
WHERE id = $5
RETURNING *;

-- name: UpdateJuiceImage :exec
UPDATE juice
SET image_url = $1
WHERE id = $2;

-- name: DeleteJuice :exec
DELETE FROM juice 
WHERE id = $1;