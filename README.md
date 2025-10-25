# LaunchDate Backend

Backend API for the LaunchDate platform - rocket launch tracking and space news.

## Features

### Rocket Launch Tracking
- 🚁 **Company Management**: Track space companies and their details
- 🚀 **Rocket Database**: Maintain comprehensive rocket specifications
- 🌍 **Launch Bases**: Manage launch sites worldwide with geo-coordinates
- 📅 **Rocket Launch Events**: Track scheduled and historical rocket launches
- 📰 **News Management**: Space news and updates

### Infrastructure
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

#### Rocket Launch Tracking

**Companies**
- `GET /api/v1/companies` - List space companies
- `POST /api/v1/companies` - Create a company
- `GET /api/v1/companies/{id}` - Get company details
- `PUT /api/v1/companies/{id}` - Update a company
- `DELETE /api/v1/companies/{id}` - Delete a company

**Rockets**
- `GET /api/v1/rockets` - List rockets
- `POST /api/v1/rockets` - Create a rocket
- `GET /api/v1/rockets/{id}` - Get rocket details
- `PUT /api/v1/rockets/{id}` - Update a rocket
- `DELETE /api/v1/rockets/{id}` - Delete a rocket

**Launch Bases**
- `GET /api/v1/launch-bases` - List launch sites
- `POST /api/v1/launch-bases` - Create a launch base
- `GET /api/v1/launch-bases/{id}` - Get launch base details
- `PUT /api/v1/launch-bases/{id}` - Update a launch base
- `DELETE /api/v1/launch-bases/{id}` - Delete a launch base

**Rocket Launches**
- `GET /api/v1/rocket-launches` - List rocket launches
- `POST /api/v1/rocket-launches` - Create a rocket launch
- `GET /api/v1/rocket-launches/{id}` - Get rocket launch details
- `PUT /api/v1/rocket-launches/{id}` - Update a rocket launch
- `DELETE /api/v1/rocket-launches/{id}` - Delete a rocket launch

**News**
- `GET /api/v1/news` - List news articles
- `POST /api/v1/news` - Create a news article
- `GET /api/v1/news/{id}` - Get news article details
- `PUT /api/v1/news/{id}` - Update a news article
- `DELETE /api/v1/news/{id}` - Delete a news article

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
