-- +goose Up
ALTER TABLE nodes ADD COLUMN IF NOT EXISTS location VARCHAR(255);
CREATE INDEX IF NOT EXISTS idx_nodes_location ON nodes(location);

-- +goose Down
DROP INDEX IF EXISTS idx_nodes_location;
ALTER TABLE nodes DROP COLUMN IF EXISTS location;
