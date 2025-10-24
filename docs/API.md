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

---

## Rocket Launch Tracking APIs

The following APIs support rocket launch tracking and space industry data management.

## Companies

Space companies and organizations.

### List Companies

#### GET /api/v1/companies

Retrieve a list of space companies.

**Query Parameters:**
- `limit` (optional): Maximum number of results (default: 50, max: 100)
- `offset` (optional): Number of results to skip (default: 0)

**Example Request:**
```bash
curl http://localhost:8080/api/v1/companies?limit=10
```

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "name": "SpaceX",
    "description": "Space Exploration Technologies Corp.",
    "founded": 2002,
    "founder": "Elon Musk",
    "headquarters": "Hawthorne, California",
    "employees": 12000,
    "website": "https://www.spacex.com",
    "imageUrl": "https://example.com/spacex.jpg",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
]
```

### Create Company

#### POST /api/v1/companies

Create a new space company.

**Request Body:**
```json
{
  "name": "Blue Origin",
  "description": "Private aerospace manufacturer",
  "founded": 2000,
  "founder": "Jeff Bezos",
  "headquarters": "Kent, Washington",
  "employees": 3500,
  "website": "https://www.blueorigin.com",
  "imageUrl": "https://example.com/blueorigin.jpg"
}
```

**Required Fields:**
- `name`

**Response:** `201 Created`
```json
{
  "id": 2,
  "name": "Blue Origin",
  "description": "Private aerospace manufacturer",
  "founded": 2000,
  "founder": "Jeff Bezos",
  "headquarters": "Kent, Washington",
  "employees": 3500,
  "website": "https://www.blueorigin.com",
  "imageUrl": "https://example.com/blueorigin.jpg",
  "created_at": "2024-10-24T10:00:00Z",
  "updated_at": "2024-10-24T10:00:00Z"
}
```

### Get Company

#### GET /api/v1/companies/:id

Get details of a specific company.

**Parameters:**
- `id` (path): Company ID

**Response:** `200 OK`

### Update Company

#### PUT /api/v1/companies/:id

Update an existing company.

**Parameters:**
- `id` (path): Company ID

**Request Body:**
All fields are optional. Only include fields you want to update.

```json
{
  "employees": 13000,
  "description": "Updated description"
}
```

**Response:** `200 OK`

### Delete Company

#### DELETE /api/v1/companies/:id

Soft delete a company.

**Parameters:**
- `id` (path): Company ID

**Response:** `204 No Content`

---

## Rockets

Rocket specifications and details.

### List Rockets

#### GET /api/v1/rockets

List all rockets with optional filtering.

**Query Parameters:**
- `active` (optional): Filter by active status (true/false)
- `company_id` (optional): Filter by company ID
- `limit` (optional): Maximum number of results (default: 50, max: 100)
- `offset` (optional): Number of results to skip (default: 0)

**Example Request:**
```bash
curl http://localhost:8080/api/v1/rockets?active=true&limit=10
```

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "name": "Falcon 9",
    "description": "Two-stage reusable rocket",
    "height": 70.0,
    "diameter": 3.7,
    "mass": 549054.0,
    "company_id": 1,
    "company": "SpaceX",
    "imageUrl": "https://example.com/falcon9.jpg",
    "active": true,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
]
```

### Create Rocket

#### POST /api/v1/rockets

Create a new rocket.

**Request Body:**
```json
{
  "name": "Starship",
  "description": "Fully reusable super heavy-lift launch vehicle",
  "height": 120.0,
  "diameter": 9.0,
  "mass": 5000000.0,
  "company_id": 1,
  "imageUrl": "https://example.com/starship.jpg",
  "active": true
}
```

**Required Fields:**
- `name`

**Response:** `201 Created`

### Get Rocket

#### GET /api/v1/rockets/:id

Get details of a specific rocket.

**Parameters:**
- `id` (path): Rocket ID

**Response:** `200 OK`

### Update Rocket

#### PUT /api/v1/rockets/:id

Update an existing rocket.

**Parameters:**
- `id` (path): Rocket ID

**Request Body:**
```json
{
  "active": false,
  "description": "Updated description"
}
```

**Response:** `200 OK`

### Delete Rocket

#### DELETE /api/v1/rockets/:id

Soft delete a rocket.

**Parameters:**
- `id` (path): Rocket ID

**Response:** `204 No Content`

---

## Launch Bases

Launch sites and facilities worldwide.

### List Launch Bases

#### GET /api/v1/launch-bases

List all launch sites.

**Query Parameters:**
- `country` (optional): Filter by country
- `limit` (optional): Maximum number of results (default: 50, max: 100)
- `offset` (optional): Number of results to skip (default: 0)

**Example Request:**
```bash
curl http://localhost:8080/api/v1/launch-bases?country=USA
```

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "name": "Kennedy Space Center",
    "location": "Merritt Island, Florida",
    "country": "USA",
    "description": "NASA's primary launch center",
    "imageUrl": "https://example.com/ksc.jpg",
    "latitude": 28.573469,
    "longitude": -80.651070,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
]
```

### Create Launch Base

#### POST /api/v1/launch-bases

Create a new launch base.

**Request Body:**
```json
{
  "name": "Vandenberg Space Force Base",
  "location": "Santa Barbara County, California",
  "country": "USA",
  "description": "U.S. Space Force Base for polar orbit launches",
  "imageUrl": "https://example.com/vandenberg.jpg",
  "latitude": 34.632,
  "longitude": -120.611
}
```

**Required Fields:**
- `name`

**Response:** `201 Created`

### Get Launch Base

#### GET /api/v1/launch-bases/:id

Get details of a specific launch base.

**Parameters:**
- `id` (path): Launch base ID

**Response:** `200 OK`

### Update Launch Base

#### PUT /api/v1/launch-bases/:id

Update an existing launch base.

**Parameters:**
- `id` (path): Launch base ID

**Request Body:**
```json
{
  "description": "Updated description",
  "latitude": 34.633,
  "longitude": -120.612
}
```

**Response:** `200 OK`

### Delete Launch Base

#### DELETE /api/v1/launch-bases/:id

Soft delete a launch base.

**Parameters:**
- `id` (path): Launch base ID

**Response:** `204 No Content`

---

## Rocket Launches

Scheduled and historical rocket launch events.

### List Rocket Launches

#### GET /api/v1/rocket-launches

List all rocket launches with optional filtering.

**Query Parameters:**
- `status` (optional): Filter by status (`scheduled`, `successful`, `failed`, `cancelled`)
- `rocket_id` (optional): Filter by rocket ID
- `launch_base_id` (optional): Filter by launch base ID
- `limit` (optional): Maximum number of results (default: 50, max: 100)
- `offset` (optional): Number of results to skip (default: 0)

**Example Request:**
```bash
curl http://localhost:8080/api/v1/rocket-launches?status=scheduled&limit=20
```

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "name": "Starlink 6-35",
    "date": "2025-10-25T10:30:00Z",
    "rocket_id": 1,
    "rocket": "Falcon 9",
    "launch_base_id": 1,
    "launchBase": "Kennedy Space Center",
    "status": "scheduled",
    "description": "Deployment of 23 Starlink satellites",
    "created_at": "2024-10-01T10:00:00Z",
    "updated_at": "2024-10-01T10:00:00Z"
  }
]
```

### Create Rocket Launch

#### POST /api/v1/rocket-launches

Create a new rocket launch event.

**Request Body:**
```json
{
  "name": "Artemis III",
  "date": "2026-09-01T14:00:00Z",
  "rocket_id": 5,
  "launch_base_id": 1,
  "status": "scheduled",
  "description": "First crewed Moon landing of the Artemis program"
}
```

**Required Fields:**
- `name`
- `date`

**Response:** `201 Created`

### Get Rocket Launch

#### GET /api/v1/rocket-launches/:id

Get details of a specific rocket launch.

**Parameters:**
- `id` (path): Rocket launch ID

**Response:** `200 OK`

### Update Rocket Launch

#### PUT /api/v1/rocket-launches/:id

Update an existing rocket launch.

**Parameters:**
- `id` (path): Rocket launch ID

**Request Body:**
```json
{
  "status": "successful",
  "date": "2026-09-02T15:30:00Z"
}
```

**Response:** `200 OK`

### Delete Rocket Launch

#### DELETE /api/v1/rocket-launches/:id

Soft delete a rocket launch.

**Parameters:**
- `id` (path): Rocket launch ID

**Response:** `204 No Content`

### Rocket Launch Status Values
- `scheduled`: Upcoming launch
- `successful`: Launch completed successfully
- `failed`: Launch failed
- `cancelled`: Launch cancelled

---

## News

Space industry news and updates.

### List News

#### GET /api/v1/news

List all news articles.

**Query Parameters:**
- `limit` (optional): Maximum number of results (default: 50, max: 100)
- `offset` (optional): Number of results to skip (default: 0)

**Example Request:**
```bash
curl http://localhost:8080/api/v1/news?limit=10
```

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "title": "SpaceX Completes 200th Successful Landing",
    "summary": "SpaceX has achieved another milestone with its 200th successful booster landing",
    "content": "Full article content in markdown format...",
    "date": "2025-10-20T12:00:00Z",
    "url": "https://www.spacex.com/updates",
    "imageUrl": "https://example.com/news1.jpg",
    "created_at": "2025-10-20T12:00:00Z",
    "updated_at": "2025-10-20T12:00:00Z"
  }
]
```

### Create News Article

#### POST /api/v1/news

Create a new news article.

**Request Body:**
```json
{
  "title": "NASA Announces New Moon Mission",
  "summary": "NASA reveals plans for additional lunar missions",
  "content": "Detailed article content with markdown support...",
  "date": "2025-10-24T08:00:00Z",
  "url": "https://www.nasa.gov/news",
  "imageUrl": "https://example.com/moon-mission.jpg"
}
```

**Required Fields:**
- `title`
- `date`

**Response:** `201 Created`

### Get News Article

#### GET /api/v1/news/:id

Get details of a specific news article.

**Parameters:**
- `id` (path): News article ID

**Response:** `200 OK`

### Update News Article

#### PUT /api/v1/news/:id

Update an existing news article.

**Parameters:**
- `id` (path): News article ID

**Request Body:**
```json
{
  "title": "Updated Title",
  "summary": "Updated summary"
}
```

**Response:** `200 OK`

### Delete News Article

#### DELETE /api/v1/news/:id

Soft delete a news article.

**Parameters:**
- `id` (path): News article ID

**Response:** `204 No Content`

---

## Complete API Example: Rocket Launch Workflow

Here's an example workflow for creating a complete rocket launch with all related entities:

```bash
# Step 1: Create a space company
curl -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -d '{
    "name": "SpaceX",
    "description": "Space Exploration Technologies Corp.",
    "founded": 2002,
    "founder": "Elon Musk",
    "headquarters": "Hawthorne, California",
    "website": "https://www.spacex.com"
  }'

# Step 2: Create a rocket
curl -X POST http://localhost:8080/api/v1/rockets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Falcon 9",
    "description": "Two-stage reusable rocket",
    "height": 70.0,
    "diameter": 3.7,
    "mass": 549054.0,
    "company_id": 1,
    "active": true
  }'

# Step 3: Create a launch base
curl -X POST http://localhost:8080/api/v1/launch-bases \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Kennedy Space Center",
    "location": "Merritt Island, Florida",
    "country": "USA",
    "latitude": 28.573469,
    "longitude": -80.651070
  }'

# Step 4: Create a rocket launch event
curl -X POST http://localhost:8080/api/v1/rocket-launches \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Starlink 6-35",
    "date": "2025-10-25T10:30:00Z",
    "rocket_id": 1,
    "launch_base_id": 1,
    "status": "scheduled",
    "description": "Deployment of 23 Starlink satellites"
  }'

# Step 5: List all scheduled launches
curl http://localhost:8080/api/v1/rocket-launches?status=scheduled
```

