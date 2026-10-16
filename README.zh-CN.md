# gin-app

[English](README.md) | 简体中文

gin-app 是基于 **Gin** 和 **Uber Fx** 的 Go 模块化 REST API 项目模板，提供常见后端能力和清晰的项目结构，方便开发业务功能。

## ✨ 功能特性

- Gin HTTP 框架
- Uber Fx 依赖注入
- Ent ORM 与 PostgreSQL
- Redis 集成
- JWT 身份认证
- Swagger（OpenAPI）接口文档
- Docker 与 Docker Compose 部署
- 数据库迁移、测试、代码检查和 Mock 工作流

## 🧩 架构概览

项目按模块组织，通过接口划分职责：

```txt
Controller -> Service -> Repository/Provider
```

- `internal/modules` 存放认证、用户、文章等业务模块。
- 模块通过接口暴露能力，保持职责边界清晰。
- Fx 负责组装各模块及其依赖。

## 📁 项目结构

```txt
cmd/
  app/                 # 应用入口
config/                # 配置
internal/
  app/                 # 应用组装入口，通过 Fx 组合各模块与基础设施
  modules/             # 业务模块：用户、文章、认证等
  router/              # Gin 引擎配置与公共、私有 API 分组
  middleware/          # HTTP 中间件：认证、日志、CORS 等
  dto/                 # 请求、响应等数据传输对象
  observability/       # 链路追踪、指标和 Sentry 模块
  tasks/               # Asynq 异步任务与 cron 定时任务
  rpc/                 # RabbitMQ RPC 与 gRPC 模块组装
pkg/
  logger/              # 日志
  postgres/            # PostgreSQL 客户端
  redis/               # Redis 客户端
  jwt/                 # JWT 工具
  oauth2/              # GitHub 与 Google OAuth2
  rabbitmq/            # RabbitMQ RPC 服务端与客户端
  grpc/                # gRPC 服务端与客户端
ent/
  schema/              # Ent 数据模型
  migrate/             # 数据库迁移文件
```

## 🛠 环境要求

- Go `>= 1.26.4`，以 `go.mod` 为准
- Docker 与 Docker Compose
- Make

## 🚀 快速开始

1. 创建本地环境配置文件：

```bash
cp .env.example .env
```

按本地环境修改 `.env`，包括数据库和 Redis 凭据、RabbitMQ 连接及有效的 Sentry DSN。`.env.example` 中包含需要替换的占位值。容器中可以完全通过环境变量配置，允许不提供 `.env` 文件。

2. 安装开发工具：

```bash
make bin-deps
```

3. 启动 PostgreSQL、Redis 和 RabbitMQ：

```bash
docker compose up -d db redis rabbitmq
```

4. 启动 API：

```bash
make run
```

## 🐳 Docker 命令

```bash
make compose-up                  # 启动 PostgreSQL/RabbitMQ，并持续输出日志
make compose-up-all              # 启动完整服务栈
make compose-up-integration-test # 启动集成测试服务栈并运行测试
make compose-down                # 停止并清理服务
```

## 🛑 优雅停机

Fx 统一处理 `SIGINT`（Ctrl+C）和 `SIGTERM`。Asynq 使用 `Start`，不通过 `Run` 单独监听这些退出信号。各模块保留原生的 `fx.OnStop` 注解和 `fx.Lifecycle` 钩子；停机时，Fx 按注册的逆序串行执行关闭操作。

应用先注册遥测模块，再通过 `infrastructure` 模块提前初始化 PostgreSQL 和 Redis，之后注册 RPC 和任务模块的钩子，最后注册 HTTP 钩子。这样可以在请求和任务执行完毕之前保留所需依赖。关闭顺序如下：

1. HTTP 停止接受新连接，等待正在处理的请求完成。
2. cron 停止调度并等待运行中的任务，随后关闭 Asynq worker。
3. 等待 gRPC 和 RabbitMQ RPC 服务端处理完成；已初始化的客户端资源在其使用者停止后关闭。
4. 关闭 Redis 和 PostgreSQL 连接。
5. 刷新 Sentry，关闭指标与链路追踪 Provider，最后刷新并关闭应用日志文件。

**Fx 的总停机预算为 75 秒**。可能阻塞的关闭操作设有本地超时，为后续清理保留时间：

| 组件 | 本地上限与超时行为 |
| --- | --- |
| HTTP | 8 秒；超时后取消请求 context，并关闭活动连接。 |
| cron | 25 秒；超时后取消任务，再等待最多 1 秒进行清理。 |
| Asynq | 关闭钩子最多等待 10 秒；库自身的 worker 停机超时保持默认 8 秒。 |
| gRPC | 5 秒；超时后发起 `Stop`，再等待最多 500 毫秒进行清理。 |
| RabbitMQ RPC 服务端 | 5 秒；超时后取消 handler，并通过带 deadline 的连接关闭操作退出。 |
| RabbitMQ RPC 客户端 | 3 秒，用于停机与关闭连接。 |
| PostgreSQL | 关闭钩子最多等待 3 秒。 |
| Sentry、指标、链路追踪 | 各 3 秒。 |
| 日志文件 | 最多等待 2 秒刷新并关闭。 |

这些时间是等待上限，正常完成后会立即继续后续清理。HTTP、gRPC 和 RabbitMQ 的关闭操作继承 Fx 的停机 context。cron 任务和 RabbitMQ handler 需要响应取消信号；Go 无法强制终止忽略 context 的 handler。超时会作为错误返回。在总停机 context 仍有效时，Fx 会在钩子报错后继续清理；如果总预算耗尽，剩余钩子可能不再执行。

Docker Compose 为应用设置了 `stop_grace_period: 80s`，长于应用的 75 秒停机预算。其他部署环境也应至少保留相同的退出宽限期。前台运行时可按 Ctrl+C，或向应用进程发送 `SIGTERM`；通过 Compose 运行时使用：

```bash
docker compose stop app
```

`SIGKILL` 会跳过优雅停机。新增模块时，应让模块的清理钩子晚于其依赖的清理钩子注册，并将新增的本地超时纳入总的串行停机预算。

## 📑 API 文档

生成 Swagger 文档：

```bash
make swag-v1
```

随后访问：

```txt
http://localhost:8080/swagger/index.html
```

## 🧪 测试

```bash
make test
# 包含 pkg 下的停机测试，并启用竞态检测：
go test -race ./pkg/... ./internal/...
```

## 🧬 数据库工作流

```bash
make schema-create User              # 创建 Ent 数据模型
make ent-gen                         # 生成 Ent 代码
make migrate-create add_user_table   # 创建迁移文件
make migrate-up                      # 执行迁移
make migrate-down                    # 回滚一次迁移
```

## 🧰 开发工具

```bash
make format            # 格式化 Go 代码
make linter-golangci   # 检查代码
make deps-audit        # 检查依赖漏洞
make mock              # 生成 Mock
make pre-commit        # 执行本地提交前检查
```

## 📜 许可证

MIT
