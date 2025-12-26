.PHONY: build run test clean docker-build docker-run help

# 变量定义
BINARY_NAME=wol-go
BUILD_DIR=build
CMD_DIR=cmd/server
GO_FILES=$(shell find . -name '*.go' -type f)

# Go 编译参数
GO=go
GOFLAGS=-v
LDFLAGS=-s -w

# 默认目标
help:
	@echo "可用命令:"
	@echo "  make build       - 编译二进制文件"
	@echo "  make run         - 运行应用"
	@echo "  make test        - 运行测试"
	@echo "  make clean       - 清理构建文件"
	@echo "  make docker-build - 构建 Docker 镜像"
	@echo "  make docker-run   - 运行 Docker 容器"

# 编译二进制文件
build:
	@echo "编译 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "编译完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 运行应用
run:
	@echo "启动 $(BINARY_NAME)..."
	$(GO) run ./$(CMD_DIR)

# 运行测试
test:
	@echo "运行测试..."
	$(GO) test -v -race -cover ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf $(BUILD_DIR)
	@echo "清理完成"

# 格式化代码
fmt:
	@echo "格式化代码..."
	$(GO) fmt ./...

# 代码检查
lint:
	@echo "代码检查..."
	golangci-lint run ./...

# 构建 Docker 镜像
docker-build:
	@echo "构建 Docker 镜像..."
	docker build -t wol-go:latest -f docker/Dockerfile .
	@echo "Docker 镜像构建完成"

# 运行 Docker 容器
docker-run:
	@echo "运行 Docker 容器..."
	docker-compose up -d

# 停止 Docker 容器
docker-stop:
	@echo "停止 Docker 容器..."
	docker-compose down

# 安装依赖
deps:
	@echo "下载依赖..."
	$(GO) mod download
	$(GO) mod tidy

# 开发模式（自动重载）
dev:
	@echo "开发模式..."
	air
