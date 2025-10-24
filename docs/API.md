# LaunchDate API Documentation

## Base URL

```
Development: http://localhost:8080
Production: https://api.launch-date.com
```

## API Version

Current version: `v1`

All endpoints are prefixed with `/api/v1` except for the health check endpoint.

## Authentication

Authentication will be implemented in a future version. Currently, all endpoints are publicly accessible.

## Response Format

All API responses follow a consistent JSON format:

**Success Response:**
```json
{
  "id": 1,
  "title": "Example Launch",
  ...
}
```

**Error Response:**
```json
{
  "error": "Error message describing what went wrong"
}
```

## Endpoints

### Health Check

#### GET /health

Check the health status of the API and its dependencies.

**Response:**
```json
{
  "status": "ok",
  "database": "ok",
  "redis": "ok"
}
```

Status values: `ok`, `degraded`, `error`

---

## Launches

### List Launches

#### GET /api/v1/launches

Retrieve a list of launches with optional filters.

**Query Parameters:**
- `status` (optional): Filter by status (`draft`, `planned`, `in-progress`, `launched`, `cancelled`)
- `priority` (optional): Filter by priority (`low`, `medium`, `high`, `critical`)
- `team_id` (optional): Filter by team ID
- `limit` (optional): Maximum number of results (default: 50, max: 100)
- `offset` (optional): Number of results to skip (default: 0)

**Example Request:**
```bash
curl http://localhost:8080/api/v1/launches?status=in-progress&limit=10
```

**Response:**
```json
[
  {
    "id": 1,
    "title": "Product Launch 2024",
    "description": "Launch of our new product line",
    "launch_date": "2024-12-15T00:00:00Z",
    "status": "in-progress",
    "priority": "high",
    "owner_id": 1,
    "team_id": 5,
    "image_url": "https://example.com/image.jpg",
    "tags": ["product", "2024"],
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-15T14:30:00Z"
  }
]
```

### Create Launch

#### POST /api/v1/launches

Create a new launch.

**Request Body:**
```json
{
  "title": "New Product Launch",
  "description": "Exciting new product",
  "launch_date": "2024-12-20T00:00:00Z",
  "status": "draft",
  "priority": "high",
  "team_id": 5,
  "image_url": "https://example.com/product.jpg",
  "tags": ["product", "q4"]
}
```

**Required Fields:**
- `title`
- `launch_date`

**Optional Fields:**
- `description`
- `status` (default: `draft`)
- `priority` (default: `medium`)
- `team_id`
- `image_url`
- `tags`

**Response:** `201 Created`
```json
{
  "id": 2,
  "title": "New Product Launch",
  "description": "Exciting new product",
  "launch_date": "2024-12-20T00:00:00Z",
  "status": "draft",
  "priority": "high",
  "owner_id": 1,
  "team_id": 5,
  "image_url": "https://example.com/product.jpg",
  "tags": ["product", "q4"],
  "created_at": "2024-01-20T10:00:00Z",
  "updated_at": "2024-01-20T10:00:00Z"
}
```

### Get Launch

#### GET /api/v1/launches/:id

Retrieve details of a specific launch.

**Parameters:**
- `id` (path): Launch ID

**Response:** `200 OK`
```json
{
  "id": 1,
  "title": "Product Launch 2024",
  "description": "Launch of our new product line",
  "launch_date": "2024-12-15T00:00:00Z",
  "status": "in-progress",
  "priority": "high",
  "owner_id": 1,
  "team_id": 5,
  "image_url": "https://example.com/image.jpg",
  "tags": ["product", "2024"],
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-15T14:30:00Z"
}
```

### Update Launch

#### PUT /api/v1/launches/:id

Update an existing launch.

**Parameters:**
- `id` (path): Launch ID

**Request Body:**
All fields are optional. Only include fields you want to update.

```json
{
  "title": "Updated Product Launch",
  "status": "in-progress",
  "priority": "critical"
}
```

**Response:** `200 OK`
```json
{
  "message": "launch updated successfully"
}
```

### Delete Launch

#### DELETE /api/v1/launches/:id

Soft delete a launch.

**Parameters:**
- `id` (path): Launch ID

**Response:** `204 No Content`

---

## Milestones

### List Milestones for Launch

#### GET /api/v1/launches/:launch_id/milestones

Retrieve all milestones for a specific launch.

**Parameters:**
- `launch_id` (path): Launch ID

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "launch_id": 1,
    "title": "Design Complete",
    "description": "Complete all design mockups",
    "due_date": "2024-10-15T00:00:00Z",
    "status": "completed",
    "order": 1,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-10-15T16:00:00Z"
  }
]
```

### Create Milestone

#### POST /api/v1/milestones

Create a new milestone.

**Request Body:**
```json
{
  "launch_id": 1,
  "title": "Beta Testing",
  "description": "Complete beta testing phase",
  "due_date": "2024-11-01T00:00:00Z",
  "status": "pending",
  "order": 2
}
```

**Required Fields:**
- `launch_id`
- `title`
- `due_date`

**Response:** `201 Created`

### Get Milestone

#### GET /api/v1/milestones/:id

Retrieve details of a specific milestone.

**Parameters:**
- `id` (path): Milestone ID

**Response:** `200 OK`

### Update Milestone

#### PUT /api/v1/milestones/:id

Update an existing milestone.

**Parameters:**
- `id` (path): Milestone ID

**Request Body:**
```json
{
  "status": "in-progress",
  "due_date": "2024-11-05T00:00:00Z"
}
```

**Response:** `200 OK`

### Delete Milestone

#### DELETE /api/v1/milestones/:id

Soft delete a milestone.

**Parameters:**
- `id` (path): Milestone ID

**Response:** `204 No Content`

---

## Tasks

### List Tasks for Launch

#### GET /api/v1/launches/:launch_id/tasks

Retrieve all tasks for a specific launch.

**Parameters:**
- `launch_id` (path): Launch ID

**Query Parameters:**
- `milestone_id` (optional): Filter by milestone ID

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "launch_id": 1,
    "milestone_id": 1,
    "title": "Create wireframes",
    "description": "Design wireframes for main pages",
    "assignee_id": 3,
    "status": "done",
    "priority": "high",
    "due_date": "2024-10-10T00:00:00Z",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-10-10T15:00:00Z"
  }
]
```

### Create Task

#### POST /api/v1/tasks

Create a new task.

**Request Body:**
```json
{
  "launch_id": 1,
  "milestone_id": 1,
  "title": "Write user documentation",
  "description": "Create comprehensive user guide",
  "assignee_id": 4,
  "status": "todo",
  "priority": "medium",
  "due_date": "2024-11-15T00:00:00Z"
}
```

**Required Fields:**
- `launch_id`
- `title`

**Response:** `201 Created`

### Get Task

#### GET /api/v1/tasks/:id

Retrieve details of a specific task.

**Parameters:**
- `id` (path): Task ID

**Response:** `200 OK`

### Update Task

#### PUT /api/v1/tasks/:id

Update an existing task.

**Parameters:**
- `id` (path): Task ID

**Request Body:**
```json
{
  "status": "in-progress",
  "assignee_id": 5
}
```

**Response:** `200 OK`

### Delete Task

#### DELETE /api/v1/tasks/:id

Soft delete a task.

**Parameters:**
- `id` (path): Task ID

**Response:** `204 No Content`

---

## Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `204 No Content`: Resource deleted successfully
- `400 Bad Request`: Invalid request data
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error
- `503 Service Unavailable`: Service degraded (e.g., database issues)

## Data Models

### Launch Status Values
- `draft`: Initial state
- `planned`: Launch is planned
- `in-progress`: Launch is in progress
- `launched`: Launch completed
- `cancelled`: Launch cancelled

### Launch Priority Values
- `low`
- `medium`
- `high`
- `critical`

### Milestone Status Values
- `pending`: Not started
- `in-progress`: In progress
- `completed`: Completed
- `blocked`: Blocked

### Task Status Values
- `todo`: Not started
- `in-progress`: In progress
- `done`: Completed
- `blocked`: Blocked

### Task Priority Values
- `low`
- `medium`
- `high`

## Caching

The API uses Redis for caching to improve performance:
- Individual resources are cached for 10 minutes
- List queries are cached for 5 minutes
- Cache is automatically invalidated on updates

## Rate Limiting

Rate limiting will be implemented in a future version.

## Pagination

List endpoints support pagination through `limit` and `offset` query parameters:
- Default limit: 50
- Maximum limit: 100
- Default offset: 0

## Error Handling

All errors return a JSON object with an `error` field:

```json
{
  "error": "Description of the error"
}
```

Common errors:
- `"invalid id"`: Invalid ID format
- `"launch not found"`: Resource not found
- `"failed to create launch"`: Internal error during creation
- `"failed to update launch"`: Internal error during update

## Best Practices

1. **Always handle errors**: Check for error responses and handle them appropriately
2. **Use appropriate HTTP methods**: GET for reads, POST for creates, PUT for updates, DELETE for deletes
3. **Include Content-Type header**: Always set `Content-Type: application/json` for POST/PUT requests
4. **Use pagination**: For large result sets, use `limit` and `offset` parameters
5. **Cache responses**: Implement client-side caching for better performance
6. **Monitor health endpoint**: Regularly check `/health` for service status

## Examples

### Create a Complete Launch with Milestones

```bash
# Step 1: Create a launch
curl -X POST http://localhost:8080/api/v1/launches \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Mobile App Launch",
    "description": "Launch of our new mobile application",
    "launch_date": "2024-12-31T00:00:00Z",
    "status": "planned",
    "priority": "high",
    "tags": ["mobile", "app"]
  }'

# Step 2: Create milestones
curl -X POST http://localhost:8080/api/v1/milestones \
  -H "Content-Type: application/json" \
  -d '{
    "launch_id": 1,
    "title": "Design Phase",
    "due_date": "2024-10-31T00:00:00Z",
    "order": 1
  }'

# Step 3: Create tasks
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "launch_id": 1,
    "milestone_id": 1,
    "title": "Create UI mockups",
    "priority": "high",
    "due_date": "2024-10-15T00:00:00Z"
  }'
```

### Update Launch Status

```bash
curl -X PUT http://localhost:8080/api/v1/launches/1 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "in-progress"
  }'
```

### Query Launches by Status

```bash
curl http://localhost:8080/api/v1/launches?status=in-progress&limit=20
```
