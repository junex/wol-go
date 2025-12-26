// API 配置
const API_BASE_URL = window.location.protocol + '//' + window.location.host + '/api';

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

    // 状态刷新间隔（毫秒）
    refreshInterval: 30000,

    // 是否启用自动刷新
    enableRefresh: true
};
