# ===================================
# Stage 1: Build
# ===================================
FROM golang:1.23-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache git ca-certificates

# 设置工作目录
WORKDIR /build

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 同步 web 目录到 internal/api/web（用于开发时修改 web，构建时自动同步）
RUN cp -r web/* internal/api/web/

# 构建应用（静态链接，无 CGO）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-s -w" -o wol-go ./cmd/server

# ===================================
# Stage 2: Runtime
# ===================================
FROM alpine:3.20

# 安装运行时依赖
RUN apk add --no-cache \
    fping \
    arp-scan \
    netcat-openbsd \
    curl \
    ca-certificates \
    && rm -rf /var/cache/apk/*

# 从构建阶段复制二进制文件
COPY --from=builder /build/wol-go /usr/local/bin/

# 创建数据目录
RUN mkdir -p /app/db /etc/cron.d

# 复制启动脚本
COPY docker/docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# 暴露端口
EXPOSE 5000

# 设置卷
VOLUME ["/app/db", "/etc/cron.d"]

# 入口点
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
