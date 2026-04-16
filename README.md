# Blog Community

一个采用前后端分离架构的博客社区项目。

## 项目结构

- `backend/`：Go 后端服务，负责用户、OAuth、鉴权、数据访问等能力
- `frontend/`：Next.js 前端应用，负责页面渲染、交互与状态管理
- `docs/`：产品、架构、接口与开发规划文档
- `.github/workflows/`：CI 工作流配置
- `Makefile`：统一的构建、测试、开发入口

## 技术栈

### Backend

- Go
- Gin
- Ent
- Wire
- PostgreSQL
- Redis
- JWT

### Frontend

- Next.js 16
- React 19
- TypeScript
- Tailwind CSS 4
- TanStack Query
- Zustand
- React Hook Form

## 快速开始

### 1. 安装依赖

前端依赖使用 `pnpm` 管理，请先确保本机已安装：

```bash
pnpm --version
```

后端依赖由 Go Modules 管理。

### 2. 配置环境变量

可参考以下示例文件：

- `backend/.env.example`
- `frontend/.env.example`

按需复制并填写本地配置。

### 3. 启动开发环境

启动后端：

```bash
make dev-backend
```

启动前端：

```bash
make dev-frontend
```

默认情况下：

- 前端地址：`http://localhost:3000`
- 后端地址：`http://localhost:8080`

## 常用命令

```bash
make build          # 构建前后端
make test           # 运行前后端测试
make lint           # 执行检查
make fmt            # 格式化代码
make dev-backend    # 启动后端开发服务
make dev-frontend   # 启动前端开发服务
```

## 文档索引

- `docs/01-PRD.md`
- `docs/02-Architecture.md`
- `docs/03-Tech-Stack.md`
- `docs/04-Frontend-Design.md`
- `docs/05-Backend-Design.md`
- `docs/06-OAuth-Auth.md`
- `docs/07-Error-Logging-Observability.md`
- `docs/08-API-Spec.md`
- `docs/09-Development-Plan.md`
- `docs/10-DevOps-Guide.md`

## 说明

- 根仓库已统一管理前后端代码
- 已配置 GitHub Actions，可结合 `.github/workflows/ci.yml` 查看 CI 流程
- 前端使用较新的 Next.js 版本，开发前建议留意 `frontend/AGENTS.md`
