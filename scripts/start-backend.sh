#!/usr/bin/env bash
# 一键启动后端服务
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKEND_DIR="$PROJECT_ROOT/backend"
ENV_FILE="$BACKEND_DIR/.env"
LOG_FILE="/tmp/blog-backend.log"
PID_FILE="/tmp/blog-backend.pid"
PORT=8081

# 颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# 检查端口占用
check_port() {
    local pid
    pid=$(lsof -t -i :"$PORT" 2>/dev/null || true)
    if [[ -n "$pid" ]]; then
        echo -e "${YELLOW}⚠ 端口 $PORT 已被占用 (PID: $pid)${NC}"
        read -rp "是否杀掉该进程? [y/N] " answer
        if [[ "$answer" =~ ^[Yy]$ ]]; then
            kill -9 "$pid" 2>/dev/null || true
            sleep 1
            echo -e "${GREEN}✓ 已释放端口 $PORT${NC}"
        else
            echo -e "${RED}✗ 取消启动${NC}"
            exit 1
        fi
    fi
}

# 检查环境文件
if [[ ! -f "$ENV_FILE" ]]; then
    echo -e "${RED}✗ 环境文件不存在: $ENV_FILE${NC}"
    echo "请先创建 .env 文件，参考:"
    cat "$BACKEND_DIR/.env.example" 2>/dev/null || true
    exit 1
fi

# Viper 会自动读取 .env 文件，这里只需确保在工作目录下运行
cd "$BACKEND_DIR"

# 检查 .env 中的端口是否匹配
env_port=$(grep "^HTTP_PORT=" "$ENV_FILE" | cut -d= -f2 | tr -d ' ')
if [[ -n "$env_port" && "$env_port" != "$PORT" ]]; then
    PORT="$env_port"
fi

echo -e "${GREEN}▶ 启动 Blog 后端服务${NC}"
echo "  项目: $BACKEND_DIR"
echo "  端口: $PORT"
echo "  日志: $LOG_FILE"
echo ""

# 检查数据库
echo -n "→ 检查 PostgreSQL... "
if nc -z 127.0.0.1 5432 </dev/null 2>/dev/null; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${YELLOW}未运行${NC}"
    echo "  提示: 请先启动 Docker 基础设施"
    echo "    docker compose -f docker-compose.infra.yml up -d postgres redis"
    echo "    或复用本机已有的 PostgreSQL 容器"
fi

echo -n "→ 检查 Redis... "
if nc -z 127.0.0.1 6379 </dev/null 2>/dev/null; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${YELLOW}未运行${NC}"
fi

# 检查端口
check_port

# 清理旧进程
if [[ -f "$PID_FILE" ]]; then
    old_pid=$(cat "$PID_FILE" 2>/dev/null)
    if kill -0 "$old_pid" 2>/dev/null; then
        kill -9 "$old_pid" 2>/dev/null || true
        sleep 1
    fi
    rm -f "$PID_FILE"
fi

# 编译检查
echo "→ 编译检查..."
cd "$BACKEND_DIR"
if ! go build -o /tmp/blog-api-service ./cmd/api-service 2>/dev/null; then
    echo -e "${RED}✗ 编译失败，请检查代码${NC}"
    exit 1
fi
echo -e "${GREEN}✓ 编译通过${NC}"

# 启动服务
echo "→ 启动服务..."
cd "$BACKEND_DIR"

nohup go run ./cmd/api-service > "$LOG_FILE" 2>&1 &
new_pid=$!
echo "$new_pid" > "$PID_FILE"

# 等待健康检查
for i in {1..30}; do
    if curl -s "http://localhost:$PORT/healthz" >/dev/null 2>&1; then
        echo ""
        echo -e "${GREEN}✓ 后端服务启动成功${NC}"
        echo "  健康检查: http://localhost:$PORT/healthz"
        echo "  API 入口: http://localhost:$PORT/api/v1"
        echo "  PID: $new_pid"
        echo ""
        echo "  查看日志: tail -f $LOG_FILE"
        echo "  停止服务: kill $new_pid"
        exit 0
    fi
    sleep 1
    echo -n "."
done

echo ""
echo -e "${RED}✗ 后端服务启动超时，请检查日志:${NC}"
echo "  tail -f $LOG_FILE"
exit 1
