-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Create teams table
CREATE TABLE IF NOT EXISTS teams (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_teams_deleted_at ON teams(deleted_at);

-- Create team_members table
CREATE TABLE IF NOT EXISTS team_members (
    id BIGSERIAL PRIMARY KEY,
    team_id BIGINT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(team_id, user_id)
);

CREATE INDEX idx_team_members_team_id ON team_members(team_id);
CREATE INDEX idx_team_members_user_id ON team_members(user_id);

-- Create launches table
CREATE TABLE IF NOT EXISTS launches (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    launch_date TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    priority VARCHAR(50) NOT NULL DEFAULT 'medium',
    owner_id BIGINT NOT NULL REFERENCES users(id),
    team_id BIGINT REFERENCES teams(id),
    image_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_launches_owner_id ON launches(owner_id);
CREATE INDEX idx_launches_team_id ON launches(team_id);
CREATE INDEX idx_launches_status ON launches(status);
CREATE INDEX idx_launches_launch_date ON launches(launch_date);
CREATE INDEX idx_launches_deleted_at ON launches(deleted_at);

-- Create launch_tags table
CREATE TABLE IF NOT EXISTS launch_tags (
    launch_id BIGINT NOT NULL REFERENCES launches(id) ON DELETE CASCADE,
    tag VARCHAR(100) NOT NULL,
    PRIMARY KEY (launch_id, tag)
);

CREATE INDEX idx_launch_tags_tag ON launch_tags(tag);

-- Create tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    launch_id BIGINT NOT NULL REFERENCES launches(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    assignee_id BIGINT REFERENCES users(id),
    status VARCHAR(50) NOT NULL DEFAULT 'todo',
    priority VARCHAR(50) NOT NULL DEFAULT 'medium',
    due_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_tasks_launch_id ON tasks(launch_id);
CREATE INDEX idx_tasks_assignee_id ON tasks(assignee_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_deleted_at ON tasks(deleted_at);

-- Create comments table
CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_comments_entity ON comments(entity_type, entity_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_deleted_at ON comments(deleted_at);

-- Create companies table
CREATE TABLE IF NOT EXISTS companies (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    founded INT,
    founder VARCHAR(255),
    headquarters VARCHAR(255),
    employees INT,
    website VARCHAR(255),
    image_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_companies_name ON companies(name);
CREATE INDEX idx_companies_deleted_at ON companies(deleted_at);

-- Create rockets table
CREATE TABLE IF NOT EXISTS rockets (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    height NUMERIC(10, 2),
    diameter NUMERIC(10, 2),
    mass NUMERIC(12, 2),
    company_id BIGINT REFERENCES companies(id),
    image_url TEXT,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_rockets_name ON rockets(name);
CREATE INDEX idx_rockets_company_id ON rockets(company_id);
CREATE INDEX idx_rockets_active ON rockets(active);
CREATE INDEX idx_rockets_deleted_at ON rockets(deleted_at);

-- Create launch_bases table
CREATE TABLE IF NOT EXISTS launch_bases (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    country VARCHAR(100),
    description TEXT,
    image_url TEXT,
    latitude NUMERIC(10, 6),
    longitude NUMERIC(10, 6),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_launch_bases_name ON launch_bases(name);
CREATE INDEX idx_launch_bases_country ON launch_bases(country);
CREATE INDEX idx_launch_bases_deleted_at ON launch_bases(deleted_at);

-- Create rocket_launches table (different from product launches)
CREATE TABLE IF NOT EXISTS rocket_launches (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    launch_date TIMESTAMP WITH TIME ZONE NOT NULL,
    rocket_id BIGINT REFERENCES rockets(id),
    launch_base_id BIGINT REFERENCES launch_bases(id),
    provider_id BIGINT REFERENCES companies(id),
    status VARCHAR(50) NOT NULL DEFAULT 'scheduled',
    description TEXT,
    cospar_id VARCHAR(50),
    sort_date VARCHAR(50),
    date_str VARCHAR(100),
    slug VARCHAR(255),
    modified TIMESTAMP WITH TIME ZONE,
    window_open TIMESTAMP WITH TIME ZONE,
    t0 TIMESTAMP WITH TIME ZONE,
    window_close TIMESTAMP WITH TIME ZONE,
    mission_description TEXT,
    launch_description TEXT,
    weather_summary TEXT,
    weather_temp NUMERIC(5, 2),
    weather_condition VARCHAR(100),
    weather_wind_mph NUMERIC(5, 2),
    weather_icon VARCHAR(100),
    weather_updated TIMESTAMP WITH TIME ZONE,
    quicktext TEXT,
    suborbital BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_rocket_launches_name ON rocket_launches(name);
CREATE INDEX idx_rocket_launches_rocket_id ON rocket_launches(rocket_id);
CREATE INDEX idx_rocket_launches_launch_base_id ON rocket_launches(launch_base_id);
CREATE INDEX idx_rocket_launches_provider_id ON rocket_launches(provider_id);
CREATE INDEX idx_rocket_launches_status ON rocket_launches(status);
CREATE INDEX idx_rocket_launches_launch_date ON rocket_launches(launch_date);
CREATE INDEX idx_rocket_launches_cospar_id ON rocket_launches(cospar_id);
CREATE INDEX idx_rocket_launches_slug ON rocket_launches(slug);
CREATE INDEX idx_rocket_launches_sort_date ON rocket_launches(sort_date);
CREATE INDEX idx_rocket_launches_t0 ON rocket_launches(t0);
CREATE INDEX idx_rocket_launches_window_open ON rocket_launches(window_open);
CREATE INDEX idx_rocket_launches_modified ON rocket_launches(modified);
CREATE INDEX idx_rocket_launches_deleted_at ON rocket_launches(deleted_at);

-- Create rocket_launch_missions table (many-to-many relationship)
CREATE TABLE IF NOT EXISTS rocket_launch_missions (
    id BIGSERIAL PRIMARY KEY,
    rocket_launch_id BIGINT NOT NULL REFERENCES rocket_launches(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_rocket_launch_missions_launch_id ON rocket_launch_missions(rocket_launch_id);

-- Create rocket_launch_tags table (many-to-many relationship)
CREATE TABLE IF NOT EXISTS rocket_launch_tags (
    id BIGSERIAL PRIMARY KEY,
    rocket_launch_id BIGINT NOT NULL REFERENCES rocket_launches(id) ON DELETE CASCADE,
    text VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_rocket_launch_tags_launch_id ON rocket_launch_tags(rocket_launch_id);
CREATE INDEX idx_rocket_launch_tags_text ON rocket_launch_tags(text);

-- Create news table
CREATE TABLE IF NOT EXISTS news (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    summary TEXT,
    content TEXT,
    news_date TIMESTAMP WITH TIME ZONE NOT NULL,
    url VARCHAR(500),
    image_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_news_title ON news(title);
CREATE INDEX idx_news_date ON news(news_date);
CREATE INDEX idx_news_deleted_at ON news(deleted_at);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_teams_updated_at BEFORE UPDATE ON teams FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_launches_updated_at BEFORE UPDATE ON launches FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tasks_updated_at BEFORE UPDATE ON tasks FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_comments_updated_at BEFORE UPDATE ON comments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_companies_updated_at BEFORE UPDATE ON companies FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_rockets_updated_at BEFORE UPDATE ON rockets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_launch_bases_updated_at BEFORE UPDATE ON launch_bases FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_rocket_launches_updated_at BEFORE UPDATE ON rocket_launches FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_rocket_launch_missions_updated_at BEFORE UPDATE ON rocket_launch_missions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_news_updated_at BEFORE UPDATE ON news FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
