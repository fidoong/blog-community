#!/usr/bin/env bash
# 一键启动前端服务
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
FRONTEND_DIR="$PROJECT_ROOT/frontend"
LOG_FILE="/tmp/blog-frontend.log"
PID_FILE="/tmp/blog-frontend.pid"
PORT=3000

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

echo -e "${GREEN}▶ 启动 Blog 前端服务${NC}"
echo "  项目: $FRONTEND_DIR"
echo "  端口: $PORT"
echo "  日志: $LOG_FILE"
echo ""

# 检查目录
if [[ ! -d "$FRONTEND_DIR" ]]; then
    echo -e "${RED}✗ 前端目录不存在: $FRONTEND_DIR${NC}"
    exit 1
fi

# 检查 pnpm
cd "$FRONTEND_DIR"
if ! command -v pnpm >/dev/null 2>&1; then
    echo -e "${RED}✗ 未安装 pnpm，请先安装: npm install -g pnpm${NC}"
    exit 1
fi

# 检查 node_modules
if [[ ! -d "$FRONTEND_DIR/node_modules" ]]; then
    echo "→ 安装依赖..."
    pnpm install
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

# 启动服务
echo "→ 启动 Next.js 开发服务器..."
cd "$FRONTEND_DIR"
nohup pnpm dev > "$LOG_FILE" 2>&1 &
new_pid=$!
echo "$new_pid" > "$PID_FILE"

# 等待服务就绪
for i in {1..30}; do
    if curl -s -o /dev/null -w "%{http_code}" "http://localhost:$PORT" 2>/dev/null | grep -q "200"; then
        echo ""
        echo -e "${GREEN}✓ 前端服务启动成功${NC}"
        echo "  访问地址: http://localhost:$PORT"
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
echo -e "${RED}✗ 前端服务启动超时，请检查日志:${NC}"
echo "  tail -f $LOG_FILE"
exit 1
