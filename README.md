# Bill Aggregation Service - How To Guide

## Overview
This service allows users to aggregate and manage their utility bills from multiple providers in one place. It supports linking accounts from different utility providers (electricity, water, internet, etc.) and provides a unified view of all bills.

## Prerequisites
- Go 1.x
- Docker and Docker Compose
- PostgreSQL 15
- Redis 7
- Make (optional, but recommended)

## Getting Started

### 1. Clone and Setup
```bash
git clone <repository-url>
cd bill-aggregator
```

### 2. Environment Setup
The service uses the following environment variables (already configured in docker-compose.yml):
- Database:
  - Host: postgres
  - Port: 5432
  - User: postgres
  - Password: postgres
  - Database: bill_aggregator
- Redis:
  - Host: redis
  - Port: 6379

### 3. Running the Service

#### Option 1: Using Docker (Recommended)
```bash
# Build and start all services
make docker-build
make docker-run

# To stop the services
make docker-stop

# To clean up (removes volumes)
make docker-clean
```

#### Option 2: Running Locally
```bash
# Start dependencies (PostgreSQL and Redis)
docker-compose up postgres redis

# Run database migrations
make migrate-up

# Build and run the application
make build
make run
```

### 4. Database Migrations
```bash
# Apply migrations
make migrate-up

# Rollback migrations
make migrate-down
```

## API Usage

### Authentication
All API endpoints require authentication using JWT tokens.

1. Register a new user:
```bash
curl -X POST http://localhost:8081/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "your_password"
  }'
```

2. Login to get a token:
```bash
curl -X POST http://localhost:8081/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "your_password"
  }'
```

### Managing Utility Accounts

1. Link a new utility account:
```bash
curl -X POST http://localhost:8081/accounts/link \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "electricity_provider",
    "credentials": {
      "username": "your_username",
      "password": "your_password"
    }
  }'
```

2. View linked accounts:
```bash
curl -X GET http://localhost:8081/accounts \
  -H "Authorization: Bearer YOUR_TOKEN"
```

3. Delete a linked account:
```bash
curl -X DELETE http://localhost:8081/accounts/{account_id} \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Managing Bills

1. View all bills:
```bash
curl -X GET http://localhost:8081/bills \
  -H "Authorization: Bearer YOUR_TOKEN"
```

2. View bills by provider:
```bash
curl -X GET http://localhost:8081/bills/{provider} \
  -H "Authorization: Bearer YOUR_TOKEN"
```

3. Refresh bills:
```bash
curl -X POST http://localhost:8081/bills/refresh \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Development

### Project Structure
```
.
├── cmd/                    # Application entry points
├── internal/              # Private application code
│   ├── adapters/         # External service adapters
│   ├── domain/           # Business logic and entities
│   ├── ports/            # Interface definitions
│   └── usecases/         # Use case implementations
├── migrations/           # Database migrations
├── docker/              # Docker-related files
└── docs/                # Documentation
```

### Running Tests
```bash
# Run all tests
make test

# Generate mocks for testing
make generate-mocks
```

### Code Quality
```bash
# Run linter
make lint
```

## Monitoring and Maintenance

### Health Checks
The service includes health check endpoints:
```bash
curl http://localhost:8081/health
```

### Logging
Logs are available through Docker:
```bash
docker-compose logs -f app
```

## Troubleshooting

### Common Issues

1. Database Connection Issues
   - Verify PostgreSQL is running: `docker-compose ps postgres`
   - Check database logs: `docker-compose logs postgres`

2. Redis Connection Issues
   - Verify Redis is running: `docker-compose ps redis`
   - Check Redis logs: `docker-compose logs redis`

3. Application Issues
   - Check application logs: `docker-compose logs app`
   - Verify environment variables are set correctly

### Debugging
1. Enable debug logging by setting the appropriate log level
2. Use the health check endpoint to verify service status
3. Check database migrations status

## Security Considerations

1. Always use HTTPS in production
2. Keep your JWT tokens secure
3. Regularly rotate provider credentials
4. Monitor for suspicious activities
5. Keep dependencies updated

## Best Practices

1. Regular Updates
   - Refresh bills regularly to ensure up-to-date information
   - Keep the application and dependencies updated

2. Account Management
   - Use strong passwords
   - Enable two-factor authentication when available
   - Regularly review linked accounts

3. Data Management
   - Regularly backup your data
   - Monitor bill amounts for unusual changes
   - Set up notifications for important events

## Support

For issues and support:
1. Check the documentation
2. Review the troubleshooting guide
3. Contact the development team
4. Check the service status page