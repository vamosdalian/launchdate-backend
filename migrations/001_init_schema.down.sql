-- Drop triggers
DROP TRIGGER IF EXISTS update_news_updated_at ON news;
DROP TRIGGER IF EXISTS update_rocket_launch_missions_updated_at ON rocket_launch_missions;
DROP TRIGGER IF EXISTS update_rocket_launches_updated_at ON rocket_launches;
DROP TRIGGER IF EXISTS update_launch_bases_updated_at ON launch_bases;
DROP TRIGGER IF EXISTS update_rockets_updated_at ON rockets;
DROP TRIGGER IF EXISTS update_companies_updated_at ON companies;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS rocket_launch_tags;
DROP TABLE IF EXISTS rocket_launch_missions;
DROP TABLE IF EXISTS news;
DROP TABLE IF EXISTS rocket_launches;
DROP TABLE IF EXISTS launch_bases;
DROP TABLE IF EXISTS rockets;
DROP TABLE IF EXISTS companies;
