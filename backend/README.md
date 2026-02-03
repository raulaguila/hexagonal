# Backend API

This repository contains the backend service for the methods application, built with **Go** and designed using **Hexagonal Architecture (Ports and Adapters)** to ensure scalability, maintainability, and testability.

## 🚀 Tech Stack

-   **Language:** [Go 1.25+](https://go.dev/)
-   **Framework:** [Fiber v2](https://gofiber.io/) (High performance web framework)
-   **Database:**
    -   [PostgreSQL](https://www.postgresql.org/) (Primary Data Store)
    -   [Redis](https://redis.io/) (Caching)
-   **ORM:** [GORM](https://gorm.io/) (with `pgx` driver)
-   **Authentication:** JWT (Access & Refresh Tokens)
-   **Documentation:** [Swagger](https://swagger.io/) (Swag)
-   **Observability:** [OpenTelemetry](https://opentelemetry.io/) & [ELK Stack](https://www.elastic.co/)
-   **Storage:** [MinIO](https://min.io/) (S3 Compatible Object Storage)
-   **Testing:** [Testcontainers](https://testcontainers.com/)

## 🏗 Architecture

The project follows the **Hexagonal Architecture** pattern:

```
backend/
├── cmd/                # Application entry points
├── config/             # Configuration files (.env, env.sh)
├── internal/
│   ├── core/           # Business logic (Domain & Ports)
│   │   ├── domain/     # Enterprise entities
│   │   ├── port/       # Interfaces (Input/Output ports)
│   │   └── service/    # Application services (Use cases)
│   └── adapter/        # Implementation details (Adapters)
│       ├── driven/     # Infrastructure (DB, Redis, Logger)
│       └── driver/     # Entry points (REST Handler, CLI)
├── pkg/                # Shared libraries / utilities
└── build/              # Dockerfiles and Migration scripts
```

## 🛠 Prerequisites

Ensure you have the following installed:

-   **Go** (v1.25 or higher)
-   **Docker** & **Docker Compose**
-   **Make** (Optional, but recommended for running commands)

## ⚡ Getting Started

### 1. Setup Environment

Initialize the configuration file:

```bash
make init
```

This will execute permissions and set up your `.env` file in the `config/` directory.

### 2. Configuration (`.env`)

Verify the `config/.env` file. Key variables include:

-   `API_PORT`: Port for the backend service (default: `9999`)
-   `POSTGRES_*`: Database credentials.
-   `REDIS_*`: Redis credentials.
-   `MINIO_*`: Object storage settings.
-   `JWT_*`: Token definition and expiry.

### 3. Run with Docker (Recommended)

Start the entire infrastructure (Postgres, Redis, MinIO, Backend, etc.):

```bash
make compose-build
```

To stop:

```bash
make compose-down
```

### 4. Run Locally

Ensure Postgres and Redis are running (or start just the services via docker):

```bash
make compose-services
```

Then run the application:

```bash
go run cmd/backend/main.go
```

## 🗄 Database Migrations

This project uses `golang-migrate` for database version control.

-   **Create a Migration:**
    ```bash
    make migrate-create name=create_users_table
    ```
-   **Run Migrations (Up):**
    ```bash
    make migrate-up
    ```
-   **Rollback (Down):**
    ```bash
    make migrate-down
    ```

## 📚 API Documentation

Swagger documentation is auto-generated.

1.  **Generate Docs:**
    ```bash
    make swag
    ```
2.  **Access UI:**
    Start the server and visit: [http://localhost:9999/swagger/index.html](http://localhost:9999/swagger/index.html)

## ✅ Quality & Testing

### Linting
Run the linter to check for code quality issues:

```bash
make lint
```

### Testing
Run unit and integration tests (uses Testcontainers for DB integration):

```bash
make test
```

### Format
Apply standard Go formatting:

```bash
make format
```

## 📦 Build for Production

To build the binary and package it:

```bash
make zip
```

This will create a `bin/backend` binary and a `.zip` file ready for deployment.

## 🤝 Contributing

1.  Fork the repository.
2.  Create a feature branch (`git checkout -b feature/amazing-feature`).
3.  Commit changes (`git commit -m 'Add amazing feature'`).
4.  Push to branch (`git push origin feature/amazing-feature`).
5.  Open a Pull Request.

---
**License:** MIT
