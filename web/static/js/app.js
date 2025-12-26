// 主应用
class WOLGO {
    constructor() {
        this.computers = [];
        this.filteredComputers = [];
        this.searchTerm = '';
        this.sortBy = 'name';
        this.statusUpdateInterval = null;
        this.selectedComputers = new Set(); // 选中的设备 MAC 地址
        this.computerModal = null; // 模态框实例
    }

    async init() {
        try {
            // 初始化模态框
            this.initComputerModal();

            // 加载计算机列表
            await this.loadComputers();

            // 渲染计算机列表
            this.renderComputers();

            // 启动状态更新
            this.startStatusUpdates();

            // 绑定事件
            this.bindEvents();
        } catch (error) {
            toast.error('加载失败: ' + error.message);
        }
    }

    initComputerModal() {
        const modalEl = document.getElementById('computerModal');
        if (modalEl) {
            this.computerModal = new bootstrap.Modal(modalEl);

            // 监听检测方式变化
            const typeSelect = document.getElementById('computer-type');
            if (typeSelect) {
                typeSelect.addEventListener('change', (e) => {
                    const portGroup = document.getElementById('tcp-port-group');
                    if (e.target.value === 'tcp') {
                        portGroup.style.display = 'block';
                    } else {
                        portGroup.style.display = 'none';
                    }
                });
            }

            // 保存按钮
            const saveBtn = document.getElementById('btn-save-computer');
            if (saveBtn) {
                saveBtn.addEventListener('click', () => this.saveComputer());
            }
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

        // 更新设备数量显示
        const deviceCount = document.getElementById('device-count');
        if (deviceCount) {
            deviceCount.textContent = `${this.filteredComputers.length} 台设备`;
        }

        if (this.filteredComputers.length === 0) {
            container.innerHTML = `
                <div class="alert alert-info" role="alert">
                    <i class="fas fa-info-circle"></i> 暂无设备，点击"添加设备"开始使用
                </div>
            `;
            this.updateBatchActionsBar();
            return;
        }

        // 批量操作工具栏
        const batchToolbar = this.createBatchToolbar();
        const cards = this.filteredComputers.map(computer => this.createComputerCard(computer)).join('');
        container.innerHTML = batchToolbar + `<div class="row row-cols-1 row-cols-md-2 row-cols-lg-3 g-4">${cards}</div>`;

        // 更新状态指示器
        this.updateAllStatuses();
        this.updateBatchActionsBar();
    }

    createComputerCard(computer) {
        const isSelected = this.selectedComputers.has(computer.mac_address);
        return `
            <div class="col" data-mac="${computer.mac_address}">
                <div class="card h-100 computer-card ${isSelected ? 'border-primary' : ''}">
                    <div class="card-body">
                        <div class="d-flex justify-content-between align-items-start mb-2">
                            <div class="d-flex align-items-center">
                                <input type="checkbox" class="form-check-input me-2 batch-checkbox"
                                    data-mac="${computer.mac_address}" ${isSelected ? 'checked' : ''}>
                                <h5 class="card-title mb-0">${this.escapeHtml(computer.name)}</h5>
                            </div>
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

    createBatchToolbar() {
        const selectedCount = this.selectedComputers.size;
        const totalCount = this.filteredComputers.length;
        const allSelected = selectedCount > 0 && selectedCount === totalCount;

        return `
            <div class="batch-toolbar mb-3 p-3 bg-light rounded" id="batch-toolbar">
                <div class="d-flex justify-content-between align-items-center flex-wrap gap-2">
                    <div class="d-flex align-items-center gap-2">
                        <input type="checkbox" class="form-check-input" id="select-all"
                            ${allSelected ? 'checked' : ''} ${totalCount === 0 ? 'disabled' : ''}>
                        <label for="select-all" class="form-check-label">
                            ${selectedCount > 0 ? `已选 ${selectedCount} / ${totalCount}` : '全选'}
                        </label>
                    </div>
                    <div class="btn-group" role="group">
                        <button class="btn btn-success btn-sm" id="batch-wake"
                            ${selectedCount === 0 ? 'disabled' : ''}>
                            <i class="fas fa-play"></i> 批量唤醒
                        </button>
                        <button class="btn btn-danger btn-sm" id="batch-sleep"
                            ${selectedCount === 0 ? 'disabled' : ''}>
                            <i class="fas fa-stop"></i> 批量关机
                        </button>
                        <button class="btn btn-info btn-sm" id="batch-status"
                            ${selectedCount === 0 ? 'disabled' : ''}>
                            <i class="fas fa-refresh"></i> 批量刷新状态
                        </button>
                    </div>
                </div>
            </div>
        `;
    }

    updateBatchActionsBar() {
        const toolbar = document.getElementById('batch-toolbar');
        if (!toolbar) return;

        // 重新渲染工具栏以更新状态
        const newToolbar = this.createBatchToolbar();
        const tempDiv = document.createElement('div');
        tempDiv.innerHTML = newToolbar;
        toolbar.replaceWith(tempDiv.firstElementChild);

        // 重新绑定批量操作事件
        this.bindBatchEvents();
    }

    bindBatchEvents() {
        // 全选/取消全选
        const selectAllCheckbox = document.getElementById('select-all');
        if (selectAllCheckbox && !selectAllCheckbox.dataset.bound) {
            selectAllCheckbox.dataset.bound = 'true';
            selectAllCheckbox.addEventListener('change', (e) => {
                this.toggleSelectAll(e.target.checked);
            });
        }

        // 批量唤醒
        const batchWakeBtn = document.getElementById('batch-wake');
        if (batchWakeBtn && !batchWakeBtn.dataset.bound) {
            batchWakeBtn.dataset.bound = 'true';
            batchWakeBtn.addEventListener('click', () => this.batchWake());
        }

        // 批量关机
        const batchSleepBtn = document.getElementById('batch-sleep');
        if (batchSleepBtn && !batchSleepBtn.dataset.bound) {
            batchSleepBtn.dataset.bound = 'true';
            batchSleepBtn.addEventListener('click', () => this.batchSleep());
        }

        // 批量刷新状态
        const batchStatusBtn = document.getElementById('batch-status');
        if (batchStatusBtn && !batchStatusBtn.dataset.bound) {
            batchStatusBtn.dataset.bound = 'true';
            batchStatusBtn.addEventListener('click', () => this.batchCheckStatus());
        }

        // 单个设备复选框
        document.querySelectorAll('.batch-checkbox').forEach(checkbox => {
            if (!checkbox.dataset.bound) {
                checkbox.dataset.bound = 'true';
                checkbox.addEventListener('change', (e) => {
                    const mac = e.target.getAttribute('data-mac');
                    if (e.target.checked) {
                        this.selectedComputers.add(mac);
                    } else {
                        this.selectedComputers.delete(mac);
                    }
                    this.updateBatchActionsBar();
                    this.renderComputers();
                });
            }
        });
    }

    toggleSelectAll(checked) {
        this.filteredComputers.forEach(computer => {
            if (checked) {
                this.selectedComputers.add(computer.mac_address);
            } else {
                this.selectedComputers.delete(computer.mac_address);
            }
        });
        this.renderComputers();
    }

    async batchWake() {
        const selected = Array.from(this.selectedComputers);
        if (selected.length === 0) return;

        if (!confirm(`确定要唤醒 ${selected.length} 台设备吗？`)) return;

        try {
            const response = await api.batchWake(selected);
            if (response.success) {
                const data = response.data;
                toast.success(
                    `批量唤醒完成: 成功 ${data.success_count} / 失败 ${data.failure_count}`
                );

                // 显示详细结果
                if (data.failure_count > 0) {
                    console.error('批量唤醒失败:', data.error_messages);
                }
            } else {
                toast.error('批量唤醒失败');
            }

            // 刷新状态
            setTimeout(() => this.updateAllStatuses(), 2000);
        } catch (error) {
            toast.error('批量唤醒失败: ' + error.message);
        }
    }

    async batchSleep() {
        const selected = Array.from(this.selectedComputers);
        if (selected.length === 0) return;

        if (!confirm(`确定要关闭 ${selected.length} 台设备吗？`)) return;

        try {
            const response = await api.batchSleep(selected);
            if (response.success) {
                const data = response.data;
                toast.success(
                    `批量关机完成: 成功 ${data.success_count} / 失败 ${data.failure_count}`
                );

                // 显示详细结果
                if (data.failure_count > 0) {
                    console.error('批量关机失败:', data.error_messages);
                }
            } else {
                toast.error('批量关机失败');
            }

            // 刷新状态
            setTimeout(() => this.updateAllStatuses(), 2000);
        } catch (error) {
            toast.error('批量关机失败: ' + error.message);
        }
    }

    async batchCheckStatus() {
        const selected = Array.from(this.selectedComputers);
        if (selected.length === 0) return;

        try {
            const response = await api.batchCheckStatus(selected);
            if (response.success) {
                const results = response.data.results;
                let onlineCount = 0;
                Object.values(results).forEach(online => {
                    if (online) onlineCount++;
                });

                toast.success(
                    `状态检查完成: 在线 ${onlineCount} / 离线 ${Object.keys(results).length - onlineCount}`
                );
            }
        } catch (error) {
            toast.error('批量状态检查失败: ' + error.message);
        }
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
                toast.success(`关机指令已发送到 ${computer.name}`);
            } else {
                await api.wakeComputer(mac);
                toast.success(`WOL 魔术包已发送到 ${computer.name}`);
            }

            // 更新状态
            setTimeout(() => this.updateAllStatuses(), 2000);
        } catch (error) {
            toast.error('操作失败: ' + error.message);
        }
    }

    showAddModal() {
        // 清空表单
        document.getElementById('computer-form').reset();
        document.getElementById('computer-mac-original').value = '';
        document.getElementById('computerModalLabel').textContent = '添加设备';

        // 隐藏 TCP 端口字段
        document.getElementById('tcp-port-group').style.display = 'none';

        // 显示模态框
        this.computerModal.show();
    }

    showEditModal(mac) {
        const computer = this.computers.find(c => c.mac_address === mac);
        if (!computer) {
            toast.error('设备不存在');
            return;
        }

        // 填充表单
        document.getElementById('computer-name').value = computer.name;
        document.getElementById('computer-mac').value = computer.mac_address;
        document.getElementById('computer-mac-original').value = computer.mac_address;
        document.getElementById('computer-ip').value = computer.ip_address;
        document.getElementById('computer-type').value = computer.test_type;

        // 处理 TCP 端口字段
        const portGroup = document.getElementById('tcp-port-group');
        const portInput = document.getElementById('computer-port');

        if (computer.test_type === 'tcp') {
            portGroup.style.display = 'block';
            portInput.value = computer.tcp_port || '';
        } else {
            portGroup.style.display = 'none';
            portInput.value = '';
        }

        document.getElementById('computerModalLabel').textContent = '编辑设备';

        // 显示模态框
        this.computerModal.show();
    }

    async saveComputer() {
        const form = document.getElementById('computer-form');
        if (!form.checkValidity()) {
            form.reportValidity();
            return;
        }

        // 收集表单数据
        const computer = {
            name: document.getElementById('computer-name').value.trim(),
            mac_address: document.getElementById('computer-mac').value.trim(),
            ip_address: document.getElementById('computer-ip').value.trim(),
            test_type: document.getElementById('computer-type').value
        };

        // TCP 端口
        if (computer.test_type === 'tcp') {
            const port = document.getElementById('computer-port').value;
            if (!port) {
                toast.error('请输入端口号');
                return;
            }
            computer.test_type = port;
            computer.tcp_port = parseInt(port);
        }

        const originalMac = document.getElementById('computer-mac-original').value;
        const isEdit = originalMac !== '';

        try {
            if (isEdit) {
                // 编辑模式
                await api.updateComputer(originalMac, computer);
                toast.success('设备已更新');
            } else {
                // 添加模式
                await api.addComputer(computer);
                toast.success('设备已添加');
            }

            // 关闭模态框
            this.computerModal.hide();

            // 重新加载设备列表
            await this.loadComputers();
            this.renderComputers();

            // 刷新状态
            setTimeout(() => this.updateAllStatuses(), 1000);
        } catch (error) {
            toast.error(isEdit ? '更新失败: ' + error.message : '添加失败: ' + error.message);
        }
    }

    async deleteComputer(mac) {
        if (!confirm('确定要删除此设备吗？')) return;

        try {
            await api.deleteComputer(mac);
            toast.success('设备已删除');
            await this.loadComputers();
            this.renderComputers();
        } catch (error) {
            toast.error('删除失败: ' + error.message);
        }
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
