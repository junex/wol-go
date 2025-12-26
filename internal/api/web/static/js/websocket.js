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
        if (this.ws && (this.ws.readyState === WebSocket.CONNECTING || this.ws.readyState === WebSocket.OPEN)) {
            console.log('[WebSocket] Already connected or connecting');
            return;
        }

        this.status = 'connecting';
        console.log('[WebSocket] Connecting to', this.url);

        try {
            this.ws = new WebSocket(this.url);

            this.ws.onopen = () => {
                this.status = 'connected';
                console.log('[WebSocket] Connected successfully');
                this.emit('connected', {});
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
                    console.log('[WebSocket] Reconnecting in', this.reconnectDelay, 'ms');
                    setTimeout(() => this.connect(), this.reconnectDelay);
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
                setTimeout(() => this.connect(), this.reconnectDelay);
            }
        }
    }

    /**
     * 断开连接
     */
    disconnect() {
        this.manualClose = true;
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

function initWebSocket() {
    if (window.Config && window.Config.WEBSOCKET_URL) {
        wsClient = new WebSocketClient(window.Config.WEBSOCKET_URL);
        wsClient.connect();
        console.log('[WebSocket] Client initialized');
    } else {
        console.log('[WebSocket] WebSocket URL not configured, using HTTP polling fallback');
    }
}

// 页面加载时初始化 WebSocket
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initWebSocket);
} else {
    initWebSocket();
}
