-- Drop product launch related tables

-- Drop dependent tables first
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS launch_tags;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS launches;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS users;

-- Drop triggers for these tables
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_teams_updated_at ON teams;
DROP TRIGGER IF EXISTS update_launches_updated_at ON launches;
DROP TRIGGER IF EXISTS update_tasks_updated_at ON tasks;
DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;
