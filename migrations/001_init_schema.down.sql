-- Drop triggers
DROP TRIGGER IF EXISTS update_news_updated_at ON news;
DROP TRIGGER IF EXISTS update_rocket_launch_missions_updated_at ON rocket_launch_missions;
DROP TRIGGER IF EXISTS update_rocket_launches_updated_at ON rocket_launches;
DROP TRIGGER IF EXISTS update_launch_bases_updated_at ON launch_bases;
DROP TRIGGER IF EXISTS update_rockets_updated_at ON rockets;
DROP TRIGGER IF EXISTS update_companies_updated_at ON companies;
DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;
DROP TRIGGER IF EXISTS update_tasks_updated_at ON tasks;
DROP TRIGGER IF EXISTS update_launches_updated_at ON launches;
DROP TRIGGER IF EXISTS update_teams_updated_at ON teams;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

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
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS launch_tags;
DROP TABLE IF EXISTS launches;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS users;
