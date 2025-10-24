package models

import (
	"time"
)

// Launch represents a product or project launch
type Launch struct {
	ID          int64      `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	LaunchDate  time.Time  `json:"launch_date" db:"launch_date"`
	Status      string     `json:"status" db:"status"`     // draft, planned, in-progress, launched, cancelled
	Priority    string     `json:"priority" db:"priority"` // low, medium, high, critical
	OwnerID     int64      `json:"owner_id" db:"owner_id"`
	TeamID      *int64     `json:"team_id,omitempty" db:"team_id"`
	ImageURL    string     `json:"image_url,omitempty" db:"image_url"`
	Tags        []string   `json:"tags,omitempty" db:"-"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Milestone represents a key milestone in a launch
type Milestone struct {
	ID          int64      `json:"id" db:"id"`
	LaunchID    int64      `json:"launch_id" db:"launch_id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	DueDate     time.Time  `json:"due_date" db:"due_date"`
	Status      string     `json:"status" db:"status"` // pending, in-progress, completed, blocked
	Order       int        `json:"order" db:"order_num"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Task represents a specific task within a milestone or launch
type Task struct {
	ID          int64      `json:"id" db:"id"`
	LaunchID    int64      `json:"launch_id" db:"launch_id"`
	MilestoneID *int64     `json:"milestone_id,omitempty" db:"milestone_id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	AssigneeID  *int64     `json:"assignee_id,omitempty" db:"assignee_id"`
	Status      string     `json:"status" db:"status"`     // todo, in-progress, done, blocked
	Priority    string     `json:"priority" db:"priority"` // low, medium, high
	DueDate     *time.Time `json:"due_date,omitempty" db:"due_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Team represents a team working on launches
type Team struct {
	ID          int64      `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// User represents a user in the system
type User struct {
	ID        int64      `json:"id" db:"id"`
	Email     string     `json:"email" db:"email"`
	Name      string     `json:"name" db:"name"`
	AvatarURL string     `json:"avatar_url,omitempty" db:"avatar_url"`
	Role      string     `json:"role" db:"role"` // admin, manager, member
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// TeamMember represents a user's membership in a team
type TeamMember struct {
	ID        int64     `json:"id" db:"id"`
	TeamID    int64     `json:"team_id" db:"team_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Role      string    `json:"role" db:"role"` // lead, member
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Comment represents a comment on a launch, milestone, or task
type Comment struct {
	ID         int64      `json:"id" db:"id"`
	EntityType string     `json:"entity_type" db:"entity_type"` // launch, milestone, task
	EntityID   int64      `json:"entity_id" db:"entity_id"`
	UserID     int64      `json:"user_id" db:"user_id"`
	Content    string     `json:"content" db:"content"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// LaunchTag represents tags associated with launches
type LaunchTag struct {
	LaunchID int64  `json:"launch_id" db:"launch_id"`
	Tag      string `json:"tag" db:"tag"`
}

// CreateLaunchRequest represents the request to create a launch
type CreateLaunchRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	LaunchDate  time.Time `json:"launch_date" binding:"required"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	TeamID      *int64    `json:"team_id"`
	ImageURL    string    `json:"image_url"`
	Tags        []string  `json:"tags"`
}

// UpdateLaunchRequest represents the request to update a launch
type UpdateLaunchRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	LaunchDate  *time.Time `json:"launch_date"`
	Status      *string    `json:"status"`
	Priority    *string    `json:"priority"`
	TeamID      *int64     `json:"team_id"`
	ImageURL    *string    `json:"image_url"`
	Tags        []string   `json:"tags"`
}

// CreateMilestoneRequest represents the request to create a milestone
type CreateMilestoneRequest struct {
	LaunchID    int64     `json:"launch_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date" binding:"required"`
	Status      string    `json:"status"`
	Order       int       `json:"order"`
}

// UpdateMilestoneRequest represents the request to update a milestone
type UpdateMilestoneRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Status      *string    `json:"status"`
	Order       *int       `json:"order"`
}

// CreateTaskRequest represents the request to create a task
type CreateTaskRequest struct {
	LaunchID    int64      `json:"launch_id" binding:"required"`
	MilestoneID *int64     `json:"milestone_id"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	AssigneeID  *int64     `json:"assignee_id"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
}

// UpdateTaskRequest represents the request to update a task
type UpdateTaskRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	AssigneeID  *int64     `json:"assignee_id"`
	Status      *string    `json:"status"`
	Priority    *string    `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	MilestoneID *int64     `json:"milestone_id"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Redis    string `json:"redis"`
}

// Company represents a space company
type Company struct {
	ID           int64      `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	Description  string     `json:"description" db:"description"`
	Founded      int        `json:"founded" db:"founded"`
	Founder      string     `json:"founder" db:"founder"`
	Headquarters string     `json:"headquarters" db:"headquarters"`
	Employees    int        `json:"employees" db:"employees"`
	Website      string     `json:"website" db:"website"`
	ImageURL     string     `json:"imageUrl" db:"image_url"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Rocket represents a rocket
type Rocket struct {
	ID          int64      `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	Height      float64    `json:"height" db:"height"`
	Diameter    float64    `json:"diameter" db:"diameter"`
	Mass        float64    `json:"mass" db:"mass"`
	CompanyID   *int64     `json:"company_id,omitempty" db:"company_id"`
	Company     string     `json:"company" db:"-"`
	ImageURL    string     `json:"imageUrl" db:"image_url"`
	Active      bool       `json:"active" db:"active"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// LaunchBase represents a launch site
type LaunchBase struct {
	ID          int64      `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Location    string     `json:"location" db:"location"`
	Country     string     `json:"country" db:"country"`
	Description string     `json:"description" db:"description"`
	ImageURL    string     `json:"imageUrl" db:"image_url"`
	Latitude    float64    `json:"latitude" db:"latitude"`
	Longitude   float64    `json:"longitude" db:"longitude"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// RocketLaunch represents a rocket launch event
type RocketLaunch struct {
	ID                 int64                 `json:"id" db:"id"`
	CosparID           string                `json:"cospar_id" db:"cospar_id"`
	SortDate           string                `json:"sort_date" db:"sort_date"`
	Name               string                `json:"name" db:"name"`
	Provider           *RocketLaunchProvider `json:"provider,omitempty" db:"-"`
	ProviderID         *int64                `json:"provider_id,omitempty" db:"provider_id"`
	Vehicle            *RocketLaunchVehicle  `json:"vehicle,omitempty" db:"-"`
	RocketID           *int64                `json:"rocket_id,omitempty" db:"rocket_id"`
	Pad                *RocketLaunchPad      `json:"pad,omitempty" db:"-"`
	LaunchBaseID       *int64                `json:"launch_base_id,omitempty" db:"launch_base_id"`
	Missions           []RocketLaunchMission `json:"missions,omitempty" db:"-"`
	MissionDescription string                `json:"mission_description" db:"mission_description"`
	LaunchDescription  string                `json:"launch_description" db:"launch_description"`
	WindowOpen         *time.Time            `json:"win_open,omitempty" db:"window_open"`
	T0                 *time.Time            `json:"t0,omitempty" db:"t0"`
	WindowClose        *time.Time            `json:"win_close,omitempty" db:"window_close"`
	DateStr            string                `json:"date_str" db:"date_str"`
	Tags               []RocketLaunchTag     `json:"tags,omitempty" db:"-"`
	Slug               string                `json:"slug" db:"slug"`
	WeatherSummary     string                `json:"weather_summary" db:"weather_summary"`
	WeatherTemp        *float32              `json:"weather_temp,omitempty" db:"weather_temp"`
	WeatherCondition   string                `json:"weather_condition" db:"weather_condition"`
	WeatherWindMPH     *float32              `json:"weather_wind_mph,omitempty" db:"weather_wind_mph"`
	WeatherIcon        string                `json:"weather_icon" db:"weather_icon"`
	WeatherUpdated     *time.Time            `json:"weather_updated,omitempty" db:"weather_updated"`
	QuickText          string                `json:"quicktext" db:"quicktext"`
	Suborbital         bool                  `json:"suborbital" db:"suborbital"`
	Modified           *time.Time            `json:"modified,omitempty" db:"modified"`
	Status             string                `json:"status" db:"status"` // scheduled, successful, failed, cancelled
	CreatedAt          time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time            `json:"deleted_at,omitempty" db:"deleted_at"`
}

// RocketLaunchProvider represents the launch service provider
type RocketLaunchProvider struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
}

// RocketLaunchVehicle represents the launch vehicle
type RocketLaunchVehicle struct {
	ID        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	CompanyID *int64 `json:"company_id,omitempty" db:"company_id"`
	Slug      string `json:"slug" db:"slug"`
}

// RocketLaunchPad represents the launch pad
type RocketLaunchPad struct {
	ID       int64                    `json:"id" db:"id"`
	Name     string                   `json:"name" db:"name"`
	Location *RocketLaunchPadLocation `json:"location,omitempty" db:"-"`
}

// RocketLaunchPadLocation represents the launch pad location
type RocketLaunchPadLocation struct {
	ID        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	State     string `json:"state" db:"state"`
	StateName string `json:"statename" db:"statename"`
	Country   string `json:"country" db:"country"`
	Slug      string `json:"slug" db:"slug"`
}

// RocketLaunchMission represents a mission
type RocketLaunchMission struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

// RocketLaunchTag represents a tag
type RocketLaunchTag struct {
	ID   int64  `json:"id" db:"id"`
	Text string `json:"text" db:"text"`
}

// News represents a news article
type News struct {
	ID        int64      `json:"id" db:"id"`
	Title     string     `json:"title" db:"title"`
	Summary   string     `json:"summary" db:"summary"`
	Content   string     `json:"content,omitempty" db:"content"`
	NewsDate  time.Time  `json:"date" db:"news_date"`
	URL       string     `json:"url" db:"url"`
	ImageURL  string     `json:"imageUrl" db:"image_url"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// CreateCompanyRequest represents the request to create a company
type CreateCompanyRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	Founded      int    `json:"founded"`
	Founder      string `json:"founder"`
	Headquarters string `json:"headquarters"`
	Employees    int    `json:"employees"`
	Website      string `json:"website"`
	ImageURL     string `json:"imageUrl"`
}

// CreateRocketRequest represents the request to create a rocket
type CreateRocketRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Height      float64 `json:"height"`
	Diameter    float64 `json:"diameter"`
	Mass        float64 `json:"mass"`
	CompanyID   *int64  `json:"company_id"`
	ImageURL    string  `json:"imageUrl"`
	Active      bool    `json:"active"`
}

// CreateLaunchBaseRequest represents the request to create a launch base
type CreateLaunchBaseRequest struct {
	Name        string  `json:"name" binding:"required"`
	Location    string  `json:"location"`
	Country     string  `json:"country"`
	Description string  `json:"description"`
	ImageURL    string  `json:"imageUrl"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

// CreateRocketLaunchRequest represents the request to create a rocket launch
type CreateRocketLaunchRequest struct {
	CosparID           string     `json:"cospar_id"`
	SortDate           string     `json:"sort_date"`
	Name               string     `json:"name" binding:"required"`
	ProviderID         *int64     `json:"provider_id"`
	RocketID           *int64     `json:"rocket_id"`
	LaunchBaseID       *int64     `json:"launch_base_id"`
	MissionDescription string     `json:"mission_description"`
	LaunchDescription  string     `json:"launch_description"`
	WindowOpen         *time.Time `json:"win_open"`
	T0                 *time.Time `json:"t0"`
	WindowClose        *time.Time `json:"win_close"`
	DateStr            string     `json:"date_str" binding:"required"`
	Slug               string     `json:"slug"`
	WeatherSummary     string     `json:"weather_summary"`
	WeatherTemp        *float32   `json:"weather_temp"`
	WeatherCondition   string     `json:"weather_condition"`
	WeatherWindMPH     *float32   `json:"weather_wind_mph"`
	WeatherIcon        string     `json:"weather_icon"`
	WeatherUpdated     *time.Time `json:"weather_updated"`
	QuickText          string     `json:"quicktext"`
	Suborbital         bool       `json:"suborbital"`
	Modified           *time.Time `json:"modified"`
	Status             string     `json:"status"`
}

// CreateNewsRequest represents the request to create a news article
type CreateNewsRequest struct {
	Title    string    `json:"title" binding:"required"`
	Summary  string    `json:"summary"`
	Content  string    `json:"content"`
	NewsDate time.Time `json:"date" binding:"required"`
	URL      string    `json:"url"`
	ImageURL string    `json:"imageUrl"`
}
