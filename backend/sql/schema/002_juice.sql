-- +goose Up
CREATE TABLE juice(
    id UUID PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    image_url TEXT,
    stock INTEGER DEFAULT 0
);

-- +goose Down
DROP TABLE juice;