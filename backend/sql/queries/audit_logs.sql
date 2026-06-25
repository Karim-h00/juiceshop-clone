-- name: AddLog :exec
INSERT INTO audit_logs (id, user_id, action, target_type, target_id, target_name, created_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);