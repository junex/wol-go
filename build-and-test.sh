#!/bin/bash
# WOL-GO - 构建和测试脚本

set -e

echo "======================================"
echo "  WOL-GO - Build and Test Script"
echo "======================================"
echo ""

# 1. 本地构建测试
echo "1. 本地构建测试..."
echo "   清理旧的构建文件..."
rm -f build/wol-go

echo "   编译 Go 程序..."
go build -o build/wol-go ./cmd/server

if [ -f build/wol-go ]; then
    SIZE=$(du -h build/wol-go | cut -f1)
    echo "   ✓ 本地构建成功!"
    echo "   二进制大小: $SIZE"
else
    echo "   ✗ 本地构建失败!"
    exit 1
fi
echo ""

# 2. 检查 Docker 是否可用
if ! command -v docker &> /dev/null; then
    echo "2. Docker 不可用，跳过 Docker 构建测试"
    echo "   提示: 请安装 Docker Desktop 后再运行 Docker 构建"
    echo ""
    echo "======================================"
    echo "  本地构建完成!"
    echo "======================================"
    echo ""
    echo "下一步:"
    echo "  1. 本地运行: ./build/wol-go"
    echo "  2. Docker 构建: docker build -t wol-go:latest ."
    exit 0
fi

# 3. Docker 镜像构建
echo "2. Docker 镜像构建..."
echo "   构建 Docker 镜像..."
docker build -t wol-go:latest .

if [ $? -eq 0 ]; then
    echo "   ✓ Docker 镜像构建成功!"

    # 获取镜像大小
    SIZE=$(docker images wol-go:latest --format "{{.Size}}")
    echo "   镜像大小: $SIZE"

    # 检查是否达到目标 (< 20MB)
    echo ""
    echo "3. 验证镜像大小..."
    SIZE_BYTES=$(docker images wol-go:latest --format "{{.Size}}")
    echo "   镜像大小: $SIZE_BYTES"

    if [[ "$SIZE_BYTES" == *"MB"* ]]; then
        SIZE_NUM=$(echo $SIZE_BYTES | sed 's/MB//')
        if (( $(echo "$SIZE_NUM < 20" | bc -l) )); then
            echo "   ✓ 镜像大小符合目标 (< 20MB)"
        else
            echo "   ⚠ 镜像大小超过目标，但仍然可用"
        fi
    fi
else
    echo "   ✗ Docker 镜像构建失败!"
    exit 1
fi
echo ""

# 4. 功能测试
echo "4. 功能测试..."
echo "   启动测试容器..."
docker run -d --name wol-go-test \
    --network host \
    -e PORT=5001 \
    -e ENABLE_LOGIN=false \
    -e ENABLE_ADD_DEL=true \
    wol-go:latest

# 等待容器启动
echo "   等待服务启动..."
sleep 3

# 检查容器状态
if docker ps | grep -q wol-go-test; then
    echo "   ✓ 容器运行成功"

    # 测试 API
    echo ""
    echo "   测试 API 端点..."

    # 测试根路径
    echo "   - 测试 GET /"
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:5001/)
    if [ "$HTTP_CODE" = "200" ]; then
        echo "     ✓ GET / 返回 $HTTP_CODE"
    else
        echo "     ✗ GET / 返回 $HTTP_CODE (预期 200)"
    fi

    # 测试 API 端点
    echo "   - 测试 GET /api/computers"
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:5001/api/computers)
    if [ "$HTTP_CODE" = "200" ]; then
        echo "     ✓ GET /api/computers 返回 $HTTP_CODE"
    else
        echo "     ✗ GET /api/computers 返回 $HTTP_CODE (预期 200)"
    fi

    # 测试添加设备
    echo "   - 测试 POST /api/computers"
    RESPONSE=$(curl -s -X POST http://localhost:5001/api/computers \
        -H "Content-Type: application/json" \
        -d '{"name":"TestPC","mac_address":"00:11:22:33:44:55","ip_address":"192.168.1.100","test_type":"icmp"}')
    if echo "$RESPONSE" | grep -q "success"; then
        echo "     ✓ POST /api/computers 成功"
    else
        echo "     ✗ POST /api/computers 失败: $RESPONSE"
    fi

    # 获取设备列表
    echo "   - 测试 GET /api/computers (验证添加)"
    RESPONSE=$(curl -s http://localhost:5001/api/computers)
    if echo "$RESPONSE" | grep -q "TestPC"; then
        echo "     ✓ 设备已成功添加到列表"
    else
        echo "     ✗ 设备未在列表中找到"
    fi

else
    echo "   ✗ 容器启动失败"
    docker logs wol-go-test
fi

# 清理测试容器
echo ""
echo "5. 清理测试环境..."
docker stop wol-go-test > /dev/null 2>&1
docker rm wol-go-test > /dev/null 2>&1
echo "   ✓ 清理完成"
echo ""

echo "======================================"
echo "  所有测试完成!"
echo "======================================"
echo ""
echo "下一步:"
echo "  1. 启动服务: docker-compose up -d"
echo "  2. 查看日志: docker-compose logs -f"
echo "  3. 访问界面: http://localhost:5000"
echo ""
