# gin-app 🚀

A production-ready **Gin + Fx** based project template for rapidly building scalable, modular, RESTful web services in Go.

This template is designed with:

- 🧱 Modular architecture (Fx DI)
- 🧠 Clean architecture layering
- 🔐 JWT authentication
- 📦 Ent ORM
- 🐳 Docker & Docker Compose
- 📑 Swagger API docs
- 🧪 Testing & mocking
- 🔁 Database migration

---

## ✨ Features

- Gin HTTP framework
- Uber Fx dependency injection
- PostgreSQL (Ent ORM)
- Redis
- JWT authentication
- Modular domain structure
- Swagger documentation
- Integration testing stack
- Database migration
- Linting & formatting
- Mock generation
- Dependency vulnerability scanning

---

## 📁 Project Structure

```txt
cmd/
  app/                 # Application entrypoint
internal/
  config/              # Configuration
  modules/             # Business modules (user, post, auth, ...)
  router/              # Router layer
  middleware/          # HTTP middlewares (auth, logging, cors, ...)
  dto/                 # Data Transfer Objects (request/response models)
pkg/
  logger/              # Logger
  postgres/            # PostgreSQL client
  redis/               # Redis client
  jwt/                 # JWT utilities
  oauth2/              # OAuth2 (Github & Google)
ent/
  schema/              # Ent schemas
  migrate/             # Migration files
```

---

## 🧩 Architecture Design

- **Module-based architecture** (inspired by NestJS)
- Fx global container with interface-based exports
- Strong decoupling via interfaces
- Dependency inversion principle (DIP)
- Clean boundaries between modules

```txt
Controller -> Service -> Domain Interface -> Infrastructure
```

---

## 🛠 Prerequisites

### Docker

```bash
brew install docker
```

### Go

```bash
go version >= 1.22
```

---

## 🐘 PostgreSQL (Docker)

```bash
docker run -d --name pg-server \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=db_name \
  -p 5432:5432 \
  -v pgdata:/var/lib/postgresql \
  postgres
```

---

## 🧠 Redis (Docker)

```bash
docker run -d --name redis \
  -p 6379:6379 \
  -v redis-data:/data \
  redis redis-server --requirepass yourpassword --appendonly yes
```

---

## ⚙️ Environment Setup

```bash
cp .env.example .env
```

---

## 🚀 Quick Start

```bash
make run
```

---

## 🐳 Docker Compose

### Start base services

```bash
make compose-up
```

### Start all services

```bash
make compose-up-all
```

### Integration test stack

```bash
make compose-up-integration-test
```

### Stop services

```bash
make compose-down
```

---

## 📑 Swagger

```bash
make swag-v1
```

Access:

```
http://localhost:8080/swagger/index.html
```

---

## 🧪 Testing

```bash
make test
```

---

## 🧬 Database

### Create schema

```bash
make schema-create User
```

### Generate ent code

```bash
make ent-gen
```

### Create migration

```bash
make migrate-create add_user_table
```

### Run migration

```bash
make migrate-up
```

### Rollback migration

```bash
make migrate-down
```

---

## 🧰 Development Tools

### Format code

```bash
make format
```

### Lint

```bash
make linter-golangci
```

### Dependency audit

```bash
make deps-audit
```

### Generate mocks

```bash
make mock
```

---

## 📦 Dependency Management

```bash
make deps
make upgrade-deps
make ls-deps
```

---

## 🔁 Pre-commit Hook

```bash
make pre-commit
```

---

## 🧠 Module Design Pattern

Each module follows:

```txt
module/
  module.go     # fx.Module definition
  service.go    # internal implementation
  repo.go       # repository
  controller.go # HTTP layer
```

External modules only depend on interfaces.

---

## 🔐 Security

- JWT authentication
- Middleware-based authorization
- Environment-based secrets
- Isolated service layers

---

## 🧱 Recommended Extensions

- OpenTelemetry tracing
- Prometheus metrics
- gRPC gateway
- CQRS pattern
- Event-driven architecture
- Message queue integration

---

## 📜 License

MIT License

---

## 🤝 Contributing

PRs are welcome. Please follow:

- Go formatting standards
- Modular structure
- Interface-first design
- Clean commit history

---

## 🧭 Philosophy

> Simple by default, scalable by design.

> Architecture should enable growth, not block it.

---

## ⭐ Why gin-app

- Enterprise-ready
- Clean architecture
- Strong boundaries
- Modular design
- High scalability
- Production-grade tooling

---

Happy coding! 🚀
