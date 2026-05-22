-- name: CreateOrder :one
INSERT INTO orders (id, total, created_at, user_id)
VALUES (gen_random_uuid(), @total, NOW(), @user_id)
RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (id, order_id, juice_id, quantity)
VALUES (gen_random_uuid(), @order_id, @juice_id, @quantity)
RETURNING *;

-- name: GetOrdersByUserID :many
SELECT id, created_at, total FROM orders
WHERE user_id = $1;

-- name: GetOrderItemsByOrderID :many
SELECT juice_id, quantity, juice.name FROM order_items
JOIN juice ON order_items.juice_id = juice.id
WHERE order_items.order_id = $1;