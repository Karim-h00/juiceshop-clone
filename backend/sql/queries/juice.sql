-- name: GetAllJuice :many
SELECT * FROM juice 
ORDER BY created_at ASC;