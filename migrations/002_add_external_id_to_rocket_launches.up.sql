-- Add external_id column to rocket_launches table
ALTER TABLE rocket_launches ADD COLUMN external_id BIGINT;

-- Create unique index on external_id to prevent duplicates
CREATE UNIQUE INDEX idx_rocket_launches_external_id ON rocket_launches(external_id) WHERE external_id IS NOT NULL;

-- Add comment to column
COMMENT ON COLUMN rocket_launches.external_id IS 'External ID from the RocketLaunch.Live API for deduplication';
