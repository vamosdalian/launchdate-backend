-- Remove external_id from rocket_launch_missions
DROP INDEX IF EXISTS idx_rocket_launch_missions_external_id;
ALTER TABLE rocket_launch_missions DROP COLUMN IF EXISTS external_id;

-- Remove external_id from rockets
DROP INDEX IF EXISTS idx_rockets_external_id;
ALTER TABLE rockets DROP COLUMN IF EXISTS external_id;

-- Remove external_id from launch_bases
DROP INDEX IF EXISTS idx_launch_bases_external_id;
ALTER TABLE launch_bases DROP COLUMN IF EXISTS external_id;

-- Remove external_id from companies
DROP INDEX IF EXISTS idx_companies_external_id;
ALTER TABLE companies DROP COLUMN IF EXISTS external_id;
