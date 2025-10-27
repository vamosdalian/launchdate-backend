-- Remove unique index on external_id
DROP INDEX IF EXISTS idx_rocket_launches_external_id;

-- Remove external_id column from rocket_launches table
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS external_id;
