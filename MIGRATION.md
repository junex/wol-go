# 从 Python 版本迁移到 Go 版本

本指南帮助你从 Python/Flask 版本的 GPTWOL 迁移到 WOL-GO。

## ✅ 兼容性保证

Go 版本完全兼容 Python 版本的数据格式和配置：

- ✅ CSV 数据文件格式相同
- ✅ Cron 文件格式相同
- ✅ 所有环境变量保持一致
- ✅ API 端点兼容

## 📋 迁移步骤

### 1. 备份现有数据

```bash
# 进入 Python 版本的目录
cd wol

# 备份数据和配置
mkdir -p backup
cp appdata/db/computers.txt backup/computers.txt
cp appdata/cron/wol backup/wol

# 备份环境变量（如果使用了自定义配置）
env | grep GPTWOL > backup/env.txt
```

### 2. 停止 Python 版本

```bash
# 停止并删除容器
docker-compose down
docker rm gptwol

# 或者
docker stop gptwol
docker rm gptwol
```

### 3. 启动 Go 版本

```bash
# 克隆或进入 Go 版本目录
cd wol-go

# 创建数据目录
mkdir -p appdata/db appdata/cron

# 启动 Go 版本
docker-compose up -d
```

### 4. 恢复数据

```bash
# 恢复计算机数据
cp backup/computers.txt appdata/db/computers.txt

# 恢复 Cron 任务
cp backup/wol appdata/cron/wol
```

### 5. 验证迁移

1. 访问 Web 界面：`http://localhost:5000`
2. 检查设备列表是否正确显示
3. 测试唤醒/关机功能
4. 检查定时任务是否正常

## 🔄 配置对照表

所有环境变量名称保持不变，无需修改 docker-compose.yml：

| Python 版本 | Go 版本 | 兼容性 |
|------------|---------|--------|
| `PORT` | `PORT` | ✅ |
| `TZ` | `TZ` | ✅ |
| `ENABLE_LOGIN` | `ENABLE_LOGIN` | ✅ |
| `USERNAME` | `USERNAME` | ✅ |
| `PASSWORD` | `PASSWORD` | ✅ |
| `ENABLE_ADD_DEL` | `ENABLE_ADD_DEL` | ✅ |
| `ENABLE_REFRESH` | `ENABLE_REFRESH` | ✅ |
| `REFRESH_INTERVAL` | `REFRESH_INTERVAL` | ✅ |
| `PING_TIMEOUT` | `PING_TIMEOUT` | ✅ |
| `ARP_TIMEOUT` | `ARP_TIMEOUT` | ✅ |
| `TCP_TIMEOUT` | `TCP_TIMEOUT` | ✅ |
| `ARP_INTERFACE` | `ARP_INTERFACE` | ✅ |

## 📊 数据文件对比

### computers.txt 格式（完全相同）

```csv
name,mac_address,ip_address,test_type
MyPC,00:11:22:33:44:55,192.168.1.100,icmp
Server,AA:BB:CC:DD:EE:FF,192.168.1.101,8080
NAS,11:22:33:44:55:66,192.168.1.102,arp
```

### Cron 文件格式（完全相同）

```cron
0 8 * * * root /usr/local/bin/wakeonlan 00:11:22:33:44:55
30 22 * * * root /usr/local/bin/wakeonlan 55:44:33:22:11:00
```

## 🆕 新功能

Go 版本新增了一些功能：

1. **更快的性能** - 启动时间从 2-3 秒降低到 0.05 秒
2. **更小的占用** - 内存从 20MB 降低到 8MB
3. **更小的镜像** - Docker 镜像从 150MB 降低到 15MB
4. **单文件部署** - 所有前端资源嵌入二进制，无需挂载卷

## 🐛 故障排查

### 问题：设备列表为空

**原因**：数据文件未正确迁移

**解决**：
```bash
# 检查数据文件
cat appdata/db/computers.txt

# 确保文件格式正确（CSV，逗号分隔）
# 确保文件在正确的位置
```

### 问题：定时任务不工作

**原因**：Cron 文件未正确迁移或格式问题

**解决**：
```bash
# 检查 Cron 文件
cat appdata/cron/wol

# 确保格式为：分 时 日 月 周 用户 命令
# 例如：0 8 * * * root /usr/local/bin/wakeonlan 00:11:22:33:44:55
```

### 问题：无法唤醒设备

**检查清单**：
1. ✅ 设备 BIOS 已启用 WOL
2. ✅ 网络模式为 `host`
3. ✅ MAC 地址格式正确（XX:XX:XX:XX:XX:XX）
4. ✅ 设备在局域网内

### 问题：状态检测不准确

**原因**：防火墙或网络配置

**解决**：
- ICMP 检测需要允许 ping
- TCP 检测需要目标端口开放
- ARP 检测需要在同一网段

## 📞 获取帮助

如果遇到问题：

1. 查看日志：`docker-compose logs -f`
2. 检查配置：`docker-compose config`
3. 提交 Issue：[GitHub Issues](https://github.com/junex/wol-go/issues)

## 🎉 迁移完成后

享受更快速、更轻量的 WOL-GO 体验！

- ⚡ 更快的启动速度
- 💾 更小的资源占用
- 🔒 相同的功能和数据兼容性
