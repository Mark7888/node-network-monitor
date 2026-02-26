-- +goose Up
ALTER TABLE nodes ADD COLUMN archived INTEGER NOT NULL DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_nodes_archived ON nodes(archived);

-- +goose Down
DROP INDEX IF EXISTS idx_nodes_archived;
ALTER TABLE nodes DROP COLUMN archived;
