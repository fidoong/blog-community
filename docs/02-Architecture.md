# 02 - 系统架构设计

> 本文档描述博客社区的整体系统架构、服务划分、数据流转及部署拓扑。

---

## 1. 架构设计原则

1. **微服务拆分**：按业务领域拆分服务，独立部署、独立扩展。
2. **前后端分离**：Next.js 负责 SSR/SEO 与 BFF，Go 负责核心微服务。
3. **事件驱动**：Kafka 解耦写操作与下游消费（搜索索引、通知、统计）。
4. **数据分层**：PG 持久化、Redis 缓存、ES 搜索、ClickHouse 分析。
5. **云原生**：容器化 + K8s + DevOps 自动化。

---

## 2. 系统拓扑图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              客户端层 (Client Layer)                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                        │
│  │   Web (Next.js) │  │   Mobile H5   │  │  Admin (Next.js) │                │
│  └──────────────┘  └──────────────┘  └──────────────┘                        │
└─────────────────────────────────────────────────────────────────────────────┘
                                       │
                                       ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            接入层 (Gateway Layer)                             │
│  Nginx (SSL/静态资源) → Kong/AWS ALB (负载均衡) → BFF Gateway (Next.js API / Go Gateway) │
│  职责：限流、鉴权、路由转发、WAF、SSL 终止                                      │
└─────────────────────────────────────────────────────────────────────────────┘
                                       │
                                       ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           服务层 (Microservices Layer)                        │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐            │
│  │ User Svc │ │ Post Svc │ │ Feed Svc │ │ Search   │ │ Notify   │            │
│  │ (Go)     │ │ (Go)     │ │ (Go)     │ │ Svc (Go) │ │ Svc (Go) │            │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘            │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐                                       │
│  │ Comment  │ │ Media    │ │ Admin    │                                       │
│  │ Svc (Go) │ │ Svc (Go) │ │ Svc (Go) │                                       │
│  └──────────┘ └──────────┘ └──────────┘                                       │
│                                                                              │
│  通信方式：内部 gRPC + Protocol Buffers；对外 REST/GraphQL（BFF 聚合）           │
│  服务发现：Consul / Kubernetes Service；配置中心：Apollo / Nacos               │
└─────────────────────────────────────────────────────────────────────────────┘
                                       │
                                       ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           数据层 (Data Layer)                                 │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐              │
│  │ PostgreSQL │  │   Redis    │  │Elasticsearch│  │   Kafka    │              │
│  │ (主业务库)  │  │ (缓存/会话) │  │  (全文搜索)  │  │ (消息队列)  │              │
│  └────────────┘  └────────────┘  └────────────┘  └────────────┘              │
│  ┌────────────┐  ┌────────────┐                                               │
│  │   MinIO    │  │ ClickHouse │                                               │
│  │ (对象存储)  │  │ (行为分析)  │                                               │
│  └────────────┘  └────────────┘                                               │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. 服务职责说明

| 服务 | 职责 | 通信协议 |
|------|------|----------|
| **User Service** | 用户注册/登录/OAuth/关注/粉丝/权限 | gRPC + HTTP |
| **Post Service** | 文章 CRUD、草稿、版本、标签、审核状态机 | gRPC + HTTP |
| **Feed Service** | 推荐流生成、关注流聚合、热榜计算 | gRPC + HTTP |
| **Comment Service** | 评论嵌套回复、评论点赞、评论审核 | gRPC + HTTP |
| **Search Service** | ES 索引管理、搜索接口、热搜词 | gRPC + HTTP |
| **Notify Service** | 站内信、WebSocket 推送、邮件通知 | gRPC + WebSocket |
| **Media Service** | 图片/视频上传、转码、CDN 分发 | HTTP |
| **Admin Service** | 运营后台 API、审核、举报处理 | HTTP |
| **BFF Gateway** | 聚合微服务接口，适配前端需求 | HTTP/GraphQL |

---

## 4. 数据流转示例：发布文章

```
用户点击发布
    │
    ▼
Next.js BFF / API Route
    │
    ▼
Post Service (Go)
    ├── 写入 PostgreSQL (Ent)
    ├── 发送 Kafka Event: PostPublished
    └── 返回成功
    │
    ▼
Kafka 消费端并行处理：
    ├── Search Consumer → 写入 Elasticsearch
    ├── Feed Consumer → 更新作者粉丝的时间线 (Redis)
    ├── Notify Consumer → 通知粉丝有新文章
    └── Analytics Consumer → 写入 ClickHouse (预留)
```

---

## 5. 部署架构

### 5.1 开发环境

- **本地开发**：Docker Compose 一键启动（PG + Redis + Kafka + ES + MinIO）
- **热重载**：前端 `next dev` (Turbopack)，后端 `air` 热编译

### 5.2 生产环境

- **容器编排**：Kubernetes（Namespace 隔离 Staging / Prod）
- **负载均衡**：Nginx Ingress Controller / AWS ALB
- **自动伸缩**：HPA 基于 CPU / 自定义指标（QPS）
- **配置管理**：ConfigMap + Secret；或接入 Apollo / Nacos
- **镜像仓库**：阿里云 ACR / 腾讯云 TCR / Docker Hub

### 5.3 CI/CD 流程

```
Git Push / PR Merge
    │
    ▼
GitHub Actions / GitLab CI
    ├── 单元测试 (Go test -race, Vitest)
    ├── 代码扫描 (golangci-lint, ESLint)
    ├── 构建 Docker 镜像
    └── 推送镜像仓库
    │
    ▼
ArgoCD / Helm
    ├── 拉取最新镜像
    ├── 滚动更新 (Rolling Update)
    └── 健康检查通过
```

---

## 6. 基础设施清单

| 组件 | 选型 | 用途 |
|------|------|------|
| **容器** | Docker | 应用容器化 |
| **编排** | Kubernetes | 服务调度与扩展 |
| **网关** | Nginx Ingress / Kong | 流量入口、限流、SSL |
| **监控** | Prometheus + Grafana | 指标采集与可视化 |
| **日志** | Loki / ELK Stack | 日志聚合与检索 |
| **链路追踪** | Jaeger + OpenTelemetry | 分布式链路追踪 |
| **告警** | Alertmanager / PagerDuty | 异常自动告警 |
