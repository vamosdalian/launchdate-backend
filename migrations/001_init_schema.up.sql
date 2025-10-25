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
CREATE TRIGGER update_companies_updated_at BEFORE UPDATE ON companies FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_rockets_updated_at BEFORE UPDATE ON rockets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_launch_bases_updated_at BEFORE UPDATE ON launch_bases FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_rocket_launches_updated_at BEFORE UPDATE ON rocket_launches FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_rocket_launch_missions_updated_at BEFORE UPDATE ON rocket_launch_missions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_news_updated_at BEFORE UPDATE ON news FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
