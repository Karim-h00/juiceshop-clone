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
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: GetOrderItemsByOrderID :many
SELECT juice_id, quantity, juice.name FROM order_items
JOIN juice ON order_items.juice_id = juice.id
WHERE order_items.order_id = $1;

-- name: GetOrderByOrderID :one
SELECT * FROM orders 
WHERE id = $1;

-- name: GetAllOrders :many
SELECT orders.*, users.username FROM orders
JOIN users ON orders.user_id = users.id
ORDER BY orders.created_at DESC
LIMIT 10
OFFSET $1;

-- name: DeleteOrderByOrderID :exec
DELETE FROM orders 
where id = $1;