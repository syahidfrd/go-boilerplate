# GO Boilerplate

[![Go Reference](https://pkg.go.dev/badge/github.com/syahidfrd/go-boilerplate.svg)](https://pkg.go.dev/github.com/syahidfrd/go-boilerplate)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/syahidfrd/go-boilerplate)](https://goreportcard.com/report/github.com/syahidfrd/go-boilerplate)

## Changelog

- **v1**: checkout to the [v1 branch](https://github.com/syahidfrd/go-boilerplate/tree/v1) <br>
  Function-based layer architecture with separate layers (domain, usecase, delivery, repository)

## Features

- 🏗️ **Modular Architecture** - Clean separation of concerns with store, service, and handler layers
- 🔐 **JWT Authentication** - Secure user authentication with middleware
- 🏥 **Health Checks** - Database and cache connectivity monitoring
- 📝 **Request Logging** - Comprehensive HTTP request logging
- ⚡ **Caching** - Redis integration for performance optimization
- 🧪 **Testing** - Unit and integration test coverage
- 🛡️ **Middleware** - CORS, recovery, real IP extraction, and request ID tracking
- 📊 **Structured Logging** - Using zerolog for structured, contextual logging
- 🔧 **Configuration** - Environment-based configuration management

### Technology Stack

- [Go 1.22+](https://golang.org/) - Programming language
- [GORM](https://gorm.io/) - ORM library for database operations
- [PostgreSQL](https://www.postgresql.org/) - Primary database
- [Redis](https://redis.io/) - Caching and session storage
- [JWT](https://github.com/golang-jwt/jwt) - Authentication tokens
- [Zerolog](https://github.com/rs/zerolog) - Structured logging
- [Testify](https://github.com/stretchr/testify) - Testing framework
- [Testcontainers](https://testcontainers.com/) - Integration testing with real containers
- [go-playground/validator](https://github.com/go-playground/validator) - Request validation
- [go-redis](https://github.com/go-redis/redis) - Redis client
- [godotenv](https://github.com/joho/godotenv) - Environment variables loader

## Quick Start

### Prerequisites

- Go 1.22+
- Docker (for databases and integration tests)
- PostgreSQL (optional if using Docker)
- Redis (optional if using Docker)

### Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/syahidfrd/go-boilerplate.git
   cd go-boilerplate
   ```

2. **Install dependencies and format code**

   ```bash
   make tidy
   ```

3. **Setup databases**

   ```bash
   # PostgreSQL
   docker run -d -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=go_boilerplate postgres

   # Redis
   docker run -d -p 6379:6379 redis
   ```

4. **Configure environment**

   ```bash
   cp .env.example .env
   # Edit .env file with your configuration
   ```

5. **Run tests**

   ```bash
   # Unit tests
   make test/unit

   # Integration tests
   make test/integration

   # All tests
   make test/all
   ```

6. **Run the application**

   ```bash
   # Option 1: Normal run
   make run bin=server

   # Option 2: Hot reload (development)
   make run/live bin=server
   ```

## API Endpoints

### Authentication

- `POST /api/auth/signup` - User registration
- `POST /api/auth/signin` - User login

### Todos (Protected)

- `GET /api/todos` - Get user's todos
- `POST /api/todos` - Create new todo
- `GET /api/todos/{id}` - Get specific todo
- `PUT /api/todos/{id}` - Update todo
- `PATCH /api/todos/{id}/toggle` - Toggle completion status
- `DELETE /api/todos/{id}` - Delete todo

### Health Check

- `GET /health` - Service health status

## Testing

This project includes comprehensive testing at multiple levels:

### Unit Tests

Run unit tests:

```bash
make test/unit
```

### Integration Tests

Run integration tests with real PostgreSQL and Redis containers:

```bash
make test/integration
```

**Requirements:**

- Docker must be running
- Integration tests use [Testcontainers](https://testcontainers.com/) to spin up real databases

### All Tests

```bash
make test/all
```

## Makefile Commands

You can find all available commands in the `Makefile`:

### Development

- `make tidy` - Format code and tidy dependencies
- `make build` - Build the application binary
- `make run bin=server` - Build and run the application
- `make run/live bin=server` - Run with live reload using Air

### Testing

- `make test/unit` - Run unit tests only
- `make test/integration` - Run integration tests (requires Docker)
- `make test/all` - Run all tests (unit + integration)
