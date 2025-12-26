# Web 资源同步说明

## 概述

本文档说明了 `web/` 目录和 `internal/api/web/` 目录之间的同步关系和维护指南。

## 目录结构

```
web/                           # 开发时的源码目录
├── static/
│   ├── css/
│   │   └── styles.css         # 样式文件
│   ├── js/
│   │   ├── api.js            # API 调用封装
│   │   ├── app.js            # 主应用逻辑
│   │   ├── config.js         # 配置文件
│   │   ├── toast.js          # Toast 通知系统
│   │   └── websocket.js      # WebSocket 客户端
│   └── images/               # 图片资源目录（当前为空）
└── templates/
    └── index.html            # 单页应用主页面

internal/api/web/              # 嵌入到二进制文件中的web资源
├── static/                    # 与web/相同结构
├── templates/
└── ...
```

## 同步机制

### 为什么需要两个目录？

1. **web/ 目录**：开发时的源码目录，便于开发和修改
2. **internal/api/web/ 目录**：嵌入到 Go 二进制文件中的资源目录

### Go 文件嵌入机制

项目使用 Go 1.16+ 的 embed 功能将 web 资源嵌入到二进制文件中：

```go
// internal/api/embed.go
package api

//go:embed web/static
var staticFiles embed.FS

//go:embed web/templates
var templateFiles embed.FS
```

### 同步步骤

当需要更新 web 资源时，请按以下步骤操作：

1. **在 web/ 目录下进行开发**
2. **测试功能正常**
3. **将 web/ 目录内容同步到 internal/api/web/ 目录**
   ```bash
   cp -r web/* internal/api/web/
   ```
4. **重新编译项目**
   ```bash
   go build -o wol-go cmd/server/main.go
   ```
5. **测试新版本**

## 重要注意事项

### 文件覆盖顺序

1. **internal/api/web/ 目录是最终使用的目录**
2. 嵌入的文件优先级高于外部文件
3. 即使存在外部的 web/ 目录，程序也使用嵌入的资源

### 常见问题

**问题**：修改了 web/ 目录的文件，但网页没有更新
**原因**：程序使用的是嵌入在二进制文件中的资源
**解决**：必须同步到 internal/api/web/ 目录并重新编译

**问题**：编译时提示 embed 相关错误
**原因**：internal/api/web/ 目录结构不匹配 embed 指令
**解决**：确保文件结构与 embed 指令一致

## 开发工作流

### 日常开发

1. 修改 `web/` 目录下的文件
2. 使用 `go run cmd/server/main.go` 启动开发服务器（使用外部文件）
3. 测试功能
4. 完成开发后，同步到 `internal/api/web/`
5. 编译正式版本

### 发布新版本

1. 确保所有功能在 `web/` 目录下测试通过
2. 执行同步命令：`cp -r web/* internal/api/web/`
3. 编译：`go build -o wol-go cmd/server/main.go`
4. 测试编译后的版本

## 文件列表检查

每次同步后，请确保以下文件都已正确复制：

- [ ] `web/static/css/styles.css`
- [ ] `web/static/js/api.js`
- [ ] `web/static/js/app.js`
- [ ] `web/static/js/config.js`
- [ ] `web/static/js/toast.js`
- [ ] `web/static/js/websocket.js`
- [ ] `web/templates/index.html`
- [ ] `web/static/images/` 目录（如果有内容）

## 版本控制

- `web/` 目录应该提交到版本控制系统
- `internal/api/web/` 目录通常不需要提交，因为它是由构建过程生成的
- 但为了确保构建可重现，可以提交同步后的版本

## 故障排除

### 文件不一致

如果发现功能异常，检查：

1. 两个目录的文件是否一致
2. 是否重新编译了项目
3. 嵌入的文件是否正确读取

### 性能问题

1. 大文件会增加二进制文件大小
2. 考虑对静态资源进行压缩
3. 使用 CDN 托管大文件（如果适用）

## 更新记录

- 2025-12-26：初始文档创建
  - 同步了 web 资源
  - 修复了添加设备功能
  - 添加了缺失的 toast.js 和 websocket.js 文件