# gin-app

English | [简体中文](README.zh-CN.md)

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
config/                # Configuration
internal/
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

- Go `>= 1.26.4` (see `go.mod`)
- Docker + Docker Compose
- Make

## 🚀 Quick start

1. Create local environment file:

```bash
cp .env.example .env
```

Update `.env` for your environment, including database and Redis credentials,
RabbitMQ, and a valid Sentry DSN. Values in `.env.example` include placeholders. In containers, configuration can
come entirely from environment variables; a missing `.env` file is allowed.

2. Install development tools:

```bash
make bin-deps
```

3. Start PostgreSQL, Redis, and RabbitMQ:

```bash
docker compose up -d db redis rabbitmq
```

4. Run the API:

```bash
make run
```

## 🐳 Docker commands

```bash
make compose-up                  # Start PostgreSQL/RabbitMQ and follow logs
make compose-up-all              # Start full stack
make compose-up-integration-test # Start integration test stack
make compose-down                # Stop and clean services
```

## 🛑 Graceful shutdown

Fx handles `SIGINT` (Ctrl+C) and `SIGTERM`. Asynq uses `Start` rather than `Run`,
so Fx is the only component managing these application termination signals.
Existing `fx.OnStop` annotations and `fx.Lifecycle` hooks run serially in reverse
registration order.

The application registers telemetry first, then eagerly initializes PostgreSQL
and Redis through the `infrastructure` module, then registers RPC and task hooks.
The HTTP hook is registered last. This keeps shared dependencies available while
requests and jobs finish, with the following shutdown order:

1. Stop accepting HTTP connections and drain active HTTP requests.
2. Stop cron scheduling and wait for running jobs, then shut down Asynq workers.
3. Drain gRPC and RabbitMQ RPC servers; close initialized client resources after
   their users have stopped.
4. Close Redis and PostgreSQL connections.
5. Flush Sentry, shut down the metrics and tracing providers, then flush and close
   application log files.

The **overall Fx stop timeout is 75 seconds**. Blocking shutdown operations use
local limits so later hooks still have time to clean up:

| Component | Local limit and timeout behavior |
| --- | --- |
| HTTP | 8 seconds; then cancel request contexts and close active connections. |
| Cron | 25 seconds; then cancel jobs and allow up to 1 second for cleanup. |
| Asynq | 10 seconds for the stop hook; the library's worker shutdown timeout remains 8 seconds. |
| gRPC | 5 seconds; then request `Stop` and allow up to 500 ms for cleanup. |
| RabbitMQ RPC server | 5 seconds; then cancel handlers and close the connection with a deadline. |
| RabbitMQ RPC client | 3 seconds for shutdown and connection closure. |
| PostgreSQL | 3 seconds for the stop hook to wait for closure. |
| Sentry, metrics, tracing | 3 seconds each. |
| Log files | 2 seconds for flush and close. |

These are upper limits, not mandatory delays. HTTP, gRPC, and RabbitMQ inherit the
Fx stop context. Cron jobs and RabbitMQ handlers must honor cancellation; Go
cannot forcibly terminate a handler that ignores its context. Timeouts are
returned as errors. Fx continues subsequent hooks after errors while the overall
stop context remains valid; exhausting that context may skip remaining hooks.

Docker Compose sets `stop_grace_period: 80s` for the application, longer than the
75-second application budget. Configure at least the same grace period in other
deployments. Use Ctrl+C for a foreground application, send `SIGTERM` to the
application process, or stop the Compose application with:

```bash
docker compose stop app
```

`SIGKILL` bypasses graceful shutdown. When adding a module, keep its cleanup hook
registered after the hooks for the dependencies it uses, and account for its
local timeout in the overall serial shutdown budget.

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
# Include package-level shutdown tests with the race detector:
go test -race ./pkg/... ./internal/...
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
