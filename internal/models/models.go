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
	Status      string     `json:"status" db:"status"` // draft, planned, in-progress, launched, cancelled
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
	Status      string     `json:"status" db:"status"` // todo, in-progress, done, blocked
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
	ID          int64      `json:"id" db:"id"`
	EntityType  string     `json:"entity_type" db:"entity_type"` // launch, milestone, task
	EntityID    int64      `json:"entity_id" db:"entity_id"`
	UserID      int64      `json:"user_id" db:"user_id"`
	Content     string     `json:"content" db:"content"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
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
