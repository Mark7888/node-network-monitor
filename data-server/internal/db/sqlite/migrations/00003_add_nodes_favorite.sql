-- +goose Up
ALTER TABLE nodes ADD COLUMN favorite INTEGER NOT NULL DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_nodes_favorite ON nodes(favorite);

-- +goose Down
DROP INDEX IF EXISTS idx_nodes_favorite;
ALTER TABLE nodes DROP COLUMN favorite;
