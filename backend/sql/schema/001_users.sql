-- +goose Up
CREATE TABLE users(
    id UUID PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user'
);

-- +goose Down
DROP TABLE users;