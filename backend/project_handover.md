# 📂 Project Handover: Go Hexagonal API

**Date:** 2026-01-31
**Status:** Phase 6 Completed (Integration & Load Testing)

This document serves as a comprehensive guide for the AI Assistant to resume work on this project seamlessly.

---

## 🏗 Architecture Overview
The project follows **Hexagonal Architecture (Ports and Adapters)** to decouple core logic from external dependencies.

- **`internal/core`**:
    - **`domain/entity`**: Pure business entities (e.g., `User`, `Auth`, `Role`, `AuditLog`).
    - **`port/input`**: UseCase interfaces (Primary/Driving ports).
    - **`port/output`**: Repository interfaces (Secondary/Driven ports).
    - **`usecase`**: Application business rules implementation.
- **`internal/adapter`**:
    - **`driver/rest`**: HTTP Handlers (Fiber), Middleware, DTOs.
    - **`driven/persistence/postgres`**: GORM implementations of repositories.
    - **`driven/storage/redis`**: Redis implementation for Caching, Rate Limiting, and Token Blacklist.
- **`internal/di`**: Dependency Injection container (`Container` struct) wiring everything together.

## 🛠 Tech Stack
- **Language**: Go 1.25+
- **Web Framework**: Fiber v2
- **Database**: PostgreSQL (GORM)
- **Cache/KV**: Redis (go-redis)
- **Object Storage**: MinIO
- **Observability**: OpenTelemetry (Fiber Contrib), Zap (LoggerX)
- **Testing**: Testify, k6 (Load Testing)

---

## ✅ Implemented Features

### 1. Authentication & Authorization (RBAC)
- **JWT Authentication**: Access (15m) and Refresh (60m) tokens.
- **RBAC**: granular permissions (e.g., `users:view`, `roles:edit`).
- **Middleware**: `RequirePermission` inspects user roles/claims.
- **Secure Logout**: **Token Blacklisting** implemented via Redis. Token is revoked upon logout call.

### 2. Security Enhancements
- **Rate Limiting**:
    - Global: 100 req/min.
    - Auth Routes: 5 req/min (Strict).
    - Backed by **Redis** (Custom `RedisStorage` adapter).
- **Audit Logging**:
    - Asynchronous `Auditor` service.
    - Logs "Who, What, Where, When" to `sys_audit_logs` (Postgres JSONB).
    - Tracks Request ID, IP, User Agent, Latency, and Status Code.

### 3. Data Integrity
- **UUIDs**: All primary keys converted to UUIDv4.
- **Soft Deletes**: Handled via GORM.

### 4. Testing
- **Integration Tests** (`test/integration`):
    - Automated setup spins up API connecting to Docker.
    - `auth_test.go`: Verifies Login -> Me -> Logout -> Fail flow.
    - **Fixes Applied**: Corrected token persistence bug in `AuthUseCase`.
- **Load Testing** (`k6/rate_limit.js`):
    - Validates Rate Limiter response (429 Too Many Requests).

---

## 📌 Current State & Memories
- **Latest Action**: Completed Phase 6. Fixed a bug where `Login` was not persisting the generated token to the DB.
- **Dependencies**: Redis is now a **HARD requirement** for the app to run (used for RateLimit store and TokenRepo).
- **Configuration**:
    - `en-US.yaml` / `pt-BR.yaml` updated with `unauthorized` and `manyRequests` keys.
    - `apperror` codes standardized (`incorrectCredentials`, `disabledUser`).

---

## 🚀 Roadmap / Next Steps

### 1. Observability (Immediate Next Step)
The dependencies are installed, but the stack is commented out in `services.compose.yml`.
- **Task**: Uncomment ElasticSearch, Kibana, APM, and OTEL Collector.
- **Goal**: Visualize Rate Limiting blocks and API latency in Grafana/Kibana.

### 2. Password Recovery
- Implement `POST /auth/forgot-password` and `POST /auth/reset-password`.
- Requires SMTP Service (e.g., Mailhog for dev).

### 3. 2FA (Multi-Factor Auth)
- Implement TOTP implementation for `ADMIN` users.

---

## 💻 Key Commands
```bash
# Start Infrastructure
make compose-up

# Build & Run Locally
go run cmd/backend/main.go

# Run Integration Tests
go test -v ./test/integration/...

# Run Load Test (requires k6)
k6 run k6/rate_limit.js
```
