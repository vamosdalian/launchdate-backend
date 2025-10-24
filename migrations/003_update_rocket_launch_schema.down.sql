-- Drop new tables
DROP TABLE IF EXISTS rocket_launch_tags;
DROP TABLE IF EXISTS rocket_launch_missions;

-- Drop new indexes
DROP INDEX IF EXISTS idx_rocket_launches_modified;
DROP INDEX IF EXISTS idx_rocket_launches_window_open;
DROP INDEX IF EXISTS idx_rocket_launches_t0;
DROP INDEX IF EXISTS idx_rocket_launches_sort_date;
DROP INDEX IF EXISTS idx_rocket_launches_slug;
DROP INDEX IF EXISTS idx_rocket_launches_provider_id;
DROP INDEX IF EXISTS idx_rocket_launches_cospar_id;

-- Drop new columns from rocket_launches table
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS suborbital;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS quicktext;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS weather_updated;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS weather_icon;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS weather_wind_mph;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS weather_condition;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS weather_temp;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS weather_summary;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS mission_description;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS launch_description;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS window_close;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS t0;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS window_open;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS provider_id;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS modified;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS slug;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS date_str;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS sort_date;
ALTER TABLE rocket_launches DROP COLUMN IF EXISTS cospar_id;
