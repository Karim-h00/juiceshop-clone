-- +goose Up
CREATE TABLE orders(
    id UUID PRIMARY KEY,
    total INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    juice_id UUID NOT NULL REFERENCES juice(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1
);

-- +goose Down
DROP TABLE order_items;
DROP TABLE orders;