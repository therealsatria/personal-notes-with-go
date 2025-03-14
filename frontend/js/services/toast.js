/**
 * Toast Notification Service
 * Provides methods to display toast notifications to the user
 */
class ToastService {
    constructor() {
        this.container = document.getElementById('toast-container');
        this.toasts = [];
        this.defaultDuration = 3000; // 3 seconds
    }

    /**
     * Create and show a toast notification
     * @param {string} message - The message to display
     * @param {string} type - The type of toast (success, error, info, warning)
     * @param {number} duration - How long to display the toast in ms
     */
    show(message, type = 'info', duration = this.defaultDuration) {
        // Create toast element
        const toast = document.createElement('div');
        toast.className = `toast toast-${type}`;
        toast.textContent = message;
        
        // Add to DOM
        this.container.appendChild(toast);
        this.toasts.push(toast);
        
        // Auto-remove after duration
        setTimeout(() => {
            this.remove(toast);
        }, duration);
        
        return toast;
    }
    
    /**
     * Remove a specific toast element
     * @param {HTMLElement} toast - The toast element to remove
     */
    remove(toast) {
        if (!toast) return;
        
        // Add fade-out class
        toast.style.opacity = '0';
        toast.style.transform = 'translateX(50px)';
        
        // Remove from DOM after animation
        setTimeout(() => {
            if (toast.parentNode === this.container) {
                this.container.removeChild(toast);
            }
            this.toasts = this.toasts.filter(t => t !== toast);
        }, 300);
    }
    
    /**
     * Show a success toast
     * @param {string} message - The message to display
     * @param {number} duration - How long to display the toast in ms
     */
    success(message, duration = this.defaultDuration) {
        return this.show(message, 'success', duration);
    }
    
    /**
     * Show an error toast
     * @param {string} message - The message to display
     * @param {number} duration - How long to display the toast in ms
     */
    error(message, duration = this.defaultDuration) {
        return this.show(message, 'error', duration);
    }
    
    /**
     * Show an info toast
     * @param {string} message - The message to display
     * @param {number} duration - How long to display the toast in ms
     */
    info(message, duration = this.defaultDuration) {
        return this.show(message, 'info', duration);
    }
    
    /**
     * Show a warning toast
     * @param {string} message - The message to display
     * @param {number} duration - How long to display the toast in ms
     */
    warning(message, duration = this.defaultDuration) {
        return this.show(message, 'warning', duration);
    }
    
    /**
     * Clear all toasts
     */
    clearAll() {
        this.toasts.forEach(toast => this.remove(toast));
    }
}

// Create and export a singleton instance
const toastService = new ToastService(); 