-- name: AddReview :one
INSERT INTO reviews(id, user_id, juice_id, rating, comment)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetJuiceReviews :many
SELECT reviews.*, users.username FROM reviews
JOIN users ON reviews.user_id = users.id
WHERE reviews.juice_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: GetReviewByID :one
SELECT * from reviews
WHERE id = $1;

-- name: DeleteReview :exec
DELETE FROM reviews WHERE id = $1;