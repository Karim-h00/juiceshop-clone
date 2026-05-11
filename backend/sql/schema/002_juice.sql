-- +goose Up
CREATE TABLE juice(
    id UUID PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    price INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    image_url TEXT NOT NULL DEFAULT 'https://placehold.co/300x300',
    stock INTEGER NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE juice;