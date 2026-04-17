# 10 - DevOps 与本地开发指南

> 本文档涵盖本地开发环境搭建、Docker Compose 配置、CI/CD 流水线与 K8s 部署规范。

---

## 1. 本地开发环境

### 1.1 前置依赖

| 工具 | 版本 | 用途 |
|------|------|------|
| Node.js | 20.x LTS | 前端运行 |
| pnpm | 9.x | 前端包管理 |
| Go | 1.23+ | 后端开发 |
| Docker | 24.x+ | 容器化服务 |
| Docker Compose | 2.x+ | 本地基础设施 |
| Make | - | 命令聚合 |

### 1.2 一键启动基础设施

项目根目录提供 `docker-compose.infra.yml`：

```yaml
version: "3.8"
services:
  postgres:
    image: postgres:15-alpine
    container_name: blog_pg
    environment:
      POSTGRES_USER: blog
      POSTGRES_PASSWORD: blog123
      POSTGRES_DB: blog
    ports:
      - "5432:5432"
    volumes:
      - pg_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    container_name: blog_redis
    ports:
      - "6379:6379"

  elasticsearch:
    image: elasticsearch:8.11.0
    container_name: blog_es
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
    ports:
      - "9200:9200"

  kafka:
    image: bitnami/kafka:3.6
    container_name: blog_kafka
    environment:
      - KAFKA_CFG_NODE_ID=0
      - KAFKA_CFG_PROCESS_ROLES=controller,broker
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@kafka:9093
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
    ports:
      - "9092:9092"

  minio:
    image: minio/minio:latest
    container_name: blog_minio
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    ports:
      - "9000:9000"
      - "9001:9001"

volumes:
  pg_data:
```

启动命令：

```bash
make infra-up
# 等价于：docker compose -f docker-compose.infra.yml up -d
```

停止命令：

```bash
make infra-down
```

---

## 2. 前端开发

### 2.1 安装与启动

```bash
cd frontend
pnpm install
pnpm dev
```

访问：http://localhost:3000

### 2.2 环境变量

```bash
# frontend/.env.local
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_APP_NAME=BlogCommunity
```

### 2.3 常用命令

```bash
pnpm build        # 生产构建
pnpm lint         # ESLint 检查
pnpm test         # Vitest 测试
pnpm typecheck    # TypeScript 类型检查
```

---

## 3. 后端开发

### 3.1 安装依赖

```bash
cd backend
go mod tidy
```

### 3.2 环境变量

```bash
# backend/.env
APP_ENV=development
APP_NAME=user-service
HTTP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=blog
DB_PASSWORD=blog123
DB_NAME=blog

REDIS_ADDR=localhost:6379

JWT_SECRET=your-super-secret-key
GITHUB_CLIENT_ID=xxx
GITHUB_CLIENT_SECRET=xxx
GOOGLE_CLIENT_ID=xxx
GOOGLE_CLIENT_SECRET=xxx
```

### 3.3 常用命令

```bash
# 运行服务（当前为 Modular Monolith 阶段，统一入口）
go run cmd/api-service/main.go

# 热重载（需安装 air）
air -c .air.toml

# 代码检查
go vet ./...
golangci-lint run

# 测试
go test -race ./...

# Ent 代码生成
go generate ./internal/user/ent

# Wire 生成
cd internal/user && wire
```

---

## 4. Makefile 参考

项目根目录 `Makefile`：

```makefile
.PHONY: infra-up infra-down dev-front dev-back lint test build

infra-up:
	docker compose -f docker-compose.infra.yml up -d

infra-down:
	docker compose -f docker-compose.infra.yml down

dev-front:
	cd frontend && pnpm dev

dev-back:
	cd backend && air

lint:
	cd frontend && pnpm lint
	cd backend && golangci-lint run

test:
	cd frontend && pnpm test
	cd backend && go test -race ./...

build:
	cd frontend && pnpm build
	cd backend && CGO_ENABLED=0 go build -o bin/ ./cmd/...
```

---

## 5. CI/CD 流水线（GitHub Actions）

### 5.1 前端 CI

```yaml
# .github/workflows/frontend-ci.yml
name: Frontend CI

on:
  push:
    paths:
      - 'frontend/**'
  pull_request:
    paths:
      - 'frontend/**'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: pnpm/action-setup@v3
        with:
          version: 9
      - uses: actions/setup-node@v4
        with:
          node-version: 20
          cache: 'pnpm'
          cache-dependency-path: frontend/pnpm-lock.yaml
      - run: cd frontend && pnpm install --frozen-lockfile
      - run: cd frontend && pnpm lint
      - run: cd frontend && pnpm typecheck
      - run: cd frontend && pnpm test
      - run: cd frontend && pnpm build
```

### 5.2 后端 CI

```yaml
# .github/workflows/backend-ci.yml
name: Backend CI

on:
  push:
    paths:
      - 'backend/**'
  pull_request:
    paths:
      - 'backend/**'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - run: cd backend && go mod tidy && go vet ./...
      - run: cd backend && go test -race ./...
      - run: cd backend && CGO_ENABLED=0 go build -o bin/ ./cmd/...
```

### 5.3 镜像构建与推送

```yaml
# .github/workflows/docker-build.yml
name: Docker Build
on:
  push:
    branches: [main]

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - run: |
          docker build -t ghcr.io/${{ github.repository }}/frontend:${{ github.sha }} -f frontend/Dockerfile frontend/
          docker build -t ghcr.io/${{ github.repository }}/gateway:${{ github.sha }} -f backend/Dockerfile --build-arg SERVICE=gateway backend/
          docker push ghcr.io/${{ github.repository }}/frontend:${{ github.sha }}
          docker push ghcr.io/${{ github.repository }}/gateway:${{ github.sha }}
```

---

## 6. K8s 部署规范

### 6.1 目录结构

```
deployments/
├── base/
│   ├── namespace.yaml
│   ├── configmap.yaml
│   ├── secret.yaml
│   ├── frontend-deployment.yaml
│   ├── frontend-service.yaml
│   ├── gateway-deployment.yaml
│   ├── gateway-service.yaml
│   └── ingress.yaml
└── overlays/
    ├── staging/
    │   └── kustomization.yaml
    └── production/
        └── kustomization.yaml
```

### 6.2 Deployment 示例（Gateway）

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gateway
spec:
  replicas: 2
  selector:
    matchLabels:
      app: gateway
  template:
    metadata:
      labels:
        app: gateway
    spec:
      containers:
        - name: gateway
          image: ghcr.io/org/blog-gateway:latest
          ports:
            - containerPort: 8080
          envFrom:
            - configMapRef:
                name: blog-config
            - secretRef:
                name: blog-secrets
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /readyz
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
```

### 6.3 Ingress 示例

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: blog-ingress
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
    - hosts:
        - blog.example.com
      secretName: blog-tls
  rules:
    - host: blog.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: frontend
                port:
                  number: 3000
          - path: /api
            pathType: Prefix
            backend:
              service:
                name: gateway
                port:
                  number: 8080
```

---

## 7. 健康检查规范

每个 Go 服务必须暴露以下接口：

| 端点 | 用途 |
|------|------|
| `GET /healthz` | LivenessProbe，仅检查进程存活 |
| `GET /readyz` | ReadinessProbe，检查 DB/Redis 等依赖可用性 |
| `GET /metrics` | Prometheus 指标采集 |

```go
//  readiness 示例
func ReadinessHandler(db *sql.DB, rdb *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := db.Ping(); err != nil {
            c.AbortWithStatus(503)
            return
        }
        if err := rdb.Ping(c.Request.Context()).Err(); err != nil {
            c.AbortWithStatus(503)
            return
        }
        c.JSON(200, gin.H{"status": "ok"})
    }
}
```

---

## 8. 分支与发布策略

| 分支 | 用途 |
|------|------|
| `main` | 稳定分支，仅接受 PR 合并 |
| `develop` | 日常开发分支 |
| `feature/*` | 新功能分支，从 develop 切出 |
| `hotfix/*` | 线上紧急修复，从 main 切出 |

**发布流程**：
1. 功能开发完成 → PR 到 `develop`
2. 里程碑结束 → PR `develop` 到 `main`
3. `main` 合并后自动打 Tag 并触发镜像构建
4. ArgoCD 检测到新 Tag 后自动同步到 Production
