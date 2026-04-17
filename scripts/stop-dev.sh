#!/usr/bin/env bash
# 一键停止所有开发服务
set -e

# 颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}▶ 停止 Blog 开发服务${NC}"
echo ""

stopped=0

# 停止前端
if [[ -f /tmp/blog-frontend.pid ]]; then
    pid=$(cat /tmp/blog-frontend.pid)
    if kill -0 "$pid" 2>/dev/null; then
        kill -9 "$pid" 2>/dev/null || true
        rm -f /tmp/blog-frontend.pid
        echo -e "  ${GREEN}✓ 前端已停止 (PID: $pid)${NC}"
        stopped=$((stopped + 1))
    else
        rm -f /tmp/blog-frontend.pid
    fi
fi

# 尝试通过端口查找并停止前端
for pid in $(lsof -t -i :3000 2>/dev/null || true); do
    if [[ -n "$pid" ]]; then
        kill -9 "$pid" 2>/dev/null || true
        echo -e "  ${GREEN}✓ 前端已停止 (PID: $pid)${NC}"
        stopped=$((stopped + 1))
    fi
done

# 停止后端
if [[ -f /tmp/blog-backend.pid ]]; then
    pid=$(cat /tmp/blog-backend.pid)
    if kill -0 "$pid" 2>/dev/null; then
        kill -9 "$pid" 2>/dev/null || true
        rm -f /tmp/blog-backend.pid
        echo -e "  ${GREEN}✓ 后端已停止 (PID: $pid)${NC}"
        stopped=$((stopped + 1))
    else
        rm -f /tmp/blog-backend.pid
    fi
fi

# 尝试通过端口查找并停止后端
for pid in $(lsof -t -i :8081 2>/dev/null || true); do
    if [[ -n "$pid" ]]; then
        kill -9 "$pid" 2>/dev/null || true
        echo -e "  ${GREEN}✓ 后端已停止 (PID: $pid)${NC}"
        stopped=$((stopped + 1))
    fi
done

if [[ $stopped -eq 0 ]]; then
    echo -e "  ${YELLOW}! 未发现运行中的服务${NC}"
else
    echo ""
    echo -e "${GREEN}✓ 共停止 $stopped 个服务${NC}"
fi
