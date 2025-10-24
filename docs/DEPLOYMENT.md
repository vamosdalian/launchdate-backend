# Deployment Guide

This guide covers various deployment options for the LaunchDate backend.

## Table of Contents

1. [Local Development](#local-development)
2. [Docker Deployment](#docker-deployment)
3. [Production Deployment](#production-deployment)
4. [Database Migrations](#database-migrations)
5. [Environment Variables](#environment-variables)
6. [Monitoring](#monitoring)

## Local Development

### Prerequisites

- Go 1.21 or later
- PostgreSQL 15
- Redis 7
- Docker (optional)

### Quick Start

1. Clone the repository:
```bash
git clone https://github.com/vamosdalian/launchdate-backend.git
cd launchdate-backend
```

2. Start dependencies with Docker:
```bash
docker run -d --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:15-alpine
docker run -d --name redis -p 6379:6379 redis:7-alpine
```

3. Run migrations:
```bash
make migrate-up
```

4. Start the server:
```bash
make run
```

The API will be available at http://localhost:8080

## Docker Deployment

### Using Docker Compose

The simplest way to deploy is using Docker Compose:

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

This will start:
- PostgreSQL database
- Redis cache
- Database migrations
- Application server

### Building Docker Image

```bash
# Build the image
docker build -t launchdate-backend:latest .

# Run the container
docker run -d \
  --name launchdate-backend \
  -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-password \
  -e REDIS_HOST=your-redis-host \
  launchdate-backend:latest
```

### Using Pre-built Image from GHCR

```bash
# Pull the latest image
docker pull ghcr.io/vamosdalian/launchdate-backend:latest

# Run the container
docker run -d \
  --name launchdate-backend \
  -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-password \
  -e REDIS_HOST=your-redis-host \
  ghcr.io/vamosdalian/launchdate-backend:latest
```

## Production Deployment

### Requirements

- PostgreSQL 15+ (managed service recommended)
- Redis 7+ (managed service recommended)
- Container orchestration platform (Kubernetes, ECS, etc.)
- Load balancer (optional, for high availability)

### Kubernetes Deployment

Create a Kubernetes deployment:

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: launchdate-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: launchdate-backend
  template:
    metadata:
      labels:
        app: launchdate-backend
    spec:
      containers:
      - name: launchdate-backend
        image: ghcr.io/vamosdalian/launchdate-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: db-secrets
              key: host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secrets
              key: password
        - name: REDIS_HOST
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: redis-host
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: launchdate-backend
spec:
  selector:
    app: launchdate-backend
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

Deploy:
```bash
kubectl apply -f deployment.yaml
```

### AWS ECS Deployment

Example task definition:

```json
{
  "family": "launchdate-backend",
  "containerDefinitions": [
    {
      "name": "launchdate-backend",
      "image": "ghcr.io/vamosdalian/launchdate-backend:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "ENVIRONMENT",
          "value": "production"
        },
        {
          "name": "DB_HOST",
          "value": "your-rds-endpoint.amazonaws.com"
        },
        {
          "name": "REDIS_HOST",
          "value": "your-elasticache-endpoint.amazonaws.com"
        }
      ],
      "secrets": [
        {
          "name": "DB_PASSWORD",
          "valueFrom": "arn:aws:secretsmanager:region:account:secret:db-password"
        }
      ],
      "healthCheck": {
        "command": ["CMD-SHELL", "curl -f http://localhost:8080/health || exit 1"],
        "interval": 30,
        "timeout": 5,
        "retries": 3
      },
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/launchdate-backend",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ],
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512"
}
```

## Database Migrations

### Running Migrations

Using migrate CLI:
```bash
# Install migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations \
  -database "postgres://user:password@host:5432/dbname?sslmode=disable" \
  up

# Rollback migrations
migrate -path migrations \
  -database "postgres://user:password@host:5432/dbname?sslmode=disable" \
  down 1
```

Using Docker:
```bash
docker run -v $(pwd)/migrations:/migrations \
  migrate/migrate \
  -path=/migrations/ \
  -database "postgres://user:password@host:5432/dbname?sslmode=disable" \
  up
```

### Creating New Migrations

```bash
make migrate-create NAME=add_user_table
```

This creates two files:
- `migrations/XXX_add_user_table.up.sql`
- `migrations/XXX_add_user_table.down.sql`

## Environment Variables

### Required Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=launchdate
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
```

### Optional Variables

```bash
# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
ENVIRONMENT=production

# Redis
REDIS_PASSWORD=
REDIS_DB=0

# S3/MinIO (for file uploads)
S3_ENDPOINT=s3.amazonaws.com
S3_ACCESS_KEY_ID=
S3_SECRET_ACCESS_KEY=
S3_BUCKET_NAME=launchdate
S3_USE_SSL=true
```

### Using Secrets

For production, use secrets management:

**Kubernetes Secrets:**
```bash
kubectl create secret generic db-secrets \
  --from-literal=host=db.example.com \
  --from-literal=password=secret
```

**AWS Secrets Manager:**
```bash
aws secretsmanager create-secret \
  --name launchdate/db \
  --secret-string '{"password":"secret","host":"db.example.com"}'
```

## Monitoring

### Health Checks

The application provides a health check endpoint:
```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "ok",
  "database": "ok",
  "redis": "ok"
}
```

### Logging

The application logs in JSON format to stdout:
```json
{
  "level": "info",
  "msg": "request processed",
  "status": 200,
  "method": "GET",
  "path": "/api/v1/launches",
  "latency": "15ms",
  "time": "2024-01-20T10:00:00Z"
}
```

### Metrics

For production monitoring, integrate with:
- **Prometheus**: Add metrics endpoint
- **Grafana**: Visualize metrics
- **CloudWatch**: AWS monitoring
- **Datadog**: Application performance monitoring

### Setting up Prometheus

Example prometheus.yml:
```yaml
scrape_configs:
  - job_name: 'launchdate-backend'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
```

## High Availability

### Database

Use managed PostgreSQL services:
- **AWS RDS**: With Multi-AZ deployment
- **Google Cloud SQL**: With high availability
- **Azure Database**: With zone redundancy

Configure read replicas for read-heavy workloads.

### Redis

Use managed Redis services:
- **AWS ElastiCache**: With automatic failover
- **Redis Enterprise Cloud**
- **Azure Cache for Redis**

### Application

Deploy multiple instances:
- Use container orchestration (Kubernetes, ECS)
- Configure auto-scaling based on CPU/memory
- Use health checks for automatic recovery
- Implement graceful shutdown

## Backup and Recovery

### Database Backups

Automated backups:
```bash
# Backup
pg_dump -h localhost -U postgres launchdate > backup.sql

# Restore
psql -h localhost -U postgres launchdate < backup.sql
```

Schedule daily backups with cron or managed service features.

### Redis Persistence

Configure Redis persistence:
```bash
# In redis.conf
save 900 1
save 300 10
save 60 10000
```

## SSL/TLS Configuration

### Using Reverse Proxy (Recommended)

Use nginx or similar:
```nginx
server {
    listen 443 ssl;
    server_name api.launch-date.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Database SSL

Enable SSL for database connection:
```bash
DB_SSLMODE=require
```

## Troubleshooting

### Common Issues

**Database connection failed:**
- Check DB_HOST and credentials
- Verify network connectivity
- Check PostgreSQL is running
- Verify firewall rules

**Redis connection failed:**
- Check REDIS_HOST
- Verify Redis is running
- Check if password is required

**Application not starting:**
- Check logs: `docker-compose logs app`
- Verify all required environment variables
- Check database migrations have run

### Debug Mode

For debugging, set log level to debug:
```go
logger.SetLevel(logrus.DebugLevel)
```

## Performance Tuning

### Database

- Add indexes for frequently queried fields
- Adjust connection pool size
- Enable query caching
- Use read replicas

### Redis

- Adjust cache TTL based on data volatility
- Use Redis cluster for high throughput
- Monitor memory usage

### Application

- Adjust worker pool size
- Tune timeout values
- Enable HTTP/2
- Use connection pooling

## Security Best Practices

1. **Use secrets management** for sensitive data
2. **Enable SSL/TLS** for all connections
3. **Implement rate limiting** to prevent abuse
4. **Use prepared statements** to prevent SQL injection
5. **Keep dependencies updated** regularly
6. **Monitor security advisories** for Go and dependencies
7. **Use least privilege** for database users
8. **Enable audit logging** in production
9. **Implement authentication** before public release
10. **Use HTTPS only** in production

## CI/CD Pipeline

The repository includes GitHub Actions workflows:

- **On Pull Request**: Run tests and lint
- **On Push to Main**: Build and push Docker image
- **On Tag**: Create release and deploy

To trigger a deployment:
```bash
git tag v1.0.0
git push origin v1.0.0
```

## Rollback

### Kubernetes
```bash
kubectl rollout undo deployment/launchdate-backend
```

### ECS
```bash
aws ecs update-service \
  --cluster my-cluster \
  --service launchdate-backend \
  --task-definition launchdate-backend:previous-version
```

### Database Migrations
```bash
migrate -path migrations \
  -database "postgres://..." \
  down 1
```

## Support

For issues and questions:
- GitHub Issues: https://github.com/vamosdalian/launchdate-backend/issues
- Documentation: See docs/ directory
