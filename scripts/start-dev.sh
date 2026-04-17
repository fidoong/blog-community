#!/usr/bin/env bash
# 一键启动完整本地开发环境（基础设施 + 后端 + 前端）
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# 颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     Blog Community 本地开发启动器       ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# ──────────────────────────────────────────
# Step 1: 基础设施 (PostgreSQL + Redis)
# ──────────────────────────────────────────
echo -e "${BLUE}[1/3] 检查基础设施...${NC}"

# 尝试检测已有服务
pg_ready=false
redis_ready=false

if nc -z 127.0.0.1 5432 </dev/null 2>/dev/null; then
    echo -e "  ${GREEN}✓ PostgreSQL 已运行 (localhost:5432)${NC}"
    pg_ready=true
fi

if nc -z 127.0.0.1 6379 </dev/null 2>/dev/null; then
    echo -e "  ${GREEN}✓ Redis 已运行 (localhost:6379)${NC}"
    redis_ready=true
fi

# 如果都没有，尝试 Docker Compose
if [[ "$pg_ready" != "true" || "$redis_ready" != "true" ]]; then
    echo -e "  ${YELLOW}→ 尝试启动 Docker Compose...${NC}"
    cd "$PROJECT_ROOT"
    if command -v docker >/dev/null 2>&1 && docker info >/dev/null 2>&1; then
        docker compose -f docker-compose.infra.yml up -d postgres redis 2>&1 | grep -E "Creating|Started|Error" || true
        sleep 3

        # 再次检测
        if nc -z 127.0.0.1 5432 </dev/null 2>/dev/null; then
            pg_ready=true
            echo -e "  ${GREEN}✓ PostgreSQL 已启动${NC}"
        fi
        if nc -z 127.0.0.1 6379 </dev/null 2>/dev/null; then
            redis_ready=true
            echo -e "  ${GREEN}✓ Redis 已启动${NC}"
        fi
    else
        echo -e "  ${YELLOW}! Docker 不可用${NC}"
    fi
fi

# 最终检查
if [[ "$pg_ready" != "true" ]]; then
    echo -e "  ${RED}✗ PostgreSQL 未就绪${NC}"
    echo "    请手动启动 PostgreSQL，或创建 blog/blog123 用户和 blog 数据库"
    exit 1
fi

if [[ "$redis_ready" != "true" ]]; then
    echo -e "  ${YELLOW}! Redis 未就绪，部分功能可能不可用${NC}"
fi

echo ""

# ──────────────────────────────────────────
# Step 2: 启动后端
# ──────────────────────────────────────────
echo -e "${BLUE}[2/3] 启动后端服务...${NC}"
"$SCRIPT_DIR/start-backend.sh"

echo ""

# ──────────────────────────────────────────
# Step 3: 启动前端
# ──────────────────────────────────────────
echo -e "${BLUE}[3/3] 启动前端服务...${NC}"
"$SCRIPT_DIR/start-frontend.sh"

echo ""
echo -e "${GREEN}══════════════════════════════════════════${NC}"
echo -e "${GREEN}  🎉 所有服务已启动！${NC}"
echo -e "${GREEN}══════════════════════════════════════════${NC}"
echo ""
echo "  前端: http://localhost:3000"
echo "  后端: http://localhost:8081"
echo "  API:  http://localhost:8081/api/v1"
echo "  健康: http://localhost:8081/healthz"
echo ""
echo "  后端日志: tail -f /tmp/blog-backend.log"
echo "  前端日志: tail -f /tmp/blog-frontend.log"
echo ""
echo "  一键停止: $SCRIPT_DIR/stop-dev.sh"
echo ""
