-- name: GetAllJuice :many
SELECT * FROM juice 
ORDER BY created_at ASC;

-- name: GetJuiceByID :one
SELECT * FROM juice
WHERE id = $1;

-- name: GetJuicesByIDs :many
SELECT name, id, price, stock FROM juice WHERE id = ANY(@ids::uuid[]);

-- name: GetJuiceByName :many
SELECT * FROM juice
WHERE name ILIKE $1;

-- name: DecrementJuiceStock :exec
UPDATE juice
SET stock = stock - $1
WHERE id = $2;