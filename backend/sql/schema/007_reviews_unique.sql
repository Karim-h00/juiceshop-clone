-- +goose Up
ALTER TABLE reviews ADD CONSTRAINT unique_user_juice_review UNIQUE (user_id, juice_id);

-- +goose Down
ALTER TABLE reviews DROP CONSTRAINT unique_user_juice_review;