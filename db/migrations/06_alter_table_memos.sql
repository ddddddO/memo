-- +migrate Up
ALTER TABLE memos ADD COLUMN is_exposed    BOOLEAN DEFAULT false;
ALTER TABLE memos ADD COLUMN exposed_at    TIMESTAMP;

-- +migrate Down
ALTER TABLE memos DROP COLUMN exposed_at;
ALTER TABLE memos DROP COLUMN is_exposed;
