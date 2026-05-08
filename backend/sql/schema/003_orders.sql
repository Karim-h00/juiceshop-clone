-- +goose Up
CREATE TABLE orders(
    id UUID PRIMARY KEY,
    total INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE users;