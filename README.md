# LaunchDate Backend

Backend API for the LaunchDate launch management platform.

## Features

- 🚀 **Launch Management**: Create and manage product/project launches
- 📍 **Milestone Tracking**: Track key milestones for each launch
- ✅ **Task Management**: Manage tasks associated with launches and milestones
- 🗃️ **PostgreSQL Database**: Robust data persistence with migrations
- ⚡ **Redis Caching**: High-performance caching layer
- 🔒 **Health Checks**: Monitor service health
- 📚 **OpenAPI Documentation**: Complete API documentation in docs/
- 🐳 **Docker Support**: Containerized deployment
- 🔄 **CI/CD**: Automated testing and deployment to GHCR

## Tech Stack

- **Language**: Go 1.21
- **Web Framework**: Gin
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **ORM**: sqlx
- **Documentation**: OpenAPI 3.0

## Getting Started

### Prerequisites

- Go 1.21 or later
- PostgreSQL 15
- Redis 7
- Docker & Docker Compose (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/vamosdalian/launchdate-backend.git
cd launchdate-backend
```

2. Copy the environment file:
```bash
cp .env.example .env
```

3. Update the `.env` file with your configuration.

### Running with Docker Compose (Recommended)

The easiest way to run the application is with Docker Compose:

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database
- Redis cache
- Database migrations
- Application server on http://localhost:8080

### Running Locally

1. Install dependencies:
```bash
go mod download
```

2. Start PostgreSQL and Redis:
```bash
# Using Docker
docker run -d --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:15-alpine
docker run -d --name redis -p 6379:6379 redis:7-alpine
```

3. Run migrations:
```bash
# Install migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/launchdate?sslmode=disable" up
```

4. Start the server:
```bash
go run cmd/server/main.go
```

The server will start on http://localhost:8080

## API Documentation

The complete API documentation is available in OpenAPI 3.0 format:
- [OpenAPI Specification](docs/openapi.yaml)

You can view the documentation using tools like:
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
- [Redoc](https://redocly.github.io/redoc/)
- [Postman](https://www.postman.com/)

### Key Endpoints

#### Health Check
- `GET /health` - Check service health

#### Launches
- `GET /api/v1/launches` - List all launches
- `POST /api/v1/launches` - Create a new launch
- `GET /api/v1/launches/{id}` - Get launch details
- `PUT /api/v1/launches/{id}` - Update a launch
- `DELETE /api/v1/launches/{id}` - Delete a launch

#### Milestones
- `GET /api/v1/launches/{launch_id}/milestones` - List milestones for a launch
- `POST /api/v1/milestones` - Create a new milestone
- `GET /api/v1/milestones/{id}` - Get milestone details
- `PUT /api/v1/milestones/{id}` - Update a milestone
- `DELETE /api/v1/milestones/{id}` - Delete a milestone

#### Tasks
- `GET /api/v1/launches/{launch_id}/tasks` - List tasks for a launch
- `POST /api/v1/tasks` - Create a new task
- `GET /api/v1/tasks/{id}` - Get task details
- `PUT /api/v1/tasks/{id}` - Update a task
- `DELETE /api/v1/tasks/{id}` - Delete a task

## Project Structure

```
.
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── api/             # HTTP handlers and routing
│   ├── config/          # Configuration management
│   ├── database/        # Database connection
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Domain models
│   ├── repository/      # Data access layer
│   └── service/         # Business logic layer
├── migrations/          # Database migrations
├── docs/                # API documentation
├── .github/
│   └── workflows/       # CI/CD workflows
├── Dockerfile           # Docker image definition
├── docker-compose.yml   # Docker Compose configuration
└── README.md
```

## Development

### Running Tests

```bash
go test -v ./...
```

### Running with Coverage

```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Linting

```bash
golangci-lint run
```

## Deployment

The application is automatically built and deployed to GitHub Container Registry (GHCR) on push to main/develop branches.

### Pull the Docker Image

```bash
docker pull ghcr.io/vamosdalian/launchdate-backend:latest
```

### Run the Docker Image

```bash
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-password \
  -e REDIS_HOST=your-redis-host \
  ghcr.io/vamosdalian/launchdate-backend:latest
```

## Environment Variables

See [.env.example](.env.example) for all available configuration options.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Contact

Project Link: [https://github.com/vamosdalian/launchdate-backend](https://github.com/vamosdalian/launchdate-backend)
