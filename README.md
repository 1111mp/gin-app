# gin-app

A production-ready starter template for building modular REST APIs in Go with **Gin** and **Uber Fx**.

It gives you a clean structure and common backend tooling so you can focus on business logic.

## ✨ Features

- Gin HTTP framework
- Uber Fx dependency injection
- Ent ORM + PostgreSQL
- Redis integration
- JWT authentication
- Swagger (OpenAPI) docs
- Docker and Docker Compose setup
- Migration, test, lint, and mock workflows

## 🧩 Architecture at a glance

This project follows a modular, interface-first structure:

```txt
Controller -> Service -> Repository/Provider
```

- `internal/modules` contains domain modules (auth, user, post, etc.)
- Modules expose interfaces to keep boundaries clean
- Fx wires dependencies across modules

## 📁 Project Structure

```txt
cmd/
  app/                 # Application entrypoint
internal/
  config/              # Configuration
  app/                 # Composition root; all modules/infrastructure are assembled here with Fx.
  modules/             # Business modules (user, post, auth, ...)
  router/              # Gin engine setup and API grouping (public vs private)
  middleware/          # HTTP middlewares (auth, logging, cors, ...)
  dto/                 # Data Transfer Objects (request/response models)
  observability/       # Tracing, metrics, sentry modules
  tasks/               # Async/background work (Asynq, cron).
  rpc/                 # RabbitMQ RPC + gRPC wiring.
pkg/
  logger/              # Logger
  postgres/            # PostgreSQL client
  redis/               # Redis client
  jwt/                 # JWT utilities
  oauth2/              # OAuth2 (Github & Google)
  rabbitmq/            # RabbitMQ RPC server/client
  grpc/                # gRPC server/client
ent/
  schema/              # Ent schemas
  migrate/             # Migration files
```

## 🛠 Prerequisites

- Go `>= 1.22`
- Docker + Docker Compose
- Make

## 🚀 Quick start

1. Create local environment file:

```bash
cp .env.example .env
```

2. Start required services:

```bash
make compose-up
```

3. Run the API:

```bash
make run
```

## 🐳 Docker commands

```bash
make compose-up                  # Start base services
make compose-up-all              # Start full stack
make compose-up-integration-test # Start integration test stack
make compose-down                # Stop and clean services
```

## 📑 API documentation

Generate Swagger docs:

```bash
make swag-v1
```

Then open:

```txt
http://localhost:8080/swagger/index.html
```

## 🧪 Testing

```bash
make test
```

## 🧬 Database workflow

```bash
make schema-create User              # Create new Ent schema
make ent-gen                         # Generate Ent code
make migrate-create add_user_table   # Create migration
make migrate-up                      # Apply migration
make migrate-down                    # Roll back migration
```

## 🧰 Development utilities

```bash
make format            # Go formatting
make linter-golangci   # Lint checks
make deps-audit        # Dependency vulnerability checks
make mock              # Generate mocks
make pre-commit        # Run local pre-commit checks
```

## 📜 License

MIT
