#!/bin/sh
set -e

echo "======================================"
echo "       WOLGO - Wake On LAN Manager"
echo "       Go Version - Production Build"
echo "======================================"

# 检查并创建数据目录
mkdir -p /app/db
mkdir -p /etc/cron.d

# 启动 cron 服务
echo "Starting cron service..."
crond -b

# 启动应用
echo "Starting WOLGO application..."
echo "======================================"
exec /usr/local/bin/wol-go
