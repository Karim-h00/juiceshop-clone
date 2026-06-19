-- +goose Up
ALTER TABLE juice ADD CONSTRAINT juice_name_unique UNIQUE (name);

-- +goose Down
ALTER TABLE juice DROP CONSTRAINT juice_name_unique;