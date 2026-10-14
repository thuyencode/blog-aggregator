-- +goose Up
ALTER TABLE
    users
ALTER COLUMN
    created_at
SET
    NOT NULL,
ALTER COLUMN
    updated_at
SET
    NOT NULL;

-- +goose Down
ALTER TABLE
    users
ALTER COLUMN
    created_at DROP NOT NULL,
ALTER COLUMN
    updated_at DROP NOT NULL;
