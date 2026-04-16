# 05 - 后端企业级设计方案（Go + Ent + Wire）

> 基于 Go 1.23+、Ent ORM、Google Wire DI 的企业级 Clean Architecture 实践。

---

## 1. 项目目录结构

```
backend/
├── api/                        # 接口定义 (protobuf / OpenAPI)
│   └── proto/
├── cmd/                        # 服务入口
│   ├── user-service/
│   │   └── main.go
│   ├── post-service/
│   │   └── main.go
│   ├── feed-service/
│   │   └── main.go
│   └── gateway/
│       └── main.go
├── internal/                   # 各服务私有代码
│   ├── user/
│   │   ├── delivery/           # HTTP/gRPC Handler
│   │   │   ├── http_handler.go
│   │   │   └── grpc_handler.go
│   │   ├── application/        # UseCase / Service
│   │   │   └── user_usecase.go
│   │   ├── domain/             # 领域实体 + Repository 接口
│   │   │   ├── user.go
│   │   │   ├── repository.go
│   │   │   └── event.go
│   │   ├── infrastructure/     # 基础设施实现
│   │   │   ├── ent_repo.go     # Ent 实现的 Repository
│   │   │   ├── redis_cache.go
│   │   │   └── kafka_publisher.go
│   │   ├── ent/
│   │   │   └── schema/
│   │   │       └── user.go
│   │   ├── wire.go             # Wire Provider 定义
│   │   └── wire_gen.go         # Wire 生成代码
│   ├── post/
│   │   └── ...
│   └── pkg/                    # 服务内公共包
├── pkg/                        # 全局公共包
│   ├── logger/
│   ├── errors/
│   ├── middleware/
│   ├── database/               # Ent + 事务封装
│   ├── cache/
│   ├── pagination/
│   ├── validator/
│   ├── response/
│   ├── tracer/
│   ├── snowflake/
│   └── hashutil/
├── configs/
├── scripts/
├── deployments/
├── tests/
└── go.mod
```

---

## 2. Clean Architecture 分层规则

依赖方向始终向内：

```
Delivery (HTTP / gRPC Handler)
    ↓
Application (UseCase)
    ↓
Domain (Entity + Repository Interface)
    ↑
Infrastructure (Ent / Redis / Kafka / External API)
```

### 各层职责

| 层级 | 职责 | 禁止行为 |
|------|------|----------|
| **Delivery** | 参数绑定、调用 UseCase、返回 HTTP/gRPC 响应 | 直接调用 Repository 或 DB |
| **Application** | 编排业务逻辑、事务管理、领域事件发布 | 直接操作 HTTP Context 或 DB SQL |
| **Domain** | 定义实体、值对象、Repository 接口、领域事件 | 依赖外部框架或基础设施包 |
| **Infrastructure** | 实现 Repository 接口、操作缓存/消息队列/外部 API | 反向依赖 Application 层 |

---

## 3. Google Wire 依赖注入

### 3.1 Provider 组织

每个服务维护独立的 `wire.go`，全局公共依赖在 `pkg/wire` 中注册 ProviderSet。

```go
// internal/user/wire.go
//go:build wireinject
// +build wireinject

package user

import (
    "github.com/google/wire"
    "<module>/internal/user/application"
    "<module>/internal/user/delivery"
    "<module>/internal/user/infrastructure"
    "<module>/pkg/database"
)

func InitializeUserHandler(cfg *Config, client *database.EntClient, rdb *RedisClient) *delivery.UserHandler {
    wire.Build(
        delivery.NewUserHandler,
        application.NewUserUseCase,
        infrastructure.NewEntUserRepo,
        infrastructure.NewRedisCache,
    )
    return &delivery.UserHandler{}
}
```

### 3.2 生成命令

```bash
cd internal/user && wire
```

生成 `wire_gen.go`，包含完整的对象组装逻辑。编译期注入，零运行时反射开销。

---

## 4. Ent + PostgreSQL 工程实践

### 4.1 Schema 定义

```go
// internal/user/ent/schema/user.go
package schema

import (
    "time"
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
)

type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("id"),
        field.String("email").Unique().MaxLen(128),
        field.String("username").Unique().MaxLen(64),
        field.String("password_hash").Optional().Sensitive(),
        field.String("avatar_url").Optional().MaxLen(512),
        field.Enum("oauth_provider").
            Values("none", "github", "google").
            Default("none"),
        field.String("oauth_id").Optional(),
        field.Time("created_at").Immutable().Default(time.Now),
        field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
    }
}

func (User) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("oauth_provider", "oauth_id").Unique(),
    }
}
```

### 4.2 代码生成

```bash
go generate ./internal/user/ent
```

或在 `ent/generate.go` 中放置：

```go
package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./schema
```

### 4.3 Repository 实现（带事务感知）

```go
// internal/user/infrastructure/ent_repo.go
package infrastructure

import (
    "context"
    "<module>/ent"
    "<module>/internal/user/domain"
    "<module>/pkg/database"
)

type entUserRepo struct {
    client *ent.Client
}

func NewEntUserRepo(client *ent.Client) domain.UserRepository {
    return &entUserRepo{client: client}
}

func (r *entUserRepo) Create(ctx context.Context, u *domain.User) error {
    client := database.ExtractTx(ctx, r.client)
    _, err := client.User.Create().
        SetEmail(u.Email).
        SetUsername(u.Username).
        SetPasswordHash(u.PasswordHash).
        Save(ctx)
    return err
}

func (r *entUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    client := database.ExtractTx(ctx, r.client)
    eu, err := client.User.Query().Where(user.Email(email)).Only(ctx)
    if ent.IsNotFound(err) {
        return nil, domain.ErrUserNotFound
    }
    if err != nil {
        return nil, err
    }
    return toDomain(eu), nil
}
```

### 4.4 事务封装

```go
// pkg/database/ent_tx.go
package database

import (
    "context"
    "fmt"
    "<module>/ent"
)

type txKey struct{}

type EntTransactor struct {
    client *ent.Client
}

func NewEntTransactor(client *ent.Client) *EntTransactor {
    return &EntTransactor{client: client}
}

func (t *EntTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
    tx, err := t.client.Tx(ctx)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    txCtx := context.WithValue(ctx, txKey{}, tx)
    if err := fn(txCtx); err != nil {
        if rerr := tx.Rollback(); rerr != nil {
            return fmt.Errorf("rollback failed: %v (original: %w)", rerr, err)
        }
        return err
    }
    return tx.Commit()
}

func ExtractTx(ctx context.Context, client *ent.Client) *ent.Client {
    if tx, ok := ctx.Value(txKey{}).(*ent.Tx); ok {
        return tx.Client()
    }
    return client
}
```

### 4.5 UseCase 层使用事务

```go
// internal/user/application/user_usecase.go
func (uc *userUseCase) Register(ctx context.Context, email, password string) (*domain.User, error) {
    hashed, err := hashutil.BcryptHash(password)
    if err != nil {
        return nil, errors.Wrap(err, errors.ErrInternal)
    }

    u := &domain.User{
        Email:        email,
        PasswordHash: hashed,
    }

    if err := uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
        if err := uc.repo.Create(txCtx, u); err != nil {
            return errors.Wrap(err, errors.ErrInternal)
        }
        return uc.eventBus.Publish(txCtx, domain.UserRegisteredEvent{UserID: u.ID})
    }); err != nil {
        return nil, err
    }

    return u, nil
}
```

---

## 5. 统一错误处理

### 5.1 业务错误封装

```go
// pkg/errors/errors.go
package errors

import "fmt"

type AppError struct {
    Code    string
    Message string
    Details map[string]any
    Cause   error
}

func (e *AppError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
    }
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Cause }

var (
    ErrInvalidInput = New("E400001", "请求参数错误")
    ErrUnauthorized = New("E401001", "未授权，请先登录")
    ErrForbidden    = New("E403001", "权限不足")
    ErrNotFound     = New("E404001", "资源不存在")
    ErrInternal     = New("E500001", "服务器内部错误")
)

func New(code, message string) *AppError {
    return &AppError{Code: code, Message: message}
}

func Wrap(cause error, appErr *AppError) *AppError {
    return &AppError{
        Code:    appErr.Code,
        Message: appErr.Message,
        Cause:   cause,
    }
}
```

### 5.2 Gin 全局错误拦截中间件

```go
// pkg/middleware/error_handler.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "<module>/pkg/errors"
    "<module>/pkg/response"
    "net/http"
    "strings"
)

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            status := http.StatusInternalServerError
            code := "E500001"
            message := "服务器内部错误"

            var appErr *errors.AppError
            if errors.As(err, &appErr) {
                code = appErr.Code
                message = appErr.Message
                status = mapStatus(code)
            }

            response.Fail(c.Writer, status, code, message)
        }
    }
}

func mapStatus(code string) int {
    switch {
    case strings.HasPrefix(code, "E400"):
        return http.StatusBadRequest
    case strings.HasPrefix(code, "E401"):
        return http.StatusUnauthorized
    case strings.HasPrefix(code, "E403"):
        return http.StatusForbidden
    case strings.HasPrefix(code, "E404"):
        return http.StatusNotFound
    default:
        return http.StatusInternalServerError
    }
}
```

### 5.3 统一响应结构

```go
// pkg/response/response.go
package response

import (
    "encoding/json"
    "net/http"
    "time"
)

type Response[T any] struct {
    Code      string `json:"code"`
    Message   string `json:"message"`
    Data      T      `json:"data"`
    Timestamp int64  `json:"timestamp"`
}

func Success[T any](w http.ResponseWriter, data T) {
    JSON(w, http.StatusOK, Response[T]{
        Code:      "OK",
        Message:   "success",
        Data:      data,
        Timestamp: time.Now().Unix(),
    })
}

func Fail(w http.ResponseWriter, status int, code, message string) {
    JSON(w, status, Response[any]{
        Code:      code,
        Message:   message,
        Timestamp: time.Now().Unix(),
    })
}

func JSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(data)
}
```

---

## 6. 中间件栈

```go
// pkg/middleware/
├── recovery.go          # Panic 恢复
├── logger.go            # Access Log
├── cors.go              # 跨域配置
├── jwt.go               # JWT Token 校验
├── oauth.go             # OAuth Session 校验
├── rate_limit.go        # Token Bucket / Redis 滑动窗口限流
├── tracing.go           # OpenTelemetry Trace 注入
├── error_handler.go     # 全局错误拦截
└── rbac.go              # Casbin 权限校验
```

中间件注册顺序（Gin）：

```go
r := gin.New()
r.Use(middleware.Recovery(logger))
r.Use(middleware.Logger(logger))
r.Use(middleware.CORS())
r.Use(middleware.Tracing())
r.Use(middleware.RateLimit(redisClient))
r.Use(middleware.ErrorHandler())
```

---

## 7. 领域事件与 Kafka

```go
// internal/user/domain/event.go
package domain

type UserRegisteredEvent struct {
    UserID uint64
    Email  string
}

// internal/user/infrastructure/kafka_publisher.go
func (p *kafkaPublisher) Publish(ctx context.Context, event domain.Event) error {
    payload, _ := json.Marshal(event)
    return p.producer.WriteMessages(ctx, kafka.Message{
        Topic: event.Topic(),
        Key:   []byte(fmt.Sprintf("%d", event.AggregateID())),
        Value: payload,
    })
}
```
