-- +goose Up
ALTER TABLE nodes ADD COLUMN IF NOT EXISTS favorite BOOLEAN NOT NULL DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_nodes_favorite ON nodes(favorite);

-- +goose Down
DROP INDEX IF EXISTS idx_nodes_favorite;
ALTER TABLE nodes DROP COLUMN IF EXISTS favorite;
