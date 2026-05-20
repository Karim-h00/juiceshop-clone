-- name: GetAllJuice :many
SELECT * FROM juice 
ORDER BY created_at ASC;

-- name: GetJuiceByID :one
SELECT * FROM juice
WHERE id = $1;

-- name: GetJuiceByName :many
SELECT * FROM juice
WHERE name ILIKE $1;