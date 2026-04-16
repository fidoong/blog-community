# 07 - 错误处理、日志、监控与可观测性

> 企业级系统的可观测性三支柱：日志（Logging）、指标（Metrics）、链路追踪（Tracing）。

---

## 1. 错误处理策略

### 1.1 分层错误处理原则

| 层级 | 处理策略 |
|------|----------|
| **Infrastructure** | 将底层错误（如 `pgx.ErrNoRows`、`redis.ErrClosed`）转换为领域错误或通用错误 |
| **Application** | 包装错误并添加上下文信息，决定事务回滚或补偿 |
| **Delivery** | 将应用错误映射为 HTTP Status + 统一 JSON 响应；不可预见的错误返回 500 并记录详情 |
| **前端** | 按状态码分类处理：401 跳转登录、403 Toast 提示、500 友好提示 |

### 1.2 错误码规范

格式：`E{HTTPStatus}{3位序号}`

| 错误码 | 含义 | HTTP 状态 |
|--------|------|-----------|
| `E400001` | 请求参数错误 | 400 |
| `E400002` | 参数校验失败 | 400 |
| `E401001` | 未授权，请先登录 | 401 |
| `E401002` | Token 已过期 | 401 |
| `E403001` | 权限不足 | 403 |
| `E404001` | 资源不存在 | 404 |
| `E404002` | 用户不存在 | 404 |
| `E409001` | 资源冲突（如邮箱已注册） | 409 |
| `E429001` | 请求过于频繁 | 429 |
| `E500001` | 服务器内部错误 | 500 |
| `E503001` | 服务暂时不可用 | 503 |

---

## 2. 日志方案

### 2.1 日志分级

| 级别 | 使用场景 |
|------|----------|
| **DEBUG** | 开发调试，详细变量输出 |
| **INFO** | 正常业务流程记录（如登录成功、文章发布） |
| **WARN** | 非致命异常（如降级处理、缓存穿透、重试） |
| **ERROR** | 业务错误（如数据库写入失败、外部 API 超时） |
| **FATAL** | 系统级致命错误（如启动依赖缺失），进程退出 |

### 2.2 结构化日志规范（Go - Zap）

```go
// pkg/logger/logger.go
package logger

import (
    "context"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var defaultLogger *zap.Logger

func Init(env string) {
    var cfg zap.Config
    if env == "production" {
        cfg = zap.NewProductionConfig()
        cfg.EncoderConfig.TimeKey = "timestamp"
        cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    } else {
        cfg = zap.NewDevelopmentConfig()
    }
    l, _ := cfg.Build()
    defaultLogger = l
}

func L() *zap.Logger { return defaultLogger }

func WithContext(ctx context.Context) *zap.Logger {
    traceID := tracer.TraceIDFromContext(ctx)
    spanID := tracer.SpanIDFromContext(ctx)
    return defaultLogger.With(
        zap.String("trace_id", traceID),
        zap.String("span_id", spanID),
    )
}
```

### 2.3 日志字段标准

每条日志必须包含以下字段：

```json
{
  "level": "error",
  "timestamp": "2024-04-16T14:30:00Z",
  "trace_id": "abc123",
  "span_id": "def456",
  "service": "user-service",
  "msg": "failed to create user",
  "error": "duplicate key value violates unique constraint",
  "user_id": 10001,
  "duration_ms": 45
}
```

### 2.4 日志收集架构

```
应用 Pod (stdout/stderr JSON 日志)
    │
    ▼
Node 级 DaemonSet (Fluent Bit / Filebeat)
    │
    ▼
Grafana Loki / ELK Stack
    │
    ▼
Grafana Dashboard / Kibana
```

---

## 3. 指标监控（Metrics）

### 3.1 核心指标

| 指标类型 | 指标名 | 说明 |
|----------|--------|------|
| **Counter** | `http_requests_total` | 按状态码、方法、路径统计总请求数 |
| **Histogram** | `http_request_duration_seconds` | 接口延迟分布（P50/P95/P99） |
| **Gauge** | `goroutines_count` | 当前 Goroutine 数量 |
| **Gauge** | `db_connections_active` | 数据库活跃连接数 |
| **Counter** | `business_events_total` | 业务事件（如文章发布数、点赞数） |

### 3.2 Prometheus 埋点示例（Go）

```go
// pkg/metrics/metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    HTTPRequests = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "http_requests_total",
        Help: "Total number of HTTP requests",
    }, []string{"method", "path", "status"})

    HTTPDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "http_request_duration_seconds",
        Help:    "HTTP request duration in seconds",
        Buckets: prometheus.DefBuckets,
    }, []string{"method", "path"})
)
```

### 3.3 Gin 中间件自动埋点

```go
func PrometheusMetrics() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())

        metrics.HTTPRequests.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
        metrics.HTTPDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)
    }
}
```

---

## 4. 链路追踪（Tracing）

### 4.1 OpenTelemetry + Jaeger 架构

```
前端请求
    │
    ▼
Next.js (OTel Web)
    │ inject trace context
    ▼
BFF Gateway (Go OTel)
    │ gRPC metadata / HTTP Header 透传
    ▼
User Service (Go OTel)
    │
    ├── PostgreSQL (pgx 自动插桩)
    ├── Redis (go-redis 自动插桩)
    └── Kafka Producer (自动插桩)
    │
    ▼
Jaeger Collector
    │
    ▼
Grafana Tempo / Jaeger UI
```

### 4.2 Go OTel 初始化

```go
// pkg/tracer/tracer.go
package tracer

import (
    "context"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/resource"
    tracesdk "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func InitTracer(serviceName, jaegerURL string) (*tracesdk.TracerProvider, error) {
    exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerURL)))
    if err != nil {
        return nil, err
    }

    tp := tracesdk.NewTracerProvider(
        tracesdk.WithBatcher(exp),
        tracesdk.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceName(serviceName),
        )),
    )
    otel.SetTracerProvider(tp)
    return tp, nil
}
```

### 4.3 Context 传播

```go
// Handler 中启动 Span
ctx, span := otel.Tracer("user-service").Start(c.Request.Context(), "RegisterUser")
defer span.End()

// UseCase / Repository 中继续透传 ctx
user, err := uc.userRepo.Create(ctx, u)
```

---

## 5. 告警规则

| 告警名称 | 触发条件 | 级别 |
|----------|----------|------|
| 接口 P99 延迟过高 | `http_request_duration_seconds{quantile="0.99"} > 1s` | Warning |
| 错误率突增 | `rate(http_requests_total{status=~"5.."}[5m]) > 0.01` | Critical |
| 服务宕机 | `up == 0` | Critical |
| DB 连接池耗尽 | `db_connections_active / db_connections_max > 0.8` | Warning |
| Goroutine 泄露 | `goroutines_count > 50000` | Warning |
| 磁盘使用率 | `disk_usage_percent > 85%` | Warning |

---

## 6. 前端可观测性

### 6.1 错误监控

- **Sentry** 或自建上报：捕获 React Error Boundary、Promise 未捕获异常
- 上报字段：用户 ID、页面 URL、错误堆栈、浏览器版本

### 6.2 性能监控

- 使用 Web Vitals 库采集 LCP、FID、CLS
- 上报至自建埋点服务或 Grafana Faro

### 6.3 日志规范

- 生产环境禁止 `console.log`，统一封装 `logger.debug/info/warn/error`
- 敏感信息（Token、密码）禁止打印
