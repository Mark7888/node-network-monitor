-- +goose Up
ALTER TABLE nodes ADD COLUMN location TEXT;
CREATE INDEX IF NOT EXISTS idx_nodes_location ON nodes(location);

-- +goose Down
DROP INDEX IF EXISTS idx_nodes_location;
ALTER TABLE nodes DROP COLUMN location;
