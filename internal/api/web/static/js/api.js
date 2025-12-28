// API 客户端
class APIClient {
    constructor() {
        this.token = localStorage.getItem('auth_token') || null;
        // 默认请求超时时间（毫秒）
        this.defaultTimeout = 10000;
    }

    // 通用请求方法（带超时）
    async request(url, options = {}, timeout = this.defaultTimeout) {
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers
        };

        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        const config = {
            ...options,
            headers
        };

        try {
            // 创建超时 Promise
            const timeoutPromise = new Promise((_, reject) => {
                setTimeout(() => reject(new Error('请求超时，请检查网络连接')), timeout);
            });

            // 发起请求
            const fetchPromise = fetch(url, config);

            // 竞速：请求 vs 超时
            const response = await Promise.race([fetchPromise, timeoutPromise]);

            // 尝试解析 JSON
            let data;
            const text = await response.text();
            try {
                data = JSON.parse(text);
            } catch {
                // 如果不是 JSON，创建一个错误响应对象
                data = { success: false, message: text || '请求失败' };
            }

            if (!response.ok) {
                throw new Error(data.error?.message || data.message || '请求失败');
            }

            return data;
        } catch (error) {
            if (error.message === '请求超时，请检查网络连接') {
                console.error('API 请求超时:', url);
            } else {
                console.error('API 请求失败:', error);
            }
            throw error;
        }
    }

    // GET 请求
    async get(url) {
        return this.request(url, { method: 'GET' });
    }

    // POST 请求
    async post(url, data) {
        return this.request(url, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }

    // PUT 请求
    async put(url, data) {
        return this.request(url, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    }

    // DELETE 请求
    async delete(url) {
        return this.request(url, { method: 'DELETE' });
    }

    // 认证相关
    async login(username, password) {
        const response = await this.post(CONFIG.endpoints.auth.login, { username, password });
        if (response.success) {
            this.token = response.data.token;
            localStorage.setItem('auth_token', this.token);
        }
        return response;
    }

    async logout() {
        const response = await this.post(CONFIG.endpoints.auth.logout);
        this.token = null;
        localStorage.removeItem('auth_token');
        return response;
    }

    async checkAuth() {
        try {
            return await this.get(CONFIG.endpoints.auth.status);
        } catch (error) {
            return { success: false, authenticated: false };
        }
    }

    // 设备管理
    async getComputers() {
        return this.get(CONFIG.endpoints.computers);
    }

    async addComputer(computer) {
        return this.post(CONFIG.endpoints.computers, computer);
    }

    async updateComputer(mac, computer) {
        return this.put(CONFIG.endpoints.computers + '/' + mac, computer);
    }

    async deleteComputer(mac) {
        return this.delete(CONFIG.endpoints.computers + '/' + mac);
    }

    async getComputerStatus(mac) {
        const response = await fetch(CONFIG.computerEndpoints(mac).status, {
            headers: { 'Authorization': this.token ? `Bearer ${this.token}` : '' }
        });
        return await response.text();
    }

    async wakeComputer(mac) {
        return this.post(CONFIG.computerEndpoints(mac).wake);
    }

    async sleepComputer(mac) {
        return this.post(CONFIG.computerEndpoints(mac).sleep);
    }

    // Cron 任务
    async getComputerCrons(mac) {
        return this.get(CONFIG.computerEndpoints(mac).crons);
    }

    async addWakeCron(mac, schedule) {
        return this.post(CONFIG.computerEndpoints(mac).wakeCron, { schedule });
    }

    async addSleepCron(mac, schedule) {
        return this.post(CONFIG.computerEndpoints(mac).sleepCron, { schedule });
    }

    async deleteWakeCron(mac) {
        return this.delete(CONFIG.computerEndpoints(mac).wakeCron);
    }

    async deleteSleepCron(mac) {
        return this.delete(CONFIG.computerEndpoints(mac).sleepCron);
    }

    // 网络扫描
    async arpScan() {
        return this.get(CONFIG.endpoints.arpScan);
    }

    // 批量操作
    async batchWake(macAddresses) {
        return this.post(CONFIG.endpoints.batch.wake, { mac_addresses: macAddresses });
    }

    async batchSleep(macAddresses) {
        return this.post(CONFIG.endpoints.batch.sleep, { mac_addresses: macAddresses });
    }

    async batchCheckStatus(macAddresses) {
        const macParam = macAddresses.join(',');
        return this.get(`${CONFIG.endpoints.batch.status}?mac=${encodeURIComponent(macParam)}`);
    }

    // 健康检查
    async healthCheck() {
        return this.get(CONFIG.endpoints.health);
    }
}

// 创建全局 API 客户端实例
const api = new APIClient();
