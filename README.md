# WOL-GO - Wake On LAN Manager (Go Version)

WOL-GO 是一个简单轻量级的局域网远程唤醒/关机图形化管理工具。

## ✨ 特性

- 🌐 **Web GUI** - 现代化的 Web 管理界面
- ⏰ **定时任务** - 支持 Cron 表达式设置自动唤醒/关机
- 🔍 **设备扫描** - ARP 网络扫描自动发现设备
- 📊 **状态监控** - 实时检查设备在线状态（ICMP/ARP/TCP）
- 🔐 **认证系统** - 可选的登录保护
- 🎨 **深色模式** - 支持明/暗主题切换
- 📱 **响应式设计** - 完美支持移动端
- 🐧 **轻量级** - Go 语言实现，资源占用极低

## 🚀 快速开始

### Docker 部署（推荐）

```bash
# 克隆仓库
git clone https://github.com/junex/wol-go.git
cd wol-go

# 启动服务
docker-compose up -d
```

访问 `http://localhost:5000` 即可使用。

### Docker 命令

```bash
docker build -t wol-go:latest .
docker run -d \
  --name=wol-go \
  --network="host" \
  -v ./appdata/db:/app/db \
  -v ./appdata/cron:/etc/cron.d \
  -e PORT=5000 \
  -e TZ=Asia/Shanghai \
  wol-go:latest
```

## ⚙️ 配置

### 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | 5000 | Web 服务端口 |
| `TZ` | UTC | 时区设置 |
| `ENABLE_LOGIN` | false | 启用登录认证 |
| `USERNAME` | admin | 用户名 |
| `PASSWORD` | admin | 密码 |
| `ENABLE_ADD_DEL` | true | 启用添加/删除设备功能 |
| `ENABLE_REFRESH` | true | 启用自动状态刷新 |
| `REFRESH_INTERVAL` | 30 | 状态检查间隔（秒）15/30/60 |
| `PING_TIMEOUT` | 300 | Ping 超时（毫秒） |
| `ARP_TIMEOUT` | 300 | ARP 超时（毫秒） |
| `TCP_TIMEOUT` | 1 | TCP 检查超时（秒） |
| `ARP_INTERFACE` | - | ARP 网络接口 |

## 📊 性能对比

| 指标 | Python 版本 | Go 版本 | 改进 |
|------|------------|---------|------|
| Docker 镜像 | ~150 MB | ~15 MB | **90% ↓** |
| 运行内存 | ~20 MB | ~8 MB | **60% ↓** |
| 启动时间 | ~2-3 秒 | ~0.05 秒 | **98% ↓** |

## 🔧 API 端点

### 认证
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/logout` - 用户登出
- `GET /api/auth/status` - 认证状态

### 设备管理
- `GET /api/computers` - 获取设备列表
- `POST /api/computers` - 添加设备
- `PUT /api/computers/:mac` - 更新设备
- `DELETE /api/computers/:mac` - 删除设备
- `GET /api/computers/:mac/status` - 检查状态
- `POST /api/computers/:mac/wake` - 唤醒设备
- `POST /api/computers/:mac/sleep` - 关机设备

### Cron 任务
- `GET /api/computers/:mac/crons` - 获取定时任务
- `POST /api/computers/:mac/crons/wake` - 添加唤醒任务
- `POST /api/computers/:mac/crons/sleep` - 添加关机任务
- `DELETE /api/computers/:mac/crons/wake` - 删除唤醒任务
- `DELETE /api/computers/:mac/crons/sleep` - 删除关机任务

### 网络
- `GET /api/network/arp-scan` - ARP 网络扫描

## 📦 项目结构

```
wol-go/
├── cmd/server/          # 应用入口
├── internal/
│   ├── api/            # HTTP 处理器、路由
│   ├── config/         # 配置管理
│   ├── model/          # 数据模型
│   ├── service/        # 业务逻辑
│   ├── repository/     # 数据访问层
│   └── pkg/            # 网络功能、验证器
├── docker/             # Docker 配置
├── go.mod
└── docker-compose.yml
```

## 🔄 从 Python 版本迁移

1. 备份数据：
   ```bash
   cp gptwol/appdata/db/computers.txt ./backup_computers.txt
   cp gptwol/appdata/cron/gptwol ./backup_cron
   ```

2. 停止 Python 版本：
   ```bash
   docker-compose down
   ```

3. 启动 Go 版本：
   ```bash
   cd wol-go
   docker-compose up -d
   ```

4. 恢复数据：
   ```bash
   cp ./backup_computers.txt appdata/db/computers.txt
   cp ./backup_cron appdata/cron/wol
   ```

## 🛠️ 开发

```bash
# 安装依赖
go mod download

# 运行
go run ./cmd/server

# 构建
go build -o build/wol-go ./cmd/server

# 测试
go test -v ./...
```

## 📚 更多文档

- **[TODO.md](TODO.md)** - 开发路线图和待办事项
- **[MIGRATION.md](MIGRATION.md)** - 从 Python 版本迁移指南
- **[PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)** - 项目技术总结
- **[DEPLOYMENT_CHECKLIST.md](DEPLOYMENT_CHECKLIST.md)** - 部署检查清单

## 📄 许可证

MIT License

## 🙏 致谢

- 原项目：[GPTWOL](https://github.com/Misterbabou/gptwol) by Misterbabou
- 当前项目：[WOL-GO](https://github.com/junex/wol-go) by junex
- 使用了 [Gorilla Mux](https://github.com/gorilla/mux) 路由器
- 前端使用了 [Bootstrap](https://getbootstrap.com/) 和 [Font Awesome](https://fontawesome.com/)
