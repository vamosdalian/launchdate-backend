-- Add new columns to rocket_launches table to match RocketLaunch.Live API structure

-- Add identifier and metadata fields
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS cospar_id VARCHAR(50);
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS sort_date VARCHAR(50);
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS date_str VARCHAR(100);
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS slug VARCHAR(255);
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS modified TIMESTAMP WITH TIME ZONE;

-- Add provider reference (launch service provider)
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS provider_id BIGINT REFERENCES companies(id);

-- Add launch window fields
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS window_open TIMESTAMP WITH TIME ZONE;
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS t0 TIMESTAMP WITH TIME ZONE;
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS window_close TIMESTAMP WITH TIME ZONE;

-- Add description fields
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS mission_description TEXT;
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS launch_description TEXT;

-- Add weather fields
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS weather_summary TEXT;
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS weather_temp NUMERIC(5, 2);
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS weather_condition VARCHAR(100);
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS weather_wind_mph NUMERIC(5, 2);
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS weather_icon VARCHAR(100);
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS weather_updated TIMESTAMP WITH TIME ZONE;

-- Add other fields
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS quicktext TEXT;
ALTER TABLE rocket_launches ADD COLUMN IF NOT EXISTS suborbital BOOLEAN DEFAULT false;

-- Rename launch_date to be clearer (keeping it for backward compatibility but using window times)
-- ALTER TABLE rocket_launches RENAME COLUMN launch_date TO legacy_launch_date;

-- Create indexes for new columns
CREATE INDEX IF NOT EXISTS idx_rocket_launches_cospar_id ON rocket_launches(cospar_id);
CREATE INDEX IF NOT EXISTS idx_rocket_launches_provider_id ON rocket_launches(provider_id);
CREATE INDEX IF NOT EXISTS idx_rocket_launches_slug ON rocket_launches(slug);
CREATE INDEX IF NOT EXISTS idx_rocket_launches_sort_date ON rocket_launches(sort_date);
CREATE INDEX IF NOT EXISTS idx_rocket_launches_t0 ON rocket_launches(t0);
CREATE INDEX IF NOT EXISTS idx_rocket_launches_window_open ON rocket_launches(window_open);
CREATE INDEX IF NOT EXISTS idx_rocket_launches_modified ON rocket_launches(modified);

-- Create table for rocket launch missions (many-to-many relationship)
CREATE TABLE IF NOT EXISTS rocket_launch_missions (
    id BIGSERIAL PRIMARY KEY,
    rocket_launch_id BIGINT NOT NULL REFERENCES rocket_launches(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_rocket_launch_missions_launch_id ON rocket_launch_missions(rocket_launch_id);

-- Create table for rocket launch tags (many-to-many relationship)
CREATE TABLE IF NOT EXISTS rocket_launch_tags (
    id BIGSERIAL PRIMARY KEY,
    rocket_launch_id BIGINT NOT NULL REFERENCES rocket_launches(id) ON DELETE CASCADE,
    text VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_rocket_launch_tags_launch_id ON rocket_launch_tags(rocket_launch_id);
CREATE INDEX IF NOT EXISTS idx_rocket_launch_tags_text ON rocket_launch_tags(text);

-- Create triggers for updated_at
DROP TRIGGER IF EXISTS update_rocket_launch_missions_updated_at ON rocket_launch_missions;
CREATE TRIGGER update_rocket_launch_missions_updated_at 
    BEFORE UPDATE ON rocket_launch_missions 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
