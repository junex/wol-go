// 主应用
class WOLGO {
    constructor() {
        this.computers = [];
        this.filteredComputers = [];
        this.searchTerm = '';
        this.sortBy = 'name';
        this.statusUpdateInterval = null;
    }

    async init() {
        try {
            // 加载计算机列表
            await this.loadComputers();

            // 渲染计算机列表
            this.renderComputers();

            // 启动状态更新
            this.startStatusUpdates();

            // 绑定事件
            this.bindEvents();
        } catch (error) {
            this.showError('加载失败: ' + error.message);
        }
    }

    async loadComputers() {
        try {
            const response = await api.getComputers();
            if (response.success) {
                this.computers = response.data;
                this.filterAndSortComputers();
            }
        } catch (error) {
            console.error('加载计算机失败:', error);
            throw error;
        }
    }

    filterAndSortComputers() {
        let filtered = this.computers;

        // 搜索过滤
        if (this.searchTerm) {
            const term = this.searchTerm.toLowerCase();
            filtered = filtered.filter(c =>
                c.name.toLowerCase().includes(term) ||
                c.mac_address.toLowerCase().includes(term) ||
                c.ip_address.includes(term)
            );
        }

        // 排序
        filtered.sort((a, b) => {
            if (this.sortBy === 'name') {
                return a.name.localeCompare(b.name);
            } else if (this.sortBy === 'ip') {
                return a.ip_address.localeCompare(b.ip_address);
            } else if (this.sortBy === 'mac') {
                return a.mac_address.localeCompare(b.mac_address);
            }
            return 0;
        });

        this.filteredComputers = filtered;
    }

    renderComputers() {
        const container = document.getElementById('computers-container');
        if (!container) return;

        if (this.filteredComputers.length === 0) {
            container.innerHTML = `
                <div class="alert alert-info" role="alert">
                    <i class="fas fa-info-circle"></i> 暂无设备，点击"添加设备"开始使用
                </div>
            `;
            return;
        }

        const cards = this.filteredComputers.map(computer => this.createComputerCard(computer)).join('');
        container.innerHTML = `<div class="row row-cols-1 row-cols-md-2 row-cols-lg-3 g-4">${cards}</div>`;

        // 更新状态指示器
        this.updateAllStatuses();
    }

    createComputerCard(computer) {
        return `
            <div class="col" data-mac="${computer.mac_address}">
                <div class="card h-100 computer-card">
                    <div class="card-body">
                        <div class="d-flex justify-content-between align-items-start mb-2">
                            <h5 class="card-title mb-0">${this.escapeHtml(computer.name)}</h5>
                            <span class="status-indicator" data-ip="${computer.ip_address}" data-test-type="${computer.test_type}"></span>
                        </div>
                        <p class="card-text">
                            <small class="text-muted">
                                <i class="fas fa-network-wired"></i> ${this.escapeHtml(computer.ip_address)}<br>
                                <i class="fas fa-microchip"></i> ${this.escapeHtml(computer.mac_address)}<br>
                                <i class="fas fa-check-circle"></i> ${this.escapeHtml(computer.test_type)}
                            </small>
                        </p>
                    </div>
                    <div class="card-footer bg-transparent border-top-0">
                        <div class="btn-group w-100" role="group">
                            <button class="btn btn-outline-primary btn-sm status-power" data-mac="${computer.mac_address}" title="唤醒/关机">
                                <i class="fas fa-power-off"></i>
                            </button>
                            <button class="btn btn-outline-secondary btn-sm btn-edit" data-mac="${computer.mac_address}" title="编辑">
                                <i class="fas fa-edit"></i>
                            </button>
                            <button class="btn btn-outline-danger btn-sm btn-delete" data-mac="${computer.mac_address}" title="删除">
                                <i class="fas fa-trash"></i>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        `;
    }

    async updateAllStatuses() {
        const indicators = document.querySelectorAll('.status-indicator');
        for (const indicator of indicators) {
            const ip = indicator.getAttribute('data-ip');
            const testType = indicator.getAttribute('data-test-type');
            this.updateStatus(ip, testType, indicator);
        }
    }

    async updateStatus(ip, testType, element) {
        if (!element) return;

        element.className = 'status-indicator';

        try {
            const status = await api.getComputerStatus(
                element.closest('[data-mac]').getAttribute('data-mac')
            );

            if (status === 'awake') {
                element.classList.add('awake');
                element.title = '在线';
            } else {
                element.classList.add('asleep');
                element.title = '离线';
            }

            // 更新电源按钮
            const card = element.closest('.card');
            const powerBtn = card.querySelector('.status-power');
            if (powerBtn) {
                if (status === 'awake') {
                    powerBtn.classList.remove('btn-success');
                    powerBtn.classList.add('btn-danger');
                    powerBtn.title = '关机';
                } else {
                    powerBtn.classList.remove('btn-danger');
                    powerBtn.classList.add('btn-success');
                    powerBtn.title = '唤醒';
                }
            }
        } catch (error) {
            element.classList.add('asleep');
            element.title = '未知';
        }
    }

    startStatusUpdates() {
        if (this.statusUpdateInterval) {
            clearInterval(this.statusUpdateInterval);
        }

        if (CONFIG.enableRefresh) {
            this.statusUpdateInterval = setInterval(() => {
                this.updateAllStatuses();
            }, CONFIG.refreshInterval);
        }
    }

    bindEvents() {
        // 搜索
        const searchInput = document.getElementById('search-input');
        if (searchInput) {
            searchInput.addEventListener('input', (e) => {
                this.searchTerm = e.target.value;
                this.filterAndSortComputers();
                this.renderComputers();
            });
        }

        // 排序
        const sortSelect = document.getElementById('sort-select');
        if (sortSelect) {
            sortSelect.addEventListener('change', (e) => {
                this.sortBy = e.target.value;
                this.filterAndSortComputers();
                this.renderComputers();
            });
        }

        // 添加设备按钮
        const addBtn = document.getElementById('btn-add');
        if (addBtn) {
            addBtn.addEventListener('click', () => this.showAddModal());
        }

        // 设备卡片事件委托
        const container = document.getElementById('computers-container');
        if (container) {
            container.addEventListener('click', (e) => {
                const powerBtn = e.target.closest('.status-power');
                const editBtn = e.target.closest('.btn-edit');
                const deleteBtn = e.target.closest('.btn-delete');

                if (powerBtn) {
                    this.togglePower(powerBtn.getAttribute('data-mac'));
                } else if (editBtn) {
                    this.showEditModal(editBtn.getAttribute('data-mac'));
                } else if (deleteBtn) {
                    this.deleteComputer(deleteBtn.getAttribute('data-mac'));
                }
            });
        }
    }

    async togglePower(mac) {
        try {
            const computer = this.computers.find(c => c.mac_address === mac);
            if (!computer) return;

            const statusIndicator = document.querySelector(`[data-mac="${mac}"] .status-indicator`);
            const isAwake = statusIndicator && statusIndicator.classList.contains('awake');

            if (isAwake) {
                await api.sleepComputer(mac);
                this.showSuccess(`关机指令已发送到 ${computer.name}`);
            } else {
                await api.wakeComputer(mac);
                this.showSuccess(`WOL 魔术包已发送到 ${computer.name}`);
            }

            // 更新状态
            setTimeout(() => this.updateAllStatuses(), 2000);
        } catch (error) {
            this.showError('操作失败: ' + error.message);
        }
    }

    showAddModal() {
        // TODO: 实现添加设备模态框
        alert('添加设备功能开发中...');
    }

    showEditModal(mac) {
        // TODO: 实现编辑设备模态框
        alert('编辑设备功能开发中... MAC: ' + mac);
    }

    async deleteComputer(mac) {
        if (!confirm('确定要删除此设备吗？')) return;

        try {
            await api.deleteComputer(mac);
            this.showSuccess('设备已删除');
            await this.loadComputers();
            this.renderComputers();
        } catch (error) {
            this.showError('删除失败: ' + error.message);
        }
    }

    showSuccess(message) {
        // TODO: 实现更好的通知系统
        alert('成功: ' + message);
    }

    showError(message) {
        alert('错误: ' + message);
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// 页面加载完成后初始化应用
document.addEventListener('DOMContentLoaded', () => {
    window.app = new WOLGO();
    window.app.init();
});
