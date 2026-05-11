-- name: GetAllJuice :many
SELECT * FROM juice 
ORDER BY created_at ASC;

-- name: GetJuiceByID :one
SELECT * FROM juice
WHERE id = $1;