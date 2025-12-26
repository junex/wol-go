/**
 * Toast 通知系统
 * 提供成功、错误、警告、信息四种类型的非侵入式通知
 */

class ToastManager {
    constructor() {
        this.container = null;
        this.toasts = new Map(); // toastId -> element
        this.toastCounter = 0;
        this.defaultDuration = 3000; // 默认显示3秒
    }

    /**
     * 初始化 Toast 容器
     */
    init() {
        if (this.container) return;

        // 创建容器
        this.container = document.createElement('div');
        this.container.className = 'toast-container position-fixed top-0 end-0 p-3';
        this.container.style.zIndex = '1050';
        document.body.appendChild(this.container);
    }

    /**
     * 显示 Toast 通知
     * @param {string} type - 类型: success, error, warning, info
     * @param {string} message - 消息内容
     * @param {number} duration - 显示时长（毫秒），0表示不自动关闭
     */
    show(type, message, duration = null) {
        this.init();

        const toastId = `toast-${this.toastCounter++}`;
        const displayDuration = duration !== null ? duration : this.defaultDuration;

        // 创建 Toast 元素
        const toastEl = document.createElement('div');
        toastEl.className = `toast`;
        toastEl.id = toastId;
        toastEl.setAttribute('role', 'alert');
        toastEl.setAttribute('aria-live', 'assertive');
        toastEl.setAttribute('aria-atomic', 'true');

        // 图标映射
        const icons = {
            success: 'fa-check-circle',
            error: 'fa-exclamation-circle',
            warning: 'fa-exclamation-triangle',
            info: 'fa-info-circle'
        };

        const icon = icons[type] || icons.info;

        toastEl.innerHTML = `
            <div class="toast-header bg-${type} text-white">
                <i class="fas ${icon} me-2"></i>
                <strong class="me-auto">${this.getTypeText(type)}</strong>
                <button type="button" class="btn-close btn-close-white" data-bs-dismiss="toast" aria-label="Close"></button>
            </div>
            <div class="toast-body">
                ${this.escapeHtml(message)}
            </div>
        `;

        // 添加到容器
        this.container.appendChild(toastEl);

        // 使用 Bootstrap Toast API
        const bsToast = new bootstrap.Toast(toastEl, {
            delay: displayDuration,
            autohide: displayDuration > 0
        });

        // 显示 Toast
        bsToast.show();

        // 保存引用
        this.toasts.set(toastId, { element: toastEl, bsToast });

        // 监听关闭事件
        toastEl.addEventListener('hidden.bs.toast', () => {
            this.toasts.delete(toastId);
            toastEl.remove();
        });

        return toastId;
    }

    /**
     * 显示成功消息
     */
    success(message, duration) {
        return this.show('success', message, duration);
    }

    /**
     * 显示错误消息
     */
    error(message, duration) {
        return this.show('error', message, duration);
    }

    /**
     * 显示警告消息
     */
    warning(message, duration) {
        return this.show('warning', message, duration);
    }

    /**
     * 显示信息消息
     */
    info(message, duration) {
        return this.show('info', message, duration);
    }

    /**
     * 隐藏指定的 Toast
     */
    hide(toastId) {
        const toast = this.toasts.get(toastId);
        if (toast) {
            toast.bsToast.hide();
        }
    }

    /**
     * 隐藏所有 Toast
     */
    hideAll() {
        this.toasts.forEach((toast) => {
            toast.bsToast.hide();
        });
    }

    /**
     * 获取类型文本
     */
    getTypeText(type) {
        const texts = {
            success: '成功',
            error: '错误',
            warning: '警告',
            info: '提示'
        };
        return texts[type] || texts.info;
    }

    /**
     * HTML 转义
     */
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// 创建全局 Toast 管理器实例
const toast = new ToastManager();

// 页面加载后初始化
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => toast.init());
} else {
    toast.init();
}
