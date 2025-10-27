-- Add external_id column to companies table
ALTER TABLE companies ADD COLUMN external_id BIGINT;

-- Create unique index on external_id to prevent duplicates
CREATE UNIQUE INDEX idx_companies_external_id ON companies(external_id) WHERE external_id IS NOT NULL;

-- Add comment to column
COMMENT ON COLUMN companies.external_id IS 'External ID from the RocketLaunch.Live API for deduplication';

-- Add external_id column to launch_bases table
ALTER TABLE launch_bases ADD COLUMN external_id BIGINT;

-- Create unique index on external_id to prevent duplicates
CREATE UNIQUE INDEX idx_launch_bases_external_id ON launch_bases(external_id) WHERE external_id IS NOT NULL;

-- Add comment to column
COMMENT ON COLUMN launch_bases.external_id IS 'External ID from the RocketLaunch.Live API for deduplication';

-- Add external_id column to rockets table
ALTER TABLE rockets ADD COLUMN external_id BIGINT;

-- Create unique index on external_id to prevent duplicates
CREATE UNIQUE INDEX idx_rockets_external_id ON rockets(external_id) WHERE external_id IS NOT NULL;

-- Add comment to column
COMMENT ON COLUMN rockets.external_id IS 'External ID from the RocketLaunch.Live API for deduplication';

-- Add external_id column to rocket_launch_missions table
ALTER TABLE rocket_launch_missions ADD COLUMN external_id BIGINT;

-- Create unique index on external_id to prevent duplicates
CREATE UNIQUE INDEX idx_rocket_launch_missions_external_id ON rocket_launch_missions(external_id) WHERE external_id IS NOT NULL;

-- Add comment to column
COMMENT ON COLUMN rocket_launch_missions.external_id IS 'External ID from the RocketLaunch.Live API for deduplication';
