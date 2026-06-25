-- name: GetAllJuice :many
SELECT juice.*,
COALESCE(AVG(reviews.rating), 0) AS avg_rating,
COUNT(reviews.id) AS reviews_count
FROM juice
LEFT JOIN reviews ON reviews.juice_id = juice.id
GROUP BY juice.id
ORDER BY juice.name ASC;

-- name: GetJuiceByID :one
SELECT * FROM juice
WHERE id = $1;

-- name: GetJuiceID :one
SELECT id from juice
WHERE name = $1;

-- name: GetJuiceDetails :one
SELECT juice.*,
COALESCE(AVG(reviews.rating), 0) AS avg_rating,
COUNT(reviews.id) AS reviews_count
FROM juice
LEFT JOIN reviews ON reviews.juice_id = juice.id
WHERE name = $1
GROUP BY juice.id;

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
    $4,
    $5,
    $6
)
RETURNING *;

-- name: UpdateJuice :one
UPDATE juice
SET
    name = $1,
    description = $2,
    price = $3,
    stock = $4,
    updated_at = $5
WHERE id = $6
RETURNING *;

-- name: UpdateJuiceImage :exec
UPDATE juice
SET 
    image_url = $1,
    updated_at = $2
WHERE id = $3;

-- name: DeleteJuice :exec
DELETE FROM juice 
WHERE id = $1;