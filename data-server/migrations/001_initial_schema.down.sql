-- Rollback initial schema

DROP INDEX IF EXISTS idx_api_keys_key_hash;
DROP INDEX IF EXISTS idx_api_keys_enabled;
DROP TABLE IF EXISTS api_keys;

DROP INDEX IF EXISTS idx_failed_timestamp;
DROP INDEX IF EXISTS idx_failed_node_id;
DROP TABLE IF EXISTS failed_measurements;

DROP INDEX IF EXISTS idx_measurements_node_timestamp;
DROP INDEX IF EXISTS idx_measurements_created_at;
DROP INDEX IF EXISTS idx_measurements_timestamp;
DROP INDEX IF EXISTS idx_measurements_node_id;
DROP TABLE IF EXISTS measurements;

DROP INDEX IF EXISTS idx_nodes_name;
DROP INDEX IF EXISTS idx_nodes_last_alive;
DROP INDEX IF EXISTS idx_nodes_status;
DROP TABLE IF EXISTS nodes;
