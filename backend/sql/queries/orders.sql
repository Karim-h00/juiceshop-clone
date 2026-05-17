-- name: GetJuicesByIDs :many
SELECT id, price FROM juice WHERE id = ANY(@ids::uuid[]);

-- name: CreateOrder :one
INSERT INTO orders (id, total, created_at, user_id)
VALUES (gen_random_uuid(), @total, NOW(), @user_id)
RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (id, order_id, juice_id, quantity)
VALUES (gen_random_uuid(), @order_id, @juice_id, @quantity)
RETURNING *;