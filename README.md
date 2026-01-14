# Go API - Hexagonal Architecture

A production-ready REST API built with Go using Hexagonal Architecture (Ports & Adapters) pattern.

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- MinIO (optional, for file storage)

### Running Locally

```bash
# Clone the repository
git clone https://github.com/raulaguila/go-api.git
cd go-api/new_api

# Copy environment file
cp config/.env.example config/.env

# Edit .env with your database credentials
# Then run:
make run

# Or without make:
go run cmd/backend/main.go
```

### Running with Docker

```bash
docker-compose up -d
```

## 📁 Project Structure

```
new_api/
├── cmd/backend/          # Application entry point
├── config/               # Configuration files
├── internal/
│   ├── adapter/
│   │   ├── driven/       # Output adapters (DB, Storage)
│   │   └── driver/       # Input adapters (REST API)
│   ├── app/              # Application layer
│   ├── core/
│   │   ├── domain/       # Entities and business rules
│   │   ├── dto/          # Data Transfer Objects
│   │   ├── port/         # Interfaces (input/output)
│   │   └── usecase/      # Business logic
│   └── di/               # Dependency injection
└── pkg/                  # Shared packages
    ├── apperror/         # Centralized error handling
    ├── logger/           # Structured logging
    ├── utils/            # Helper functions
    └── validator/        # Validation utilities
```

## 🔧 Configuration

Environment variables (see `config/.env`):

| Variable | Description | Default |
|----------|-------------|---------|
| `API_PORT` | Server port | `9000` |
| `ENVIRONMENT` | development/production | `development` |
| `LOG_LEVEL` | debug/info/warn/error | `info` |
| `POSTGRES_HOST` | Database host | `localhost` |
| `POSTGRES_PORT` | Database port | `5432` |
| `POSTGRES_USER` | Database user | `root` |
| `POSTGRES_PASS` | Database password | `root` |
| `POSTGRES_BASE` | Database name | `api` |

## 📚 API Documentation

Swagger UI is available at: `http://localhost:9000/swagger/`

### Main Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/auth` | Login |
| `GET` | `/auth` | Get current user |
| `PUT` | `/auth` | Refresh token |
| `GET` | `/user` | List users |
| `POST` | `/user` | Create user |
| `GET` | `/user/:id` | Get user by ID |
| `PUT` | `/user/:id` | Update user |
| `DELETE` | `/user` | Delete users |
| `GET` | `/profile` | List profiles |
| `POST` | `/profile` | Create profile |
| `GET` | `/health` | Health check |

## 🏗️ Architecture

This project follows **Hexagonal Architecture** (Ports & Adapters):

```
┌─────────────────────────────────────────────────────────────┐
│                      Driver Adapters                         │
│                  (REST, gRPC, CLI, GraphQL)                 │
├─────────────────────────────────────────────────────────────┤
│                       Input Ports                            │
│                    (Use Case Interfaces)                     │
├─────────────────────────────────────────────────────────────┤
│                      Application Core                        │
│              (Entities, Use Cases, Domain Logic)            │
├─────────────────────────────────────────────────────────────┤
│                      Output Ports                            │
│                  (Repository Interfaces)                     │
├─────────────────────────────────────────────────────────────┤
│                     Driven Adapters                          │
│               (PostgreSQL, MinIO, External APIs)             │
└─────────────────────────────────────────────────────────────┘
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package tests
go test ./internal/core/usecase/...
```

## 🛠️ Development

### Available Make Commands

```bash
make run          # Run the application
make build        # Build the binary
make test         # Run tests
make lint         # Run linter
make swagger      # Generate Swagger docs
make docker-up    # Start with Docker
make docker-down  # Stop Docker containers
```

### Adding New Features

See [ADDING_FEATURES.md](./ADDING_FEATURES.md) for a complete guide on how to add new functionality following the hexagonal architecture.

### Architecture Details

See [ARCHITECTURE.md](./ARCHITECTURE.md) for detailed documentation on the project architecture.

## 📄 License

MIT License - see LICENSE file for details.

## 👥 Contributors

- [Raul del Aguila](https://github.com/raulaguila)
