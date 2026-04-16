# 03 - 技术选型方案

> 本文档锁定前后端核心技术栈、版本号及选型理由。

---

## 1. 前端技术栈（Next.js 生态）

| 层级 | 技术选型 | 版本建议 | 说明 |
|------|----------|----------|------|
| **框架** | Next.js | 15.x | App Router + React 19 + RSC，原生 SSR/SSG/ISR |
| **语言** | TypeScript | 5.5+ | 严格模式 (`strict: true`)，类型安全 |
| **样式引擎** | Tailwind CSS | 4.x | 原子化样式，Vite 插件化配置，无 `tailwind.config.js` |
| **CSS 变量** | CSS Variables | - | 支持 Light/Dark 主题切换 |
| **UI 组件库** | shadcn/ui + Radix UI | latest | 无样式依赖、可深度定制、accessibility 友好 |
| **图标** | Lucide React | latest | 简洁现代图标库 |
| **字体** | Geist (Vercel) | latest | Next.js 官方推荐字体 |
| **客户端状态** | Zustand | 4.5+ | 极简，无样板代码，支持 persist 中间件 |
| **服务端状态** | TanStack Query | 5.x | 数据缓存、自动重试、乐观更新 |
| **路由缓存** | Next.js Cache + revalidatePath | - | ISR / Server Actions 缓存策略 |
| **表单** | React Hook Form | 7.5+ | 高性能非受控表单 |
| **表单校验** | Zod | 3.23+ | TypeScript-first Schema 校验 |
| **富文本编辑器** | Tiptap / Editor.js | 2.x / 2.x | 可扩展块编辑器，支持 Markdown 导出 |
| **代码高亮** | Prism.js / Shiki | - | 技术文章代码块高亮 |
| **图表** | Recharts / Tremor | latest | 创作者数据分析图表 |
| **HTTP 客户端** | Ky | 1.x | 轻量、基于 fetch、TypeScript 友好 |
| **构建工具** | Turbopack | built-in | Next.js 15 默认，极速 HMR |
| **测试框架** | Vitest + React Testing Library | latest | 单元/集成测试 |

---

## 2. 后端技术栈（Go 生态）

| 层级 | 技术选型 | 版本建议 | 说明 |
|------|----------|----------|------|
| **语言** | Go | 1.23+ | 高并发、编译快、云原生首选 |
| **Web 框架** | Gin / Echo | 1.10+ / 4.12+ | 推荐 **Gin**（生态丰富）或 **Echo**（简洁） |
| **微服务框架** | Kratos (可选) | 2.x | bilibili 开源，更企业级的微服务框架 |
| **gRPC** | google.golang.org/grpc | 1.67+ | 内部服务通信 |
| **Protobuf** | protoc / buf | latest | 接口定义与代码生成 |
| **ORM** | Ent | 0.14+ | Meta 开源，强类型 Schema 与代码生成 |
| **PG 驱动** | pgx (database/sql) | 5.x | 高性能 PostgreSQL 驱动 |
| **依赖注入** | Wire | 0.6+ | Google 编译期依赖注入 |
| **缓存** | go-redis / rueidis | 9.x / 1.x | Redis 客户端，支持 Cluster |
| **消息队列** | Kafka | 3.x | segmentio/kafka-go 客户端 |
| **搜索引擎** | Elasticsearch | 8.x | olivere/elastic 或官方客户端 |
| **对象存储** | MinIO | latest | MinIO Go SDK |
| **配置管理** | Viper | 1.19+ | 环境变量/配置文件/远程配置统一读取 |
| **日志** | Zap / Slog | 1.27+ / 内置 | 结构化日志，高性能 |
| **链路追踪** | OpenTelemetry | 1.x | OTel SDK + Jaeger Exporter |
| **指标** | Prometheus | 1.20+ | client_golang 埋点 |
| **JWT** | golang-jwt/jwt | 5.x | JWT 签发与校验 |
| **OAuth2** | golang.org/x/oauth2 | latest | OAuth2 客户端 |
| **权限控制** | Casbin | 2.x | RBAC/ABAC 权限模型 |
| **参数校验** | go-playground/validator | 10.x | Struct tag 绑定校验 |
| **密码哈希** | bcrypt (golang.org/x/crypto) | latest | 密码安全存储 |
| **API 文档** | swaggo/swag | 1.16+ | 自动生成 Swagger/OpenAPI 文档 |
| **测试** | testify + go-sqlmock | latest | 单元测试、Mock 数据库 |

---

## 3. 基础设施选型

| 组件 | 选型 | 版本/说明 |
|------|------|-----------|
| **数据库** | PostgreSQL | 15+，主从复制，读写分离预留 |
| **缓存** | Redis | 7.x，支持 RediSearch/RedisJSON（可选） |
| **消息队列** | Apache Kafka | 3.x，3 节点集群 |
| **搜索引擎** | Elasticsearch | 8.x，3 节点集群 |
| **对象存储** | MinIO | 分布式部署，兼容 S3 协议 |
| **分析数据库** | ClickHouse | 24.x，用户行为分析（预留） |
| **容器** | Docker | 24.x+ |
| **编排** | Kubernetes | 1.29+ |
| **CI/CD** | GitHub Actions + ArgoCD | - |
| **监控** | Prometheus + Grafana | - |
| **日志** | Grafana Loki / ELK | - |
| **链路追踪** | Jaeger + OpenTelemetry | - |

---

## 4. 版本锁定建议

在 `frontend/package.json` 与 `backend/go.mod` 中，建议对核心依赖做**最小版本锁定**，避免自动升级导致的破坏性变更：

- **Go**: 使用 `go 1.23` 并在 CI 中校验
- **Node.js**: 使用 `.nvmrc` 锁定 `20.x` 或 `22.x` LTS
- **PNPM**: 使用 `packageManager` 字段锁定 `9.x`
