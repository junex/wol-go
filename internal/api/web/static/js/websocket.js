/**
 * WebSocket 客户端
 * 处理实时状态更新和事件推送
 */

class WebSocketClient {
    constructor(url) {
        this.url = url;
        this.ws = null;
        this.reconnectDelay = 3000; // 3 秒后重连
        this.manualClose = false;
        this.messageHandlers = new Map();
        this.status = 'disconnected'; // disconnected, connecting, connected
        this.reconnectTimer = null;
        this.lastConnectTime = 0;
        this.visibilityHandler = null;

        // 默认消息处理器
        this.defaultHandlers = {
            'status': this.handleStatusUpdate.bind(this),
            'computer_added': this.handleComputerAdded.bind(this),
            'computer_updated': this.handleComputerUpdated.bind(this),
            'computer_deleted': this.handleComputerDeleted.bind(this),
            'cron_added': this.handleCronAdded.bind(this),
            'cron_deleted': this.handleCronDeleted.bind(this),
        };
    }

    /**
     * 连接 WebSocket
     */
    connect() {
        // 清除之前的重连定时器
        if (this.reconnectTimer) {
            clearTimeout(this.reconnectTimer);
            this.reconnectTimer = null;
        }

        if (this.ws && (this.ws.readyState === WebSocket.CONNECTING || this.ws.readyState === WebSocket.OPEN)) {
            console.log('[WebSocket] Already connected or connecting');
            return;
        }

        this.status = 'connecting';
        this.lastConnectTime = Date.now();
        console.log('[WebSocket] Connecting to', this.url);

        try {
            this.ws = new WebSocket(this.url);

            this.ws.onopen = () => {
                this.status = 'connected';
                console.log('[WebSocket] Connected successfully');
                this.emit('connected', {});

                // 连接成功后重新获取设备列表
                this.emit('reconnect', { timestamp: this.lastConnectTime });
            };

            this.ws.onmessage = (event) => {
                this.handleMessage(event.data);
            };

            this.ws.onclose = (event) => {
                this.status = 'disconnected';
                console.log('[WebSocket] Connection closed:', event.code, event.reason);
                this.emit('disconnected', { code: event.code, reason: event.reason });

                // 自动重连（如果不是手动关闭）
                if (!this.manualClose) {
                    this.scheduleReconnect();
                }
            };

            this.ws.onerror = (error) => {
                console.error('[WebSocket] Error:', error);
                this.emit('error', { error });
            };

        } catch (error) {
            console.error('[WebSocket] Failed to connect:', error);
            this.status = 'disconnected';

            // 自动重连
            if (!this.manualClose) {
                this.scheduleReconnect();
            }
        }
    }

    /**
     * 安排重连
     */
    scheduleReconnect() {
        if (this.reconnectTimer) {
            clearTimeout(this.reconnectTimer);
        }

        console.log('[WebSocket] Reconnecting in', this.reconnectDelay, 'ms');
        this.reconnectTimer = setTimeout(() => {
            this.reconnectTimer = null;
            this.connect();
        }, this.reconnectDelay);
    }

    /**
     * 开始监听页面可见性变化
     */
    startVisibilityHandler() {
        if (this.visibilityHandler) {
            return; // 已经在监听
        }

        this.visibilityHandler = () => this.handleVisibilityChange();
        document.addEventListener('visibilitychange', this.visibilityHandler);
        console.log('[WebSocket] Visibility handler started');
    }

    /**
     * 停止监听页面可见性变化
     */
    stopVisibilityHandler() {
        if (this.visibilityHandler) {
            document.removeEventListener('visibilitychange', this.visibilityHandler);
            this.visibilityHandler = null;
            console.log('[WebSocket] Visibility handler stopped');
        }
    }

    /**
     * 处理页面可见性变化
     */
    handleVisibilityChange() {
        if (document.hidden) {
            console.log('[WebSocket] Page hidden');
            // 页面隐藏时不需要特殊处理，WebSocket 会自动处理
        } else {
            console.log('[WebSocket] Page visible');

            // 页面重新可见时，检查连接状态
            const timeSinceLastConnect = Date.now() - this.lastConnectTime;

            // 如果超过 30 秒没有连接尝试，强制重新连接
            if (timeSinceLastConnect > 30000 || !this.isConnected()) {
                console.log('[WebSocket] Page visible after long time, reconnecting...');
                this.emit('page_visible', { timeSinceLastConnect });

                // 如果连接已断开，尝试重新连接
                if (!this.isConnected()) {
                    this.connect();
                } else {
                    // 连接正常，也触发重新获取数据
                    this.emit('reconnect', { timestamp: Date.now() });
                }
            }
        }
    }

    /**
     * 断开连接
     */
    disconnect() {
        this.manualClose = true;

        // 清除重连定时器
        if (this.reconnectTimer) {
            clearTimeout(this.reconnectTimer);
            this.reconnectTimer = null;
        }

        // 停止可见性监听
        this.stopVisibilityHandler();

        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
        this.status = 'disconnected';
        console.log('[WebSocket] Disconnected manually');
    }

    /**
     * 处理接收到的消息
     */
    handleMessage(data) {
        try {
            const message = JSON.parse(data);
            console.log('[WebSocket] Received message:', message);

            const { type, payload } = message;

            // 调用类型特定的处理器
            if (this.defaultHandlers[type]) {
                this.defaultHandlers[type](payload);
            }

            // 调用自定义处理器
            if (this.messageHandlers.has(type)) {
                this.messageHandlers.get(type)(payload);
            }

            // 触发通用消息事件
            this.emit('message', message);

        } catch (error) {
            console.error('[WebSocket] Failed to parse message:', error);
        }
    }

    /**
     * 发送消息
     */
    send(type, payload) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const message = { type, payload };
            this.ws.send(JSON.stringify(message));
            console.log('[WebSocket] Sent message:', message);
        } else {
            console.warn('[WebSocket] Cannot send message, WebSocket is not connected');
        }
    }

    /**
     * 注册消息处理器
     */
    on(type, handler) {
        this.messageHandlers.set(type, handler);
    }

    /**
     * 移除消息处理器
     */
    off(type) {
        this.messageHandlers.delete(type);
    }

    /**
     * 触发事件（用于内部使用）
     */
    emit(type, data) {
        if (this.messageHandlers.has(type)) {
            this.messageHandlers.get(type)(data);
        }
    }

    /**
     * 获取连接状态
     */
    getStatus() {
        return this.status;
    }

    /**
     * 检查是否已连接
     */
    isConnected() {
        return this.status === 'connected' && this.ws && this.ws.readyState === WebSocket.OPEN;
    }

    // ===== 默认消息处理器 =====

    /**
     * 处理设备状态更新
     */
    handleStatusUpdate(payload) {
        console.log('[WebSocket] Status update:', payload);
        const { mac_address, online } = payload;

        // 更新 UI 中的设备状态
        this.updateComputerStatusUI(mac_address, online);
    }

    /**
     * 处理设备添加事件
     */
    handleComputerAdded(payload) {
        console.log('[WebSocket] Computer added:', payload);
        this.emit('computer_added', payload);
    }

    /**
     * 处理设备更新事件
     */
    handleComputerUpdated(payload) {
        console.log('[WebSocket] Computer updated:', payload);
        this.emit('computer_updated', payload);
    }

    /**
     * 处理设备删除事件
     */
    handleComputerDeleted(payload) {
        console.log('[WebSocket] Computer deleted:', payload);
        this.emit('computer_deleted', payload);
    }

    /**
     * 处理 Cron 任务添加事件
     */
    handleCronAdded(payload) {
        console.log('[WebSocket] Cron added:', payload);
        this.emit('cron_added', payload);
    }

    /**
     * 处理 Cron 任务删除事件
     */
    handleCronDeleted(payload) {
        console.log('[WebSocket] Cron deleted:', payload);
        this.emit('cron_deleted', payload);
    }

    // ===== UI 更新辅助方法 =====

    /**
     * 更新设备状态 UI
     */
    updateComputerStatusUI(macAddress, online) {
        // 更新状态指示器
        const statusIndicator = document.querySelector(`[data-mac="${macAddress}"] .status-indicator`);
        if (statusIndicator) {
            statusIndicator.className = `status-indicator status-${online ? 'online' : 'offline'}`;
        }

        // 更新状态文本
        const statusText = document.querySelector(`[data-mac="${macAddress}"] .status-text`);
        if (statusText) {
            statusText.textContent = online ? '在线' : '离线';
        }

        // 更新操作按钮
        const wakeBtn = document.querySelector(`[data-mac="${macAddress}"] .btn-wake`);
        const sleepBtn = document.querySelector(`[data-mac="${macAddress}"] .btn-sleep`);

        if (wakeBtn && sleepBtn) {
            if (online) {
                wakeBtn.classList.add('d-none');
                sleepBtn.classList.remove('d-none');
            } else {
                wakeBtn.classList.remove('d-none');
                sleepBtn.classList.add('d-none');
            }
        }
    }
}

// 创建全局 WebSocket 客户端实例（如果配置了 WebSocket URL）
let wsClient = null;
let wsInitPromise = null;

// 初始化 WebSocket 并返回 Promise
function initWebSocket() {
    if (wsInitPromise) {
        return wsInitPromise; // 防止重复初始化
    }

    wsInitPromise = new Promise((resolve) => {
        if (wsClient) {
            resolve(wsClient);
            return;
        }

        if (!window.Config?.WEBSOCKET_URL) {
            console.log('[WebSocket] WebSocket URL not configured, using HTTP polling fallback');
            resolve(null);
            return;
        }

        wsClient = new WebSocketClient(window.Config.WEBSOCKET_URL);
        wsClient.connect();
        // 启动页面可见性监听
        wsClient.startVisibilityHandler();
        console.log('[WebSocket] Client initialized');

        // 等待连接完成再 resolve
        const checkConnected = () => {
            if (wsClient.isConnected()) {
                resolve(wsClient);
            } else {
                setTimeout(checkConnected, 100);
            }
        };
        checkConnected();
    });

    return wsInitPromise;
}

// 获取 WebSocket 客户端（异步等待）
function getWebSocketClient() {
    return wsInitPromise || initWebSocket();
}

// 页面加载时初始化 WebSocket
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initWebSocket);
} else {
    initWebSocket();
}
