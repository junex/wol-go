// API 配置
const API_BASE_URL = window.location.protocol + '//' + window.location.host + '/api';

// WebSocket URL (自动转换协议)
const WS_PROTOCOL = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
const WS_URL = WS_PROTOCOL + '//' + window.location.host + '/api/ws';

// 应用配置
const CONFIG = {
    // API 端点
    endpoints: {
        // 认证
        auth: {
            login: `${API_BASE_URL}/auth/login`,
            logout: `${API_BASE_URL}/auth/logout`,
            status: `${API_BASE_URL}/auth/status`
        },
        // 设备
        computers: `${API_BASE_URL}/computers`,
        // 批量操作
        batch: {
            wake: `${API_BASE_URL}/computers/batch/wake`,
            sleep: `${API_BASE_URL}/computers/batch/sleep`,
            status: `${API_BASE_URL}/computers/batch/status`
        },
        // 健康检查
        health: `${API_BASE_URL}/health`,
        // ARP 扫描
        arpScan: `${API_BASE_URL}/network/arp-scan`
    },

    // 构建设备相关端点
    computerEndpoints: (mac) => ({
        base: `${API_BASE_URL}/computers/${mac}`,
        status: `${API_BASE_URL}/computers/${mac}/status`,
        wake: `${API_BASE_URL}/computers/${mac}/wake`,
        sleep: `${API_BASE_URL}/computers/${mac}/sleep`,
        crons: `${API_BASE_URL}/computers/${mac}/crons`,
        wakeCron: `${API_BASE_URL}/computers/${mac}/crons/wake`,
        sleepCron: `${API_BASE_URL}/computers/${mac}/crons/sleep`
    }),

    // WebSocket URL
    WEBSOCKET_URL: WS_URL,

    // 状态刷新间隔（毫秒）
    refreshInterval: 30000,

    // 是否启用自动刷新
    enableRefresh: true
};

// 导出到全局
window.Config = CONFIG;
