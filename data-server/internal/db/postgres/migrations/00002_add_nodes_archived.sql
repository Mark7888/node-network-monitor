-- +goose Up
ALTER TABLE nodes ADD COLUMN IF NOT EXISTS archived BOOLEAN NOT NULL DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_nodes_archived ON nodes(archived);

-- +goose Down
DROP INDEX IF EXISTS idx_nodes_archived;
ALTER TABLE nodes DROP COLUMN IF EXISTS archived;
