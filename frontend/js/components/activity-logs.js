/**
 * Activity Logs Component
 * Handles the display and interaction with activity logs
 */
class ActivityLogs {
    constructor() {
        // DOM elements
        this.logsTable = document.getElementById('activity-logs-table');
        this.logsBody = document.getElementById('activity-logs-body');
        this.entityTypeFilter = document.getElementById('entity-type-filter');
        this.actionFilter = document.getElementById('action-filter');
        this.clearLogsBtn = document.getElementById('clear-logs-btn');
        this.prevPageBtn = document.getElementById('prev-page-btn');
        this.nextPageBtn = document.getElementById('next-page-btn');
        this.pageInfo = document.getElementById('page-info');

        // Pagination state
        this.currentPage = 1;
        this.limit = 20;
        this.totalPages = 1;
        this.totalLogs = 0;

        // Filters
        this.entityType = '';
        this.action = '';

        // Initialize
        this.init();
    }

    /**
     * Initialize the component
     */
    init() {
        // Add event listeners
        this.entityTypeFilter.addEventListener('change', () => {
            this.entityType = this.entityTypeFilter.value;
            this.currentPage = 1;
            this.loadLogs();
        });

        this.actionFilter.addEventListener('change', () => {
            this.action = this.actionFilter.value;
            this.currentPage = 1;
            this.loadLogs();
        });

        this.clearLogsBtn.addEventListener('click', () => this.confirmClearOldLogs());

        this.prevPageBtn.addEventListener('click', () => {
            if (this.currentPage > 1) {
                this.currentPage--;
                this.loadLogs();
                // Scroll to top of the logs table
                this.scrollToTop();
            }
        });

        this.nextPageBtn.addEventListener('click', () => {
            if (this.currentPage < this.totalPages) {
                this.currentPage++;
                this.loadLogs();
                // Scroll to top of the logs table
                this.scrollToTop();
            }
        });

        // Load logs when the page is shown
        document.querySelector('a[data-page="activity-logs"]').addEventListener('click', () => {
            this.loadLogs();
        });

        // Initial load if we're on the activity logs page
        if (window.location.hash === '#activity-logs') {
            this.loadLogs();
        }
    }

    /**
     * Scroll to the top of the logs table
     */
    scrollToTop() {
        if (this.logsTable) {
            this.logsTable.scrollIntoView({ behavior: 'smooth' });
        }
    }

    /**
     * Load activity logs based on current filters and pagination
     */
    async loadLogs() {
        try {
            this.logsBody.innerHTML = '<tr><td colspan="5" class="loading">Loading logs...</td></tr>';
            
            // First, get the total count for pagination
            let countUrl = `/activity-logs/count`;
            
            // Apply filters if set
            if (this.entityType) {
                countUrl = `/activity-logs/entity-type/${this.entityType}/count`;
            } else if (this.action) {
                countUrl = `/activity-logs/action/${this.action}/count`;
            }
            
            try {
                const countResponse = await apiService.get(countUrl);
                this.totalLogs = countResponse.count || 0;
                this.totalPages = Math.ceil(this.totalLogs / this.limit);
            } catch (error) {
                console.error('Error getting log count:', error);
                this.totalLogs = 0;
                this.totalPages = 1;
            }
            
            // Now get the actual logs
            let url = `/activity-logs?limit=${this.limit}&offset=${(this.currentPage - 1) * this.limit}`;
            
            // Apply filters if set
            if (this.entityType) {
                url = `/activity-logs/entity-type/${this.entityType}?limit=${this.limit}&offset=${(this.currentPage - 1) * this.limit}`;
            } else if (this.action) {
                url = `/activity-logs/action/${this.action}?limit=${this.limit}&offset=${(this.currentPage - 1) * this.limit}`;
            }
            
            const response = await apiService.get(url);
            
            if (response.length === 0) {
                this.logsBody.innerHTML = '<tr><td colspan="5" class="empty-message">No logs found</td></tr>';
                this.updatePagination(0);
                return;
            }
            
            this.renderLogs(response);
            this.updatePagination(response.length);
        } catch (error) {
            console.error('Error loading logs:', error);
            toastService.show('Failed to load activity logs', 'error');
            this.logsBody.innerHTML = '<tr><td colspan="5" class="error-message">Failed to load logs</td></tr>';
        }
    }

    /**
     * Render logs in the table
     * @param {Array} logs - Array of log objects
     */
    renderLogs(logs) {
        this.logsBody.innerHTML = '';
        
        logs.forEach(log => {
            const row = document.createElement('tr');
            
            // Format timestamp
            const timestamp = new Date(log.timestamp);
            const formattedDate = timestamp.toLocaleDateString();
            const formattedTime = timestamp.toLocaleTimeString();
            
            row.innerHTML = `
                <td>${formattedDate} ${formattedTime}</td>
                <td>${this.formatAction(log.action)}</td>
                <td>${this.formatEntityType(log.entityType)}</td>
                <td>${log.description}</td>
                <td>${log.ipAddress}</td>
            `;
            
            this.logsBody.appendChild(row);
        });
    }

    /**
     * Format action for display
     * @param {string} action - Action string
     * @returns {string} Formatted action
     */
    formatAction(action) {
        const actionMap = {
            'create': '<span class="badge badge-success">Create</span>',
            'read': '<span class="badge badge-info">Read</span>',
            'update': '<span class="badge badge-warning">Update</span>',
            'delete': '<span class="badge badge-danger">Delete</span>',
            'check': '<span class="badge badge-secondary">Check</span>',
            'generate': '<span class="badge badge-primary">Generate</span>'
        };
        
        return actionMap[action] || action;
    }

    /**
     * Format entity type for display
     * @param {string} entityType - Entity type string
     * @returns {string} Formatted entity type
     */
    formatEntityType(entityType) {
        const entityMap = {
            'note': '<span class="entity-type note">Note</span>',
            'category': '<span class="entity-type category">Category</span>',
            'encryption': '<span class="entity-type encryption">Encryption</span>',
            'key': '<span class="entity-type key">Key</span>'
        };
        
        return entityMap[entityType] || entityType;
    }

    /**
     * Update pagination controls
     * @param {number} count - Number of logs loaded
     */
    updatePagination(count) {
        // Update buttons state
        this.prevPageBtn.disabled = this.currentPage === 1;
        this.nextPageBtn.disabled = this.currentPage >= this.totalPages;
        
        // Update page info
        this.pageInfo.textContent = `Page ${this.currentPage} of ${this.totalPages} (${this.totalLogs} logs)`;
    }

    /**
     * Confirm and clear old logs
     */
    confirmClearOldLogs() {
        const days = prompt('Delete logs older than how many days?', '30');
        
        if (days === null) {
            return; // User cancelled
        }
        
        const daysNum = parseInt(days);
        if (isNaN(daysNum) || daysNum <= 0) {
            toastService.show('Please enter a valid number of days', 'error');
            return;
        }
        
        if (confirm(`Are you sure you want to delete logs older than ${daysNum} days? This action cannot be undone.`)) {
            this.clearOldLogs(daysNum);
        }
    }

    /**
     * Clear logs older than specified days
     * @param {number} days - Number of days
     */
    async clearOldLogs(days) {
        try {
            const response = await apiService.delete(`/activity-logs/older-than/${days}`);
            
            toastService.show(`Successfully deleted ${response.rowsAffected} old logs`, 'success');
            
            // Reload logs
            this.currentPage = 1;
            this.loadLogs();
        } catch (error) {
            console.error('Error clearing logs:', error);
            toastService.show('Failed to clear old logs', 'error');
        }
    }
}

// Initialize the component
const activityLogs = new ActivityLogs(); 